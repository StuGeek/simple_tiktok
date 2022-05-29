package controller

import (
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 用户信息表users
type UserDao struct {
	Id            int64  `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Name          string `json:"name"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
	Token         string `json:"token"`
}

func (UserDao) TableName() string {
	return "users"
}

// 视频信息表videos
type VideoDao struct {
	Id            int64  `json:"id,omitempty" gorm:"primary_key;AUTO_INCREMENT"`
	AuthorId      int64  `json:"author_id,omitempty"`
	PlayUrl       string `json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty"`
	Title         string `json:"title,omitempty"`
	PublishTime   int64  `json:"publish_time,omitempty"`
}

func (VideoDao) TableName() string {
	return "videos"
}

// 点赞视频信息表favorite_videos
type FavoriteVideoDao struct {
	Token   string `json:"token"`    // 用户的token
	VideoId int64  `json:"video_id"` // 用户喜欢的视频Id
}

func (FavoriteVideoDao) TableName() string {
	return "favorite_videos"
}

// 评论信息表comments
type CommentDao struct {
	Id          int64  `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	UserId      int64  `json:"user_id"`
	VideoId     int64  `json:"video_id"`
	Content     string `json:"content"`
	CreateDate  string `json:"create_date"`
	PublishTime int64  `json:"publish_time"`
}

func (CommentDao) TableName() string {
	return "comments"
}

// 关注信息表follows
type FollowDao struct {
	UserId   int64 `json:"user_id"`
	ToUserId int64 `json:"to_user_id"`
}

func (FollowDao) TableName() string {
	return "follows"
}

var (
	globalDb *gorm.DB     // 全局数据库操作指针
	dbMutex  sync.RWMutex // 操作数据库锁
)

// 初始化数据库，进行自动迁移或建表，初始化用户信息map
func InitDB() {
	dsn := sql_username + ":" + sql_password + "@tcp(127.0.0.1:3306)/" + sql_dbname + "?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	globalDb, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	globalDb.AutoMigrate(&UserDao{})
	globalDb.AutoMigrate(&FavoriteVideoDao{})
	globalDb.AutoMigrate(&VideoDao{})
	globalDb.AutoMigrate(&CommentDao{})
	globalDb.AutoMigrate(&FollowDao{})

	InitUserInfo()
}

// 初始化账号信息
func InitUserInfo() {
	// 获取所有账号信息
	var users []UserDao

	dbMutex.Lock()
	globalDb.Find(&users)
	dbMutex.Unlock()

	// 遍历所有账号
	for _, user := range users {
		// 存储每个账号的token和User的对应关系
		usersLoginInfo[user.Token] = User{
			Id:            user.Id,
			Name:          user.Name,
			FollowCount:   user.FollowCount,
			FollowerCount: user.FollowerCount,
			IsFollow:      user.IsFollow,
		}
		// 存储每个账号的Id和Token的对应关系
		userIdToToken[user.Id] = user.Token
	}
}

// 根据用户的token获取用户点赞的视频Id和Video结构体对应的map
func GetFavoriteVideoByToken(token string) map[int64]Video {
	// 用favorite_videos表和videos表查询出特定token对应用户所点赞的视频
	var favoriteVideos []VideoDao
	dbMutex.Lock()
	globalDb.Joins("inner join favorite_videos on videos.id = favorite_videos.video_id").Where("favorite_videos.token = ?", token).Find(&favoriteVideos)
	dbMutex.Unlock()

	var favoriteVideoInfo = make(map[int64]Video)

	// 存储用户点赞视频的视频Id和视频，并设置IsFavorite为true
	for _, video := range favoriteVideos {
		favoriteVideoInfo[video.Id] = Video{
			Id:            video.Id,
			Author:        usersLoginInfo[userIdToToken[video.AuthorId]],
			PlayUrl:       video.PlayUrl,
			CoverUrl:      video.CoverUrl,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			IsFavorite:    true,
			Title:         video.Title,
		}
	}

	return favoriteVideoInfo
}

// 调用Feed接口时，初始化这个用户可以获取的视频信息
func InitVideoInfo(lastestTime int64, token string) ([]Video, int64) {
	// 找到投稿时间不晚于lastestTime的投稿视频，按投稿时间倒序排列，最多30个
	var videos []VideoDao
	dbMutex.Lock()
	globalDb.Where("publish_time <= ?", lastestTime).Order("publish_time desc").Find(&videos).Limit(30)
	dbMutex.Unlock()
	var nextTime int64

	// 获取用户点赞的视频列表
	favoriteVideoInfo := GetFavoriteVideoByToken(token)

	var videoList []Video
	for _, videoDao := range videos {
		// 如果视频被点赞过，isFavorite设置为true，否则设置为false
		_, isFavorite := favoriteVideoInfo[videoDao.Id]

		videoList = append(videoList, Video{
			Id:            videoDao.Id,
			Author:        usersLoginInfo[userIdToToken[videoDao.AuthorId]],
			PlayUrl:       videoDao.PlayUrl,
			CoverUrl:      videoDao.CoverUrl,
			FavoriteCount: videoDao.FavoriteCount,
			CommentCount:  videoDao.CommentCount,
			IsFavorite:    isFavorite,
			Title:         videoDao.Title,
		})

		// 退出循环时，记录下本次返回的视频中，发布最早的时间，作为下次请求时的latest_time
		// nextTime = videoDao.PublishTime
		nextTime = time.Now().Unix()
	}

	// 返回获取的视频列表和下次请求时的latest_time
	return videoList, nextTime
}

// 根据用户的Id获取存储关注的用户Id的map
func GetFollowById(userId int64) map[int64]struct{} {
	// 从follows关注信息表中查询出这个用户关注的所有用户的Id
	var follows []FollowDao
	dbMutex.Lock()
	globalDb.Where("user_id = ?", userId).Find(&follows)
	dbMutex.Unlock()

	var followsInfo = make(map[int64]struct{})

	// 存储关注用户的用户Id
	for _, follow := range follows {
		followsInfo[follow.ToUserId] = struct{}{}
	}

	return followsInfo
}

// 根据登录用户的Id初始化账号信息，主要是设置这个账号对每个用户的IsFollow属性
func InitUserInfoById(userId int64) {
	followList := GetFollowById(userId)

	// 遍历所有存储在usersLoginInfo中账号信息
	for token, user := range usersLoginInfo {
		// 判断登录用户是否关注了这个用户
		_, isFollow := followList[user.Id]

		// 存储每个账号的token和User的对应关系，设置IsFollow属性
		usersLoginInfo[token] = User{
			Id:            user.Id,
			Name:          user.Name,
			FollowCount:   user.FollowCount,
			FollowerCount: user.FollowerCount,
			IsFollow:      isFollow,
		}
	}
}

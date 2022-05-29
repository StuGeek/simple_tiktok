package controller

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 数据库中的用户信息表users的基本元素
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

// 数据库中的视频信息表videos的基本元素
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

// 数据库中的点赞视频信息表favorite_videos的基本元素
type FavoriteVideoDao struct {
	Token   string `json:"token"`    // 用户的token
	VideoId int64  `json:"video_id"` // 用户喜欢的视频Id
}

func (FavoriteVideoDao) TableName() string {
	return "favorite_videos"
}

var globalDb *gorm.DB // 全局数据库操作指针

// 初始化数据库，进行连接，账号、评论信息初始化等工作
func InitDB() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/simple_tiktok?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	globalDb, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	InitUserInfo()
	// InitCommentInfo()
}

// 初始化账号信息
func InitUserInfo() {
	// 自动迁移，没有则以UserDao结构建表
	globalDb.AutoMigrate(&UserDao{})

	// 获取所有账号信息
	var users []UserDao
	globalDb.Find(&users)

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

// 根据用户的token获取用户点赞的视频
func GetFavoriteVideoByToken(token string) map[int64]Video {
	// 自动迁移，没有则以FavoriteVideoDao结构建表
	globalDb.AutoMigrate(&FavoriteVideoDao{})

	// 用favorite_videos表和videos表查询出特定token对应用户所点赞的视频
	var favoriteVideos []VideoDao
	globalDb.Joins("inner join favorite_videos on videos.id = favorite_videos.video_id").Where("favorite_videos.token = ?", token).Find(&favoriteVideos)

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

// 调用Feed接口时，初始化这个用户获取到的视频信息
func InitVideoInfo(lastestTime int64, token string) ([]Video, int64) {
	// 自动迁移，没有则以VideoDao结构建表
	globalDb.AutoMigrate(&VideoDao{})

	// 找到投稿时间不晚于lastestTime的投稿视频，按投稿时间倒序排列，最多30个
	var videos []VideoDao
	globalDb.Where("publish_time <= ?", lastestTime).Order("publish_time desc").Find(&videos).Limit(30)
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

// 初始化这个用户获取到的评论信息
func InitCommentInfo() {
	globalDb.AutoMigrate(&Comment{})
	// globalDb.Create(&Comment{})
}

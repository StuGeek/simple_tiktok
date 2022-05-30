package repository

import (
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var SqlUsername = "root"                 // 数据库的用户名
var SqlPassword = "123456"               // 数据库的密码
var SqlDBName = "simple_tiktok"          // 使用的数据库名
var SqlDemoDBName = "demo_simple_tiktok" // 导入demo数据使用的数据库名

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

var GlobalDB *gorm.DB    // 全局数据库操作指针
var DBMutex sync.RWMutex // 操作数据库锁

// 初始化数据库，进行自动迁移或建表，初始化用户信息map
func InitDB(dbName string) {
	dsn := SqlUsername + ":" + SqlPassword + "@tcp(127.0.0.1:3306)/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	GlobalDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	GlobalDB.AutoMigrate(&UserDao{})
	GlobalDB.AutoMigrate(&VideoDao{})
	GlobalDB.AutoMigrate(&FavoriteVideoDao{})
	GlobalDB.AutoMigrate(&CommentDao{})
	GlobalDB.AutoMigrate(&FollowDao{})
}

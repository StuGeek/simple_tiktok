package controller

import (
	"sync"

	"github.com/RaymondCode/simple-demo/repository"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var serverUrl = "http://172.26.44.188:8080/" // 服务器的url

var sqlUsername = "root"                 // 数据库的用户名
var sqlPassword = "123456"               // 数据库的密码
var sqlDBName = "simple_tiktok"          // 使用的数据库名
var sqlDemoDBName = "demo_simple_tiktok" // 导入demo数据使用的数据库名

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

type Video struct {
	Id            int64  `json:"id,omitempty"`
	Author        User   `json:"author,omitempty"`
	PlayUrl       string `json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty"`
	Title         string `json:"title,omitempty"`
}

type Comment struct {
	Id         int64  `json:"id"`
	User       User   `json:"user"`
	Content    string `json:"content"`
	CreateDate string `json:"create_date"`
}

type User struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

var (
	globalDb *gorm.DB     // 全局数据库操作指针
	dbMutex  sync.RWMutex // 操作数据库锁
)

// 初始化数据库，进行自动迁移或建表，初始化用户信息map
func InitDB() {
	dsn := sqlUsername + ":" + sqlPassword + "@tcp(127.0.0.1:3306)/" + sqlDBName + "?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	globalDb, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	globalDb.AutoMigrate(&repository.UserDao{})
	globalDb.AutoMigrate(&repository.VideoDao{})
	globalDb.AutoMigrate(&repository.FavoriteVideoDao{})
	globalDb.AutoMigrate(&repository.CommentDao{})
	globalDb.AutoMigrate(&repository.FollowDao{})

	InitUserInfo()
}

package repository

import (
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"simple_tiktok/global"
)

var (
	globalDB       *gorm.DB   // 全局数据库操作指针
	usersMutex     sync.Mutex // 操作users表时用到的锁
	videosMutex    sync.Mutex // 操作videos表时用到的锁
	relationsMutex sync.Mutex // 操作relations表时用到的锁
	favoritesMutex sync.Mutex // 操作favorites表时用到的锁
	commentsMutex  sync.Mutex // 操作comments表时用到的锁
)

// 初始化数据库，进行自动迁移或建表
func InitDB(dbName string) {
	dsn := global.SqlUsername + ":" + global.SqlPassword + "@tcp(127.0.0.1:3306)/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	globalDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	sqlDB, err := globalDB.DB()
	if err != nil {
		panic("connect db server failed.")
	}

	// 使用连接池
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Second * 600)

	globalDB.AutoMigrate(&UserDao{}, &VideoDao{}, &FavoriteDao{}, &CommentDao{}, &RelationDao{})
}

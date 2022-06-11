package repository

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"simple_tiktok/global"
)

var GlobalDB *gorm.DB // 全局数据库操作指针

// 初始化数据库，进行自动迁移或建表
func InitDB(dbName string) {
	dsn := global.SqlUsername + ":" + global.SqlPassword + "@tcp(127.0.0.1:3306)/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	GlobalDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	sqlDB, err := GlobalDB.DB()
	if err != nil {
		panic("connect db server failed.")
	}

	// 使用连接池
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Second * 600)

	GlobalDB.AutoMigrate(&UserDao{}, &VideoDao{}, &FavoriteDao{}, &CommentDao{}, &RelationDao{})
}

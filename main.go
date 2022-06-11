package main

import (
	"os"
	"simple_tiktok/global"
	"simple_tiktok/repository"
	"simple_tiktok/service"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	initRouter(r)

	// 如果命令行参数为--demo，使用demo数据库，并导入demo数据
	if len(os.Args) == 1 {
		// 初始化数据库
		repository.InitDB(global.SqlDBName)
	} else if len(os.Args) == 2 && os.Args[1] == "--demo" {
		// 初始化Demo数据库
		repository.InitDB(global.SqlDemoDBName)
		// 向数据库导入Demo数据
		repository.InitDemoData()
	}

	// 初始化账号信息
	service.InitUserInfo()

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

package main

import (
	"os"
	"simple_tiktok/controller"
	"simple_tiktok/repository"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	initRouter(r)

	// 如果命令行参数为--demo，使用demo数据库，并导入demo数据
	if len(os.Args) == 1 {
		// 初始化数据库
		repository.InitDB(repository.SqlDBName)
	} else if len(os.Args) == 2 && os.Args[1] == "--demo" {
		// 初始化Demo数据库
		repository.InitDB(repository.SqlDemoDBName)
		// 向数据库导入Demo数据
		controller.InitDemoData()
	}

	// 初始化账号信息
	controller.InitUserInfo()

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

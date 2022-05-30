package main

import (
	"github.com/RaymondCode/simple-demo/controller"
	"github.com/RaymondCode/simple-demo/repository"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	initRouter(r)
	// 初始化数据库
	repository.InitDB(repository.SqlDemoDBName)
	// 初始化导入Demo数据的数据库
	controller.InitDemoData()
	// 初始化账号信息
	controller.InitUserInfo()

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

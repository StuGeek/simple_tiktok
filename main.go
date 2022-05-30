package main

import (
	"github.com/RaymondCode/simple-demo/controller"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	initRouter(r)
	// 初始化数据库
	// controller.InitDB()
	// 初始化导入Demo数据的数据库
	controller.InitDemoData()

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

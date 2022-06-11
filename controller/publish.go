package controller

import (
	"net/http"
	"strconv"

	"simple_tiktok/global"
	"simple_tiktok/service"

	"github.com/gin-gonic/gin"
)

type VideoListResponse struct {
	global.Response
	VideoList []global.Video `json:"video_list"`
}

// 发布视频行为，需要用户是登录状态
func Publish(c *gin.Context) {
	token := c.PostForm("token")
	// 判断是否处于登录状态，不是则直接返回，取消发布视频
	if _, exist := service.GetExistUserByToken(token); !exist {
		c.JSON(http.StatusOK, global.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}

	var finalName string
	var err error

	// 调用service层的发布视频服务，并获取保存的视频文件名，如果出错，返回错误响应
	if finalName, err = service.PublishVideoAction(token, c); err != nil {
		c.JSON(http.StatusOK, global.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	// 调用service层没出错则返回响应成功
	c.JSON(http.StatusOK, global.Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
}

// 根据user_id返回发布作品列表
func PublishList(c *gin.Context) {
	userIdStr := c.Query("user_id")
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)

	// 根据用户id调用service层的获取用户发布的视频列表服务，如果出错，返回服务器错误响应
	videoList, err := service.GetPublishListById(userId)
	if err != nil {
		c.JSON(http.StatusOK, global.Response{StatusCode: 1, StatusMsg: "Internal Server Error! Get publish list failed"})
		return
	}

	// 调用service层没出错则返回响应成功和作品视频列表
	c.JSON(http.StatusOK, VideoListResponse{
		Response: global.Response{
			StatusCode: 0,
		},
		VideoList: videoList,
	})
}

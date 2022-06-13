package controller

import (
	"net/http"
	"strconv"

	"simple_tiktok/global"
	"simple_tiktok/service"

	"github.com/gin-gonic/gin"
)

// 点赞或取消点赞行为
func FavoriteAction(c *gin.Context) {
	token := c.Query("token")
	videoIdStr := c.Query("video_id")
	actionType := c.Query("action_type")
	videoId, err := strconv.ParseInt(videoIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, global.Response{StatusCode: 1, StatusMsg: "Get videoId error"})
		return
	}

	// 如果是点赞行为
	if actionType == "1" {
		// 根据点赞用户的token和点赞视频的Id调用service层的点赞服务，如果出错，返回错误信息响应
		if errMsg := service.LikeAction(token, videoId); errMsg != "" {
			c.JSON(http.StatusOK, global.Response{StatusCode: 1, StatusMsg: errMsg})
			return
		}
	} else if actionType == "2" {
		// 根据用户的token和视频的Id调用service层的取消点赞服务，如果出错，返回错误信息响应
		if errMsg := service.CancelLikeAction(token, videoId); errMsg != "" {
			c.JSON(http.StatusOK, global.Response{StatusCode: 1, StatusMsg: errMsg})
			return
		}
	}
	// 调用service层没出错则返回响应成功
	c.JSON(http.StatusOK, global.Response{StatusCode: 0})
}

// 返回点赞喜欢列表
func FavoriteList(c *gin.Context) {
	userIdStr := c.Query("user_id")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, global.Response{StatusCode: 1, StatusMsg: "Get userId error"})
		return
	}

	// 根据用户的Id调用service层的获取喜欢视频列表服务，如果出错，返回服务器错误响应
	videoList, errMsg := service.GetFavoriteListById(userId)
	if errMsg != "" {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: global.Response{StatusCode: 1, StatusMsg: errMsg},
		})
		return
	}

	// 调用service层没出错则返回响应成功和喜欢视频列表
	c.JSON(http.StatusOK, VideoListResponse{
		Response:  global.Response{StatusCode: 0},
		VideoList: videoList,
	})
}

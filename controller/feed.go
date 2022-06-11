package controller

import (
	"net/http"
	"strconv"
	"time"

	"simple_tiktok/global"
	"simple_tiktok/service"

	"github.com/gin-gonic/gin"
)

type FeedResponse struct {
	global.Response
	VideoList []global.Video `json:"video_list"`
	NextTime  int64          `json:"next_time"`
}

// 获取按投稿时间倒序的视频列表，单次最多30个
func Feed(c *gin.Context) {
	latestTimeStr := c.Query("latest_time")
	var latestTime int64
	// latest_time不填表示当前时间
	if latestTimeStr != "" {
		lTime, _ := strconv.ParseInt(latestTimeStr, 10, 64)
		latestTime = lTime
	} else {
		latestTime = time.Now().Unix()
	}

	token := c.Query("token")

	// 根据最新投稿时间戳和用户token，调用service层的获取视频列表服务，并获取下次请求时的latest_time
	videoList, nextTime, err := service.GetFeedVideoList(latestTime, token)
	// 如果服务出错，返回服务器错误响应
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: global.Response{StatusCode: 1, StatusMsg: "Internal Server Error! Get feed videos failed"},
		})
		return
	}

	// 调用service层没出错则返回响应成功，以及视频列表和下次请求时的latest_time
	c.JSON(http.StatusOK, FeedResponse{
		Response:  global.Response{StatusCode: 0},
		VideoList: videoList,
		NextTime:  nextTime,
	})
}

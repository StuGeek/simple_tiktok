package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list"`
	NextTime  int64   `json:"next_time"`
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

	// 处在登录状态每次Feed时，都需要初始化一次账号信息，更新当前用户的关注状态
	token := c.Query("token")
	if token != "" {
		InitUserInfoById(usersLoginInfo[token].Id)
	}

	// 根据最新投稿时间戳和用户token，返回用户视频列表和下次请求时的latest_time
	videoList, nextTime := InitVideoInfo(latestTime, token)

	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: videoList,
		NextTime:  nextTime,
	})
}

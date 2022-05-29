package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
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

	// 根据最新投稿时间戳和用户token，返回用户视频列表和下次请求时的latest_time
	videoList, nextTime := InitVideoInfo(latestTime, token)

	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: videoList,
		NextTime:  nextTime,
	})
}

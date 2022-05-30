package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/RaymondCode/simple-demo/repository"
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

// 调用Feed接口时，初始化这个用户可以获取的视频信息
func InitVideoInfo(lastestTime int64, token string) ([]Video, int64) {
	// 找到投稿时间不晚于lastestTime的投稿视频，按投稿时间倒序排列，最多30个
	var videos []repository.VideoDao
	dbMutex.Lock()
	globalDb.Where("publish_time <= ?", lastestTime).Order("publish_time desc").Find(&videos).Limit(30)
	dbMutex.Unlock()
	var nextTime int64

	// 获取用户点赞的视频列表
	favoriteVideoInfo := GetFavoriteVideoByToken(token)

	var videoList []Video
	for _, videoDao := range videos {
		// 如果视频被点赞过，isFavorite设置为true，否则设置为false
		_, isFavorite := favoriteVideoInfo[videoDao.Id]

		videoList = append(videoList, Video{
			Id:            videoDao.Id,
			Author:        usersLoginInfo[userIdToToken[videoDao.AuthorId]],
			PlayUrl:       videoDao.PlayUrl,
			CoverUrl:      videoDao.CoverUrl,
			FavoriteCount: videoDao.FavoriteCount,
			CommentCount:  videoDao.CommentCount,
			IsFavorite:    isFavorite,
			Title:         videoDao.Title,
		})

		// 退出循环时，记录下本次返回的视频中，发布最早的时间，作为下次请求时的latest_time
		nextTime = videoDao.PublishTime
		// nextTime = time.Now().Unix()
	}

	// 返回获取的视频列表和下次请求时的latest_time
	return videoList, nextTime
}

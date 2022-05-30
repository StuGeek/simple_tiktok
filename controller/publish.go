package controller

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/RaymondCode/simple-demo/repository"
	"github.com/gin-gonic/gin"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

// 发布视频行为，需要用户是登录状态
func Publish(c *gin.Context) {
	token := c.PostForm("token")
	// 判断是否处于登录状态，不是则直接返回，取消发布视频
	if _, exist := usersLoginInfo[token]; !exist {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}

	// 获取发布的视频文件数据
	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	// 将视频文件经过文件路径和命名处理后存入本地
	filename := filepath.Base(data.Filename)
	user := usersLoginInfo[token]
	finalName := fmt.Sprintf("%d_%s", user.Id, filename)
	saveFile := filepath.Join("./public/", finalName)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	// 获取视频标题
	title := c.PostForm("title")
	repository.DBMutex.Lock()
	// 将视频信息存入数据库中，投稿时间为当前时间
	repository.GlobalDB.Create(&repository.VideoDao{
		AuthorId:      user.Id,
		PlayUrl:       serverUrl + "static/" + finalName,
		FavoriteCount: 0,
		CommentCount:  0,
		IsFavorite:    false,
		Title:         title,
		PublishTime:   time.Now().Unix(),
	})
	repository.DBMutex.Unlock()

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
}

// 返回发布作品列表
func PublishList(c *gin.Context) {
	token := c.Query("token")
	userIdStr := c.Query("user_id")

	userId, _ := strconv.ParseInt(userIdStr, 10, 64)

	// 从数据库中根据用户id获取这个用户发布的视频列表
	var videoDaoList []repository.VideoDao
	repository.DBMutex.Lock()
	repository.GlobalDB.Where("author_id = ?", userId).Find(&videoDaoList)
	repository.DBMutex.Unlock()

	// 获取这个用户点赞的视频列表
	favoriteVideoInfo := GetFavoriteVideoByToken(token)

	var videoList []Video
	for _, videoDao := range videoDaoList {
		// 判断这个用户是否有给自己发布的视频点赞
		_, isFavorite := favoriteVideoInfo[videoDao.Id]

		// 将从数据库中取出的视频列表加入到videoList中，并最后返回
		videoList = append(videoList, Video{
			Id:            videoDao.Id,
			Author:        usersLoginInfo[token],
			PlayUrl:       videoDao.PlayUrl,
			CoverUrl:      videoDao.CoverUrl,
			FavoriteCount: videoDao.FavoriteCount,
			CommentCount:  videoDao.CommentCount,
			IsFavorite:    isFavorite,
			Title:         videoDao.Title,
		})
	}

	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: videoList,
	})
}

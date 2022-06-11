package service

import (
	"fmt"
	"path/filepath"
	"simple_tiktok/global"
	"simple_tiktok/repository"
	"sync"

	"github.com/gin-gonic/gin"
)

// 投稿视频，返回投稿的文件名和可能的错误
func PublishVideoAction(token string, c *gin.Context) (string, error) {
	// 获取发布的视频文件数据
	data, err := c.FormFile("data")
	if err != nil {
		return "", err
	}

	// 将视频文件经过文件路径和命名处理后存入本地
	filename := filepath.Base(data.Filename)
	userId, _ := repository.GetUserIdByToken(token)
	finalName := fmt.Sprintf("%d_%s", userId, filename)
	saveFile := filepath.Join("./public/", finalName)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		return "", err
	}

	// 获取视频标题
	title := c.PostForm("title")
	// 将视频信息存入数据库中，投稿时间为当前时间
	repository.CreateVideo(userId, finalName, title)

	return finalName, nil
}

// 根据作者Id获取发布视频列表
func GetPublishListById(authorId int64) ([]global.Video, error) {
	var videoDaoList []repository.VideoDao
	var favoriteVideoInfo map[int64]global.Video
	var err error

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		// 从数据库中根据用户id获取这个用户发布的视频列表
		videoDaoList, err = repository.QueryVideoByAuthorId(authorId)
	}()
	go func() {
		defer wg.Done()
		// 获取这个用户点赞的视频列表
		authorToken, _ := repository.GetUserTokenById(authorId)
		favoriteVideoInfo, err = GetFavoriteVideosByToken(authorToken)
	}()

	wg.Wait()

	if err != nil {
		return []global.Video{}, err
	}

	// 发布视频列表的作者为同一个
	author, _ := repository.GetUserById(authorId)

	var videoList []global.Video
	for _, videoDao := range videoDaoList {
		// 判断这个用户是否有给自己发布的视频点赞
		_, isFavorite := favoriteVideoInfo[videoDao.Id]

		// 将从数据库中取出的视频列表加入到videoList中，并最后返回
		videoList = append(videoList, global.Video{
			Id:            videoDao.Id,
			Author:        author,
			PlayUrl:       videoDao.PlayUrl,
			CoverUrl:      videoDao.CoverUrl,
			FavoriteCount: videoDao.FavoriteCount,
			CommentCount:  videoDao.CommentCount,
			IsFavorite:    isFavorite,
			Title:         videoDao.Title,
		})
	}

	return videoList, nil
}

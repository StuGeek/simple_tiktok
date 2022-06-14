package service

import (
	"fmt"
	"path/filepath"
	"simple_tiktok/global"
	"simple_tiktok/repository"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// 投稿视频，返回投稿的文件名和可能的错误信息
func PublishVideoAction(token string, c *gin.Context) (string, string) {
	// 用户要处于登录状态，登录凭证需要有效
	user, errMsg := GetUserByToken(token)
	if errMsg != "" {
		return "", errMsg
	}

	// 获取发布的视频文件数据
	data, err := c.FormFile("data")
	if err != nil {
		return "", "Internal Server Error! Get file data fail"
	}

	// 将视频文件经过文件路径和命名处理后存入本地，加用户序号和时间戳避免同名文件覆盖
	filename := filepath.Base(data.Filename)
	finalName := fmt.Sprintf("%d_%d_%s", user.Id, time.Now().Unix(), filename)
	saveFile := filepath.Join(global.SavaFilePath, finalName)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		return "", "Internal Server Error! Save file fail"
	}

	// 获取视频标题
	title := c.PostForm("title")
	// 将视频信息存入数据库中，投稿时间为当前时间
	if err := repository.CreateVideo(user.Id, finalName, title); err != nil {
		return "", "Internal Server Error! Create video record fail"
	}

	return finalName, ""
}

// 根据token和作者Id获取发布视频列表
func GetPublishList(token string, authorId int64) ([]global.Video, string) {
	var tokenUserFavoriteMap = make(map[int64]struct{})
	var publishVideoList []repository.VideoDao
	var errMsg string = ""

	// 根据用户Id获取作者的User结构体
	authorDao, err := repository.QueryUserById(authorId)
	if err != nil {
		return nil, "Internal Server Error! Query user failed"
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		// 获取登录用户
		user, errMsg := GetUserByToken(token)
		// 查看发布视频列表可以不处在登录状态，不用返回错误信息
		if errMsg == "" {
			// 如果登录凭证有效，获取登录用户的点赞列表
			tokenUserFavoriteList, err := repository.QueryAllFavoriteVideoByUserId(user.Id)
			if err != nil {
				errMsg = "Internal Server Error! Query follow users failed"
			} else {
				// 记录登录用户点赞的视频Id
				for _, favorite := range tokenUserFavoriteList {
					tokenUserFavoriteMap[favorite.Id] = struct{}{}
				}
			}

		}
	}()
	go func() {
		defer wg.Done()
		// 获取查看的这个用户发布的视频列表
		publishVideoList, err = repository.QueryVideoByAuthorId(authorId)
		if err != nil {
			errMsg = "Internal Server Error! Query video failed"
		}
	}()

	wg.Wait()

	if errMsg != "" {
		return nil, errMsg
	}

	// 发布视频列表的作者为同一个
	author := userDaoToUser(&authorDao)

	// 将从数据库中取出的视频列表加入到videoList中，并根据登录用户是否点赞设置IsFavorite属性最后返回
	var videoList []global.Video
	for _, videoDao := range publishVideoList {
		// 判断登录用户是否有给查看用户发布的视频点赞
		_, isFavorite := tokenUserFavoriteMap[videoDao.Id]
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

	return videoList, ""
}

package service

import (
	"simple_tiktok/global"
	"simple_tiktok/repository"
	"sync"
)

// 根据用户token和视频Id，进行点赞，返回可能的错误信息
func LikeAction(token string, videoId int64) string {
	// 点赞需要用户已经登录，且登录凭证有效
	user, errMsg := GetUserByToken(token)
	if errMsg != "" {
		return errMsg
	}

	// 获取用户是否已经给这个视频点过赞
	isFavorite, err := repository.QueryIsFavorite(user.Id, videoId)
	if err != nil {
		return "Internal Server Error! Query favorite failed"
	}
	// 如果之前已经点过赞了，直接返回
	if isFavorite {
		return "The video has been liked"
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		// 将数据库中videos视频信息表的这个视频的总点赞数加一，如果更新失败，记录错误信息
		err := repository.AddOneVideoFavoriteCountById(videoId)
		if err != nil {
			errMsg = "Internal Server Error! Add one video favorite count failed"
		}
	}()
	go func() {
		defer wg.Done()
		// 在数据库的favorites点赞视频表中创建相应点赞记录，如果插入失败，记录错误信息
		err := repository.CreateFavorite(user.Id, videoId)
		if err != nil {
			errMsg = "Internal Server Error! Create favorite failed"
		}
	}()
	wg.Wait()

	// 如果有错误信息，说明出错，返回错误信息
	if errMsg != "" {
		return errMsg
	}

	// 没有错误信息则返回空字符串
	return ""
}

// 根据用户token和视频Id，取消点赞，返回可能的错误信息
func CancelLikeAction(token string, videoId int64) string {
	// 取消点赞需要用户已经登录，且登录凭证有效
	user, errMsg := GetUserByToken(token)
	if errMsg != "" {
		return errMsg
	}

	// 获取用户是否已经给这个视频点过赞
	isFavorite, err := repository.QueryIsFavorite(user.Id, videoId)
	if err != nil {
		return "Internal Server Error! Query favorite failed"
	}
	// 如果之前没有点过赞，直接返回
	if !isFavorite {
		return "Didn't like the video before"
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		// 将数据库中videos视频信息表的这个视频的总点赞数减一，如果更新失败，记录错误信息
		err := repository.SubOneVideoFavoriteCountById(videoId)
		if err != nil {
			errMsg = "Internal Server Error! Sub one video favorite count failed"
		}
	}()
	go func() {
		defer wg.Done()
		// 删除点赞记录，如果删除失败，记录错误信息
		err := repository.DeleteFavorite(user.Id, videoId)
		if err != nil {
			errMsg = "Internal Server Error! Delete favorite failed"
		}
	}()
	wg.Wait()

	// 如果有错误信息，说明出错，返回错误信息
	if errMsg != "" {
		return errMsg
	}

	// 没有错误信息则返回空字符串
	return ""
}

// 根据用户Id获取用户的点赞视频列表，返回视频列表和可能的错误信息
func GetFavoriteListById(userId int64) ([]global.Video, string) {
	// 获取点赞的视频列表
	favoriteList, err := repository.QueryAllFavoriteVideoByUserId(userId)
	if err != nil {
		return nil, "Internal Server Error! Query favorite failed"
	}

	// 将所有点赞的视频加入videoList中，设置IsFavorite属性为true，并最后返回
	var videoList []global.Video
	for _, favoriteVideo := range favoriteList {
		author := global.User{
			Id:   favoriteVideo.AuthorId,
			Name: favoriteVideo.Name,
		}

		videoList = append(videoList, global.Video{
			Id:            favoriteVideo.Id,
			Author:        author,
			PlayUrl:       favoriteVideo.PlayUrl,
			CoverUrl:      favoriteVideo.CoverUrl,
			FavoriteCount: favoriteVideo.FavoriteCount,
			CommentCount:  favoriteVideo.CommentCount,
			IsFavorite:    true,
			Title:         favoriteVideo.Title,
		})
	}

	return videoList, ""
}

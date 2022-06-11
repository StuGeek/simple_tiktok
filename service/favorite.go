package service

import (
	"simple_tiktok/global"
	"simple_tiktok/repository"
	"sync"
)

// 根据用户token和视频Id，进行点赞，返回可能的错误信息
func LikeAction(token string, videoId int64) string {
	var errMsg string = ""

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
		// 在数据库的favorite_videos点赞视频表中创建相应点赞记录，如果插入失败，记录错误信息
		err := repository.CreateFavorite(token, videoId)
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
	var errMsg string = ""

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
		err := repository.DeleteFavorite(token, videoId)
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

// 根据用户Id获取用户的喜欢视频列表，返回视频列表和可能的错误
func GetFavoriteVideosById(userId int64) ([]global.Video, error) {
	// 根据token获取这个用户点赞的视频列表
	userToken, _ := repository.GetUserTokenById(userId)
	favoriteVideos, err := repository.QueryFavoriteVideosByToken(userToken)
	if err != nil {
		return []global.Video{}, err
	}

	// 将所有点赞的视频加入videoList中，设置IsFavorite属性为true，并最后返回
	var videoList []global.Video
	for _, favoriteVideo := range favoriteVideos {
		author, _ := repository.GetUserById(favoriteVideo.AuthorId)
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

	return videoList, nil
}

// 根据用户的token获取用户点赞的视频Id和Video结构体对应的map
func GetFavoriteVideosByToken(token string) (map[int64]global.Video, error) {
	// 根据token查询对应用户所点赞的视频
	favoriteVideos, err := repository.QueryFavoriteVideosByToken(token)
	if err != nil {
		return map[int64]global.Video{}, err
	}

	var favoriteVideoInfo = make(map[int64]global.Video)

	// 存储用户点赞视频的视频Id和视频，并设置IsFavorite为true
	for _, favoriteVideo := range favoriteVideos {
		author, _ := repository.GetUserById(favoriteVideo.AuthorId)
		favoriteVideoInfo[favoriteVideo.Id] = global.Video{
			Id:            favoriteVideo.Id,
			Author:        author,
			PlayUrl:       favoriteVideo.PlayUrl,
			CoverUrl:      favoriteVideo.CoverUrl,
			FavoriteCount: favoriteVideo.FavoriteCount,
			CommentCount:  favoriteVideo.CommentCount,
			IsFavorite:    true,
			Title:         favoriteVideo.Title,
		}
	}

	return favoriteVideoInfo, nil
}

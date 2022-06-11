package service

import (
	"simple_tiktok/global"
	"simple_tiktok/repository"
	"sync"
)

// 获取投稿时间不晚于lastestTime的投稿视频，按投稿时间倒序排列，最多30个
func GetFeedVideoList(latestTime int64, token string) ([]global.Video, int64, error) {
	var videos []repository.VideoDao
	var favoriteVideoInfo map[int64]global.Video
	var err error = nil

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		// 根据最新投稿时间戳获取视频列表
		videos, err = repository.QueryVideoByPublishTime(latestTime)
	}()
	go func() {
		defer wg.Done()
		// 获取用户点赞的视频列表
		favoriteVideoInfo, err = GetFavoriteVideosByToken(token)
	}()

	wg.Wait()

	// 查询出错则返回错误
	if err != nil {
		return []global.Video{}, 0, err
	}

	// 没有视频则返回nil
	if len(videos) == 0 {
		return nil, 0, nil
	}

	var videoList []global.Video
	for _, videoDao := range videos {
		// 如果视频被点赞过，isFavorite设置为true，否则设置为false
		_, isFavorite := favoriteVideoInfo[videoDao.Id]
		author, _ := repository.GetUserById(videoDao.AuthorId)

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

	// 返回获取的视频列表和下次请求时的latest_time
	return videoList, videos[len(videos)-1].PublishTime, nil
}

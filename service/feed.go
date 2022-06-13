package service

import (
	"simple_tiktok/global"
	"simple_tiktok/repository"
	"sync"
)

// 获取投稿时间不晚于lastestTime的投稿视频，按投稿时间倒序排列，最多30个，返回视频列表和下次请求时的latest_time
func GetFeedVideoList(latestTime int64, token string) ([]global.Video, int64, string) {
	var favoriteVideoMap = map[int64]struct{}{}
	var feedVideos []repository.VideoJoinUser
	var errMsg string = ""

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		var err error
		// 根据最新投稿时间戳获取视频列表
		feedVideos, err = repository.QueryVideoByPublishTime(latestTime)
		if err != nil {
			errMsg = "Internal Server Error! Query video failed"
		}
	}()
	go func() {
		defer wg.Done()
		// 获取用户点赞的视频列表
		user, _, err := repository.QueryUserByToken(token)
		if err == nil {
			favoriteVideoList, err := repository.QueryAllFavoriteVideoByUserId(user.Id)
			if err != nil {
				errMsg = "Internal Server Error! Query favorite video failed"
			}

			for _, video := range favoriteVideoList {
				favoriteVideoMap[video.Id] = struct{}{}
			}
		} else {
			errMsg = "Internal Server Error! Query user failed"
		}
	}()

	wg.Wait()

	// 查询出错则返回错误
	if errMsg != "" {
		return []global.Video{}, 0, errMsg
	}

	// 没有视频则返回nil
	if len(feedVideos) == 0 {
		return nil, 0, "There are no videos"
	}

	var videoList []global.Video
	for _, videoDao := range feedVideos {
		// 如果视频被点赞过，isFavorite设置为true，否则设置为false
		_, isFavorite := favoriteVideoMap[videoDao.Id]

		author := global.User{
			Id:   videoDao.AuthorId,
			Name: videoDao.Name,
		}
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
	return videoList, feedVideos[len(feedVideos)-1].PublishTime, ""
}

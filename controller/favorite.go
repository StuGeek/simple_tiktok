package controller

import (
	"net/http"
	"strconv"

	"github.com/RaymondCode/simple-demo/repository"
	"github.com/gin-gonic/gin"
)

// 点赞或取消点赞行为
func FavoriteAction(c *gin.Context) {
	token := c.Query("token")
	videoIdStr := c.Query("video_id")
	actionType := c.Query("action_type")

	videoId, _ := strconv.ParseInt(videoIdStr, 10, 64)

	// 根据token获取这个用户点赞的视频列表
	favoriteVideos := GetFavoriteVideoByToken(token)

	// 点赞需要用户已经登录
	if _, exist := usersLoginInfo[token]; exist {
		var video repository.VideoDao
		// 获取用户是否已经给这个视频点过赞
		_, isFavorite := favoriteVideos[videoId]
		// 如果是点赞行为且之前没有给这个视频点过赞
		if actionType == "1" && !isFavorite {
			// 将数据库中videos视频信息表的这个视频的总点赞数加一
			repository.DBMutex.Lock()
			repository.GlobalDB.Where("id = ?", videoId).First(&video).Update("favorite_count", video.FavoriteCount+1).Update("is_favorite", true)
			// 在数据库的favorite_videos点赞视频表中创建相应点赞记录
			repository.GlobalDB.Create(&repository.FavoriteVideoDao{
				Token:   token,
				VideoId: videoId,
			})
			repository.DBMutex.Unlock()
		} else if actionType == "2" && isFavorite {
			repository.DBMutex.Lock()
			// 如果是取消点赞行为且之前给这个视频点过赞了，更新视频总点赞数，删除点赞记录
			repository.GlobalDB.Where("id = ?", videoId).First(&video).Update("favorite_count", video.FavoriteCount-1).Update("is_favorite", false)
			repository.GlobalDB.Where("token = ? and video_id = ?", token, videoId).Delete(&repository.FavoriteVideoDao{})
			repository.DBMutex.Unlock()
		}

		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// 返回点赞喜欢列表
func FavoriteList(c *gin.Context) {
	token := c.Query("token")

	// 获取这个用户点赞的视频列表
	favoriteVideos := GetFavoriteVideoByToken(token)

	// 将所有点赞的视频加入videoList中并最后返回
	var videoList []Video
	for _, favoriteVideo := range favoriteVideos {
		videoList = append(videoList, favoriteVideo)
	}

	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: videoList,
	})
}

// 根据用户的token获取用户点赞的视频Id和Video结构体对应的map
func GetFavoriteVideoByToken(token string) map[int64]Video {
	// 用favorite_videos表和videos表查询出特定token对应用户所点赞的视频
	var favoriteVideos []repository.VideoDao
	repository.DBMutex.Lock()
	repository.GlobalDB.Joins("inner join favorite_videos on videos.id = favorite_videos.video_id").Where("favorite_videos.token = ?", token).Find(&favoriteVideos)
	repository.DBMutex.Unlock()

	var favoriteVideoInfo = make(map[int64]Video)

	// 存储用户点赞视频的视频Id和视频，并设置IsFavorite为true
	for _, video := range favoriteVideos {
		favoriteVideoInfo[video.Id] = Video{
			Id:            video.Id,
			Author:        usersLoginInfo[userIdToToken[video.AuthorId]],
			PlayUrl:       video.PlayUrl,
			CoverUrl:      video.CoverUrl,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			IsFavorite:    true,
			Title:         video.Title,
		}
	}

	return favoriteVideoInfo
}

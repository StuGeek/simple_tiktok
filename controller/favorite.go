package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// FavoriteAction no practical effect, just check if token is valid
// 点赞行为
func FavoriteAction(c *gin.Context) {
	token := c.Query("token")
	videoIdStr := c.Query("video_id")
	actionTypeStr := c.Query("action_type")

	videoId, _ := strconv.ParseInt(videoIdStr, 10, 64)
	actionType, _ := strconv.ParseInt(actionTypeStr, 10, 64)

	// 根据token获取这个用户点赞的视频列表
	favoriteVideos := GetFavoriteVideoByToken(token)

	// 如果用户存在
	if _, exist := usersLoginInfo[token]; exist {
		var video VideoDao
		// 获取用户是否已经给这个视频点过赞
		_, isFavorite := favoriteVideos[videoId]
		// 如果是点赞行为且之前没有给这个视频点过赞
		if actionType == 1 && !isFavorite {
			// 将数据库中videos视频信息表的这个视频的总点赞数加一
			globalDb.Where("id = ?", videoId).First(&video).Update("favorite_count", video.FavoriteCount+1).Update("is_favorite", true)
			// 在数据库的favorite_videos点赞视频表中创建相应点赞记录
			globalDb.Create(&FavoriteVideoDao{
				Token:   token,
				VideoId: videoId,
			})
		} else if actionType == 2 && isFavorite {
			// 如果是取消点赞行为且之前给这个视频点过赞了，更新视频总点赞数，删除点赞记录
			globalDb.Where("id = ?", videoId).First(&video).Update("favorite_count", video.FavoriteCount-1).Update("is_favorite", false)
			globalDb.Where("token = ? and video_id = ?", token, videoId).Delete(&FavoriteVideoDao{})
		}

		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		// 用户不存在则不能点赞
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// FavoriteList all users have same favorite video list
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

package controller

import (
	"net/http"
	"strconv"

	"simple_tiktok/global"
	"simple_tiktok/service"

	"github.com/gin-gonic/gin"
)

// 点赞或取消点赞行为
func FavoriteAction(c *gin.Context) {
	token := c.Query("token")
	videoIdStr := c.Query("video_id")
	actionType := c.Query("action_type")
	videoId, _ := strconv.ParseInt(videoIdStr, 10, 64)

	// 点赞需要用户已经登录
	if _, exist := service.GetExistUserByToken(token); exist {
		// 根据token调用service层的获取这个用户点赞的视频列表服务
		favoriteVideos, err := service.GetFavoriteVideosByToken(token)
		// 如果服务出错，返回服务器错误响应
		if err != nil {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: global.Response{StatusCode: 1, StatusMsg: "Internal Server Error! Get favorite videos failed"},
			})
			return
		}
		// 获取用户是否已经给这个视频点过赞
		_, isFavorite := favoriteVideos[videoId]

		// 如果是点赞行为
		if actionType == "1" {
			// 如果之前已经点过赞了，直接返回
			if isFavorite {
				c.JSON(http.StatusOK, global.Response{StatusCode: 1, StatusMsg: "The video has been liked"})
				return
			}

			// 根据点赞用户的token和点赞视频的Id调用service层的点赞服务，如果出错，返回错误信息响应
			if errMsg := service.LikeAction(token, videoId); errMsg != "" {
				c.JSON(http.StatusOK, global.Response{StatusCode: 1, StatusMsg: errMsg})
				return
			}
		} else if actionType == "2" {
			// 如果是取消点赞行为且之前没有点过赞，直接返回
			if !isFavorite {
				c.JSON(http.StatusOK, global.Response{StatusCode: 1, StatusMsg: "Didn't like the video before"})
				return
			}

			// 根据用户的token和视频的Id调用service层的取消点赞服务，如果出错，返回错误信息响应
			if errMsg := service.CancelLikeAction(token, videoId); errMsg != "" {
				c.JSON(http.StatusOK, global.Response{StatusCode: 1, StatusMsg: errMsg})
				return
			}
		}
		// 调用service层没出错则返回响应成功
		c.JSON(http.StatusOK, global.Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, global.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// 返回点赞喜欢列表
func FavoriteList(c *gin.Context) {
	userIdStr := c.Query("user_id")
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)

	// 根据用户的Id调用service层的获取喜欢视频列表服务，如果出错，返回服务器错误响应
	videoList, err := service.GetFavoriteVideosById(userId)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: global.Response{StatusCode: 1, StatusMsg: "Internal Server Error! Get favorite videos failed"},
		})
		return
	}

	// 调用service层没出错则返回响应成功和喜欢视频列表
	c.JSON(http.StatusOK, VideoListResponse{
		Response: global.Response{
			StatusCode: 0,
		},
		VideoList: videoList,
	})
}

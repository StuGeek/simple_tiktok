package controller

import (
	"net/http"
	"strconv"

	"simple_tiktok/global"
	"simple_tiktok/service"

	"github.com/gin-gonic/gin"
)

type CommentListResponse struct {
	global.Response
	CommentList []global.Comment `json:"comment_list"`
}

type CommentActionResponse struct {
	global.Response
	Comment global.Comment `json:"comment"`
}

// 评论或取消评论行为
func CommentAction(c *gin.Context) {
	token := c.Query("token")
	actionType := c.Query("action_type")
	videoIdStr := c.Query("video_id")
	videoId, err := strconv.ParseInt(videoIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, global.Response{StatusCode: 1, StatusMsg: "Get videoId error"})
		return
	}

	// 如果是发布评论
	if actionType == "1" {
		text := c.Query("comment_text")
		// 调用service层的发布评论服务
		comment, errMsg := service.PublishComment(token, videoId, text)
		// 如果服务出错，返回错误信息响应
		if errMsg != "" {
			c.JSON(http.StatusOK, global.Response{StatusCode: 1, StatusMsg: errMsg})
			return
		}

		// 服务没出错则返回正常响应
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: global.Response{StatusCode: 0},
			Comment:  *comment,
		})
		return
	} else if actionType == "2" {
		commentIdStr := c.Query("comment_id")
		commentId, _ := strconv.ParseInt(commentIdStr, 10, 64)
		// 如果是取消评论行为，调用service层的取消评论服务
		errMsg := service.CancelComment(token, commentId, videoId)
		// 如果服务出错，返回错误信息响应
		if errMsg != "" {
			c.JSON(http.StatusOK, global.Response{StatusCode: 1, StatusMsg: errMsg})
			return
		}
	}
	// 服务没出错则返回正常响应
	c.JSON(http.StatusOK, global.Response{StatusCode: 0})
}

// 查看视频的所有评论，按发布时间倒序
func CommentList(c *gin.Context) {
	videoIdStr := c.Query("video_id")
	videoId, err := strconv.ParseInt(videoIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, global.Response{StatusCode: 1, StatusMsg: "Get videoId error"})
		return
	}

	// 根据视频Id调用service层的获取视频所有评论服务
	commentList, errMsg := service.GetCommentList(videoId)
	// 如果服务出错，返回服务器错误响应
	if errMsg != "" {
		c.JSON(http.StatusOK, CommentListResponse{
			Response: global.Response{StatusCode: 1, StatusMsg: errMsg},
		})
		return
	}

	c.JSON(http.StatusOK, CommentListResponse{
		Response:    global.Response{StatusCode: 0},
		CommentList: commentList,
	})
}

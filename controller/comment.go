package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list"`
}

type CommentActionResponse struct {
	Response
	Comment Comment `json:"comment"`
}

// 评论或取消评论行为
func CommentAction(c *gin.Context) {
	token := c.Query("token")
	actionType := c.Query("action_type")

	// userId := c.Query("user_id")
	videoIdStr := c.Query("video_id")

	videoId, _ := strconv.ParseInt(videoIdStr, 10, 64)

	// 评论需要用户已经登录
	if user, exist := usersLoginInfo[token]; exist {
		// 如果是发布评论
		if actionType == "1" {
			text := c.Query("comment_text")

			now := time.Now()
			createDate := now.Format("01-02")
			publishTime := now.Unix()

			// 向评论信息表中插入相应的评论记录
			newCommentDao := CommentDao{
				UserId:      usersLoginInfo[token].Id,
				VideoId:     videoId,
				Content:     text,
				CreateDate:  createDate,
				PublishTime: publishTime,
			}

			var video VideoDao

			dbMutex.Lock()
			globalDb.Create(&newCommentDao)
			// 更新视频信息表中相应视频的评论数加一
			globalDb.Where("id = ?", videoId).First(&video).Update("comment_count", video.CommentCount+1)
			dbMutex.Unlock()

			c.JSON(http.StatusOK, CommentActionResponse{Response: Response{StatusCode: 0},
				Comment: Comment{
					Id:         newCommentDao.Id,
					User:       user,
					Content:    text,
					CreateDate: createDate,
				}})
			return
		} else if actionType == "2" {
			// commentIdStr := c.Query("comment_id")

			// commentId, _ := strconv.ParseInt(commentIdStr, 10, 64)
			// globalDb.Where("id = ?", commentId).Delete(&CommentDao{})

			var video VideoDao

			dbMutex.Lock()
			// 如果是取消评论，则从评论信息表中删除相应的记录，并更新视频信息表中相应视频的评论数减一
			globalDb.Where("user_id = ? and video_id = ?", usersLoginInfo[token].Id, videoId).Delete(&CommentDao{})
			globalDb.Where("id = ?", videoId).First(&video).Update("comment_count", video.CommentCount-1)
			dbMutex.Unlock()
		}
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// 查看视频的所有评论，按发布时间倒序
func CommentList(c *gin.Context) {
	videoIdStr := c.Query("video_id")

	videoId, _ := strconv.ParseInt(videoIdStr, 10, 64)

	// 从评论信息表中根据视频id获取按发布时间倒序的所有评论
	var comments []CommentDao
	dbMutex.Lock()
	globalDb.Where("video_id = ?", videoId).Order("publish_time desc").Find(&comments)
	dbMutex.Unlock()

	// 从记录账号信息的map中获取发布评论的User信息，返回相应的评论列表
	var commentList []Comment
	for _, commentDao := range comments {
		commentList = append(commentList, Comment{
			Id:         commentDao.Id,
			User:       usersLoginInfo[userIdToToken[commentDao.UserId]],
			Content:    commentDao.Content,
			CreateDate: commentDao.CreateDate,
		})
	}

	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: commentList,
	})
}

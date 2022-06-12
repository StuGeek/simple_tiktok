package service

import (
	"simple_tiktok/global"
	"simple_tiktok/repository"
	"sync"
	"time"
)

// 根据发布评论的用户Id，视频Id，评论内容发布评论，返回发布的评论Id，创建日期，和可能的错误信息
func PublishComment(userId int64, videoId int64, content string) (int64, string, string) {
	var commentCount int64
	var err error = nil
	commentCount, err = repository.QueryCommentCountByVideoId(videoId)
	if err != nil {
		return 0, "", "Internal Server Error! Query comment count failed"
	}

	// 评论数不能超过单个视频的评论最大值
	if global.MaxCommentCount < 0 && commentCount >= global.MaxCommentCount {
		return 0, "", "The number of comments of the video has reached the maximum"
	}

	// 获取当前日期，当前时间
	now := time.Now()
	createDate := now.Format("01-02")
	publishTime := now.Unix()

	var newCommentId int64
	var errMsg string = ""

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		// 向评论信息表中插入相应的评论记录，如果插入失败，记录错误信息
		newCommentId, err = repository.CreateComment(userId, videoId, content, createDate, publishTime)
		if err != nil {
			errMsg = "Internal Server Error! Create comment failed"
		}

	}()
	go func() {
		defer wg.Done()
		// 更新视频信息表中相应视频的评论数加一，如果更新失败，记录错误信息
		err = repository.AddOneVideoCommentCountById(videoId)
		if err != nil {
			errMsg = "Internal Server Error! Add one video comment count failed"
		}
	}()
	wg.Wait()

	// 如果有错误信息，说明出错，返回错误信息
	if errMsg != "" {
		return 0, "", errMsg
	}

	// 没有错误信息则返回新发布的评论Id和发布日期
	return newCommentId, createDate, ""
}

// 根据评论的用户Id，视频Id删除评论，返回可能的错误信息
func CancelComment(commentId int64, videoId int64) string {
	var errMsg string = ""

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		// 如果是取消评论，则从评论信息表中删除相应的记录，如果删除失败，记录错误信息
		err := repository.DeleteComment(commentId)
		if err != nil {
			errMsg = "Internal Server Error! Delete comment failed"
		}
	}()
	go func() {
		defer wg.Done()
		// 更新视频信息表中相应视频的评论数减一，如果更新失败，记录错误信息
		err := repository.SubOneVideoCommentCountById(videoId)
		if err != nil {
			errMsg = "Internal Server Error! Sub one video comment count failed"
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

// 根据视频Id获取视频的评论列表，返回评论列表和可能的错误
func GetCommentList(videoId int64) ([]global.Comment, error) {
	// 从评论信息表中根据视频id获取按发布时间倒序的所有评论
	comments, err := repository.QueryCommentByVideoId(videoId)
	if err != nil {
		return []global.Comment{}, err
	}

	// 从记录账号信息的map中获取发布评论的User信息，返回相应的评论列表
	var commentList []global.Comment
	for _, commentDao := range comments {
		user, _ := repository.GetUserById(commentDao.UserId)
		commentList = append(commentList, global.Comment{
			Id:         commentDao.Id,
			User:       user,
			Content:    commentDao.Content,
			CreateDate: commentDao.CreateDate,
		})
	}

	return commentList, nil
}

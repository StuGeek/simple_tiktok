package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// 评论信息表comments
type CommentDao struct {
	Id          int64  `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	UserId      int64  `json:"user_id"`
	VideoId     int64  `json:"video_id" gorm:"index"`
	Content     string `json:"content"`
	CreateDate  string `json:"create_date"`
	PublishTime int64  `json:"publish_time" gorm:"index:,sort:desc"`
}

func (CommentDao) TableName() string {
	return "comments"
}

// 评论表内连接用户表，多了一个用户名表项，用来获取评论列表返回结果
type CommentJoinUser struct {
	CommentId  int64  `json:"comment_id"`
	UserId     int64  `json:"user_id"`
	Content    string `json:"content"`
	CreateDate string `json:"create_date"`
	UserName   string `json:"user_name"`
}

// 根据评论Id从comments表中找到评论者的Id
func QueryUserIdByCommentId(commentId int64) (int64, error) {
	var comment CommentDao
	err := globalDB.Where("Id = ?", commentId).First(&comment).Error
	if err != nil {
		// 如果没找到评论就返回0和nil
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("comment record not found!", err)
			return 0, nil
		} else {
			fmt.Println("QueryUserIdByCommentId(int64) failed!", err)
			return 0, err
		}
	}

	return comment.UserId, nil
}

// 根据视频Id从comments表中返回视频的所有评论和对应的用，并按评论时间倒序排列
func QueryAllCommentByVideoId(videoId int64) ([]CommentJoinUser, error) {
	// 用comments表和users表查询出特定视频Id的评论和对应的评论用户
	var commentJoinUser []CommentJoinUser
	usersMutex.Lock()
	commentsMutex.Lock()
	err := globalDB.Model(&CommentDao{}).Where("video_id = ?", videoId).Select("comments.id as comment_id, comments.user_id, comments.content, comments.create_date, users.name as user_name").Joins("inner join users on users.id = comments.user_id").Order("publish_time desc").Scan(&commentJoinUser).Error
	commentsMutex.Unlock()
	usersMutex.Unlock()
	if err != nil {
		// 如果没找到评论就返回空评论列表和nil
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("comment record not found!", err)
			return nil, nil
		} else {
			fmt.Println("QueryAllCommentByVideoId(int64) failed!", err)
			return nil, err
		}
	}

	return commentJoinUser, nil
}

// 向comments表中插入一条评论，返回插入的评论Id和可能的错误
func CreateComment(userId int64, videoId int64, content string, createDate string, publishTime int64) (int64, error) {
	newCommentDao := CommentDao{
		UserId:      userId,
		VideoId:     videoId,
		Content:     content,
		CreateDate:  createDate,
		PublishTime: publishTime,
	}

	commentsMutex.Lock()
	err := globalDB.Create(&newCommentDao).Error
	commentsMutex.Unlock()
	if err != nil {
		fmt.Println("Create comment failed!", err)
		return 0, err
	}

	// 返回插入的评论的Id
	newCommentId := newCommentDao.Id
	return newCommentId, nil
}

// 根据评论Id从comments表中删除一条评论，返回可能的错误
func DeleteComment(commentId int64) error {
	commentsMutex.Lock()
	err := globalDB.Where("id = ?", commentId).Delete(&CommentDao{}).Error
	commentsMutex.Unlock()
	if err != nil {
		fmt.Println("Delete comment failed!", err)
		return err
	}

	return nil
}

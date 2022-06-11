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

// 根据视频Id从comments表中返回视频的所有评论，并按评论时间倒序排列
func QueryCommentByVideoId(videoId int64) ([]CommentDao, error) {
	var comments []CommentDao
	err := GlobalDB.Where("video_id = ?", videoId).Order("publish_time desc").Find(&comments).Error
	if err != nil {
		// 如果没找到评论就返回空评论列表和nil
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("comment record not found!", err)
			return []CommentDao{}, nil
		} else {
			fmt.Println("QueryCommentByVideoId(int64) failed!", err)
			return []CommentDao{}, err
		}
	}

	return comments, nil
}

// 向comments表中插入一条评论
func CreateComment(userId int64, videoId int64, content string, createDate string, publishTime int64) (int64, error) {
	newCommentDao := CommentDao{
		UserId:      userId,
		VideoId:     videoId,
		Content:     content,
		CreateDate:  createDate,
		PublishTime: publishTime,
	}

	err := GlobalDB.Create(&newCommentDao).Error
	if err != nil {
		fmt.Println("Create comment failed!", err)
		return 0, err
	}

	// 返回插入的评论的Id
	newCommentId := newCommentDao.Id
	return newCommentId, nil
}

// 根据评论Id从comments表中删除一条评论
func DeleteComment(commentId int64) error {
	err := GlobalDB.Where("id = ?", commentId).Delete(&CommentDao{}).Error
	if err != nil {
		fmt.Println("Delete comment failed!", err)
		return err
	}

	return nil
}

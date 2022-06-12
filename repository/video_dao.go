package repository

import (
	"errors"
	"fmt"
	"simple_tiktok/global"
	"time"

	"gorm.io/gorm"
)

// 视频信息表videos
type VideoDao struct {
	Id            int64  `json:"id,omitempty" gorm:"primary_key;AUTO_INCREMENT"`
	AuthorId      int64  `json:"author_id,omitempty" gorm:"index"`
	PlayUrl       string `json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty"`
	Title         string `json:"title,omitempty"`
	PublishTime   int64  `json:"publish_time,omitempty" gorm:"index:,sort:desc"`
}

func (VideoDao) TableName() string {
	return "videos"
}

// 从videos表中返回所有投稿时间不大于最新投稿时间戳的视频，并按投稿时间倒序排列，最多30个
func QueryVideoByPublishTime(latestTime int64) ([]VideoDao, error) {
	var videoList []VideoDao
	err := GlobalDB.Where("publish_time <= ?", latestTime).Order("publish_time desc").Limit(30).Find(&videoList).Error
	if err != nil {
		// 如果没找到视频就返回空视频列表和nil
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("video record not found!", err)
			return []VideoDao{}, nil
		} else {
			fmt.Println("QueryVideoByPublishTime(int64) failed!", err)
			return []VideoDao{}, err
		}
	}

	return videoList, nil
}

// 根据作者Id从videos表中这个用户发布的所有视频
func QueryVideoByAuthorId(authorId int64) ([]VideoDao, error) {
	// 从数据库中根据用户id获取这个用户发布的视频列表
	var videoList []VideoDao
	err := GlobalDB.Where("author_id = ?", authorId).Find(&videoList).Error
	if err != nil {
		// 如果没找到视频就返回空视频列表和nil
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("video record not found!", err)
			return []VideoDao{}, nil
		} else {
			fmt.Println("QueryVideoByAuthorId(int64) failed!", err)
			return []VideoDao{}, err
		}
	}

	return videoList, nil
}

// 根据视频Id查询视频的评论数
func QueryCommentCountByVideoId(videoId int64) (int64, error) {
	var video VideoDao
	err := GlobalDB.Where("id = ?", videoId).First(&video).Error
	if err != nil {
		// 如果没找到视频就返回0和nil
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("video record not found!", err)
			return 0, nil
		} else {
			fmt.Println("QueryCommentCountByVideoId(int64) failed!", err)
			return 0, err
		}
	}

	return video.CommentCount, nil
}

// 向videos表中插入一条视频
func CreateVideo(authorId int64, finalName string, title string) {
	// 将视频信息存入数据库中，投稿时间为当前时间
	GlobalDB.Create(&VideoDao{
		AuthorId:      authorId,
		PlayUrl:       global.ServerUrl + "static/" + finalName,
		FavoriteCount: 0,
		CommentCount:  0,
		IsFavorite:    false,
		Title:         title,
		PublishTime:   time.Now().Unix(),
	})
}

// 根据视频Id给这个视频的点赞数加一
func AddOneVideoFavoriteCountById(videoId int64) error {
	var video VideoDao
	err := GlobalDB.Where("id = ?", videoId).First(&video).Update("favorite_count", video.FavoriteCount+1).Update("is_favorite", true).Error
	if err != nil {
		fmt.Println("AddOneFavoriteCountById failed", err)
		return err
	}

	return nil
}

// 根据视频Id给这个视频的点赞数减一
func SubOneVideoFavoriteCountById(videoId int64) error {
	var video VideoDao
	err := GlobalDB.Where("id = ?", videoId).First(&video).Update("favorite_count", video.FavoriteCount-1).Update("is_favorite", false).Error
	if err != nil {
		fmt.Println("SubOneVideoFavoriteCountById failed", err)
		return err
	}

	return nil
}

// 根据视频Id给这个视频的评论数加一
func AddOneVideoCommentCountById(videoId int64) error {
	var video VideoDao
	err := GlobalDB.Where("id = ?", videoId).First(&video).Update("comment_count", video.CommentCount+1).Error
	if err != nil {
		fmt.Println("AddOneVideoCommentCountById failed", err)
		return err
	}

	return nil
}

// 根据视频Id给这个视频的评论数减一
func SubOneVideoCommentCountById(videoId int64) error {
	var video VideoDao
	err := GlobalDB.Where("id = ?", videoId).First(&video).Update("comment_count", video.CommentCount-1).Error
	if err != nil {
		fmt.Println("SubOneVideoCommentCountById failed", err)
		return err
	}

	return nil
}

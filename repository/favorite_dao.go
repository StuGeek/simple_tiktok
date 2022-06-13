package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// 点赞视频信息表favorites
type FavoriteDao struct {
	UserId  int64 `json:"user_id" gorm:"index"`  // 用户的Id
	VideoId int64 `json:"video_id" gorm:"index"` // 用户点赞的视频Id
}

func (FavoriteDao) TableName() string {
	return "favorites"
}

// 根据用户Id和视频Id查询用户是否给这个视频点过赞
func QueryIsFavorite(userId int64, videoId int64) (bool, error) {
	err := globalDB.Where("user_id = ? and video_id = ?", userId, videoId).First(&FavoriteDao{}).Error
	if err != nil {
		// 如果没找到视频就返回空视频列表和nil
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		} else {
			fmt.Println("QueryIsFavorite(int64, int64) failed!", err)
			return false, err
		}
	}

	return true, nil
}

// 根据用户Id查询用户所有点赞的视频列表和对应的作者
func QueryAllFavoriteVideoByUserId(userId int64) ([]VideoJoinUser, error) {
	// 用favorite_videos表、videos表和users表查询出特定Id对应用户所点赞的视频和视频的作者
	var videoJoinUser []VideoJoinUser
	err := globalDB.Model(&FavoriteDao{}).Where("favorites.user_id = ?", userId).Select("videos.id, videos.author_id, videos.play_url, videos.cover_url, videos.favorite_count, videos.comment_count, videos.title, users.name").Joins("inner join videos on videos.id = favorites.video_id").Joins("inner join users on users.id = videos.author_id").Scan(&videoJoinUser).Error
	if err != nil {
		// 如果没找到视频就返回空视频列表和nil
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("favorite videos record not found!", err)
			return nil, nil
		} else {
			fmt.Println("QueryAllFavoriteVideoByUserId(int64) failed!", err)
			return nil, err
		}
	}

	return videoJoinUser, nil
}

// 向favorites表中插入一条点赞记录
func CreateFavorite(userId int64, videoId int64) error {
	err := globalDB.Create(&FavoriteDao{UserId: userId, VideoId: videoId}).Error
	if err != nil {
		fmt.Println("CreateFavorite failed", err)
		return err
	}

	return nil
}

// 从favorites表中删除一条点赞记录
func DeleteFavorite(userId int64, videoId int64) error {
	err := globalDB.Where("user_id = ? and video_id = ?", userId, videoId).Delete(&FavoriteDao{}).Error
	if err != nil {
		fmt.Println("DeleteFavorite failed", err)
		return err
	}

	return nil
}

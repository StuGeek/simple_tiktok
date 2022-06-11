package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// 点赞视频信息表favorites
type FavoriteDao struct {
	Token   string `json:"token" gorm:"index"`    // 用户的token
	VideoId int64  `json:"video_id" gorm:"index"` // 用户点赞的视频Id
}

func (FavoriteDao) TableName() string {
	return "favorites"
}

// 根据用户token查询所有用户点赞的视频
func QueryFavoriteVideosByToken(token string) ([]VideoDao, error) {
	// 用favorite_videos表和videos表查询出特定token对应用户所点赞的视频
	var favoriteVideos []VideoDao
	err := GlobalDB.Joins("inner join favorites on videos.id = favorites.video_id").Where("favorites.token = ?", token).Find(&favoriteVideos).Error
	if err != nil {
		// 如果没找到视频就返回空视频列表和nil
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("favorite videos record not found!", err)
			return []VideoDao{}, nil
		} else {
			fmt.Println("QueryFavoriteVideosByToken(string) failed!", err)
			return []VideoDao{}, err
		}
	}

	return favoriteVideos, nil
}

// 向favorite_videos表中插入一条点赞记录
func CreateFavorite(token string, videoId int64) error {
	err := GlobalDB.Create(&FavoriteDao{Token: token, VideoId: videoId}).Error
	if err != nil {
		fmt.Println("CreateFavorite failed", err)
		return err
	}

	return nil
}

// 从favorite_videos表中删除一条点赞记录
func DeleteFavorite(token string, videoId int64) error {
	err := GlobalDB.Where("token = ? and video_id = ?", token, videoId).Delete(&FavoriteDao{}).Error
	if err != nil {
		fmt.Println("DeleteFavorite failed", err)
		return err
	}

	return nil
}

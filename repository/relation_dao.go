package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// 关系信息表relations
type RelationDao struct {
	UserId   int64 `json:"user_id" gorm:"index"`    // 关注者的用户Id
	ToUserId int64 `json:"to_user_id" gorm:"index"` // 被关注者的用户Id
}

func (RelationDao) TableName() string {
	return "relations"
}

// 根据用户Id从relations表中返回用户所有的关注用户
func QueryAllFollowsById(userId int64) ([]RelationDao, error) {
	var follows []RelationDao
	err := GlobalDB.Where("user_id = ?", userId).Find(&follows).Error
	if err != nil {
		// 如果没找到关注用户就返回空关注列表和nil
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("follows record not found!", err)
			return []RelationDao{}, nil
		} else {
			fmt.Println("QueryAllFollowsById(int64) failed", err)
			return []RelationDao{}, err
		}
	}

	return follows, nil
}

// 根据用户Id从relations表中返回用户所有的粉丝用户
func QueryAllFollowersById(userId int64) ([]RelationDao, error) {
	var followers []RelationDao
	err := GlobalDB.Where("to_user_id = ?", userId).Find(&followers).Error
	if err != nil {
		// 如果没找到粉丝用户就返回空粉丝列表和nil
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("follows record not found!", err)
			return []RelationDao{}, nil
		} else {
			fmt.Println("QueryAllFollowsById(int64) failed", err)
			return []RelationDao{}, err
		}
	}

	return followers, nil
}

// 向relations表中插入一条关系记录
func CreateRelation(userId int64, toUserId int64) error {
	err := GlobalDB.Create(&RelationDao{UserId: userId, ToUserId: toUserId}).Error
	if err != nil {
		fmt.Println("CreateRelation failed", err)
		return err
	}

	return nil
}

// 从relations表中删除一条关系记录
func DeleteRelation(userId int64, toUserId int64) error {
	err := GlobalDB.Where("user_id = ? and to_user_id = ?", userId, toUserId).Delete(&RelationDao{}).Error
	if err != nil {
		fmt.Println("DeleteRelation failed", err)
		return err
	}

	return nil
}

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

// 根据关注用户Id和被关注Id查询是否已经关注过用户
func QueryIsFollow(userId int64, toUserId int64) (bool, error) {
	relationsMutex.Lock()
	err := globalDB.Where("user_id = ? and to_user_id = ?", userId, toUserId).First(&RelationDao{}).Error
	relationsMutex.Unlock()
	if err != nil {
		// 如果没找到记录就返回false和nil
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		} else {
			fmt.Println("QueryIsFollow(int64, int64) failed!", err)
			return false, err
		}
	}

	return true, nil
}

// 根据用户Id查找这个用户的关注列表
func QueryFollowUserListById(userId int64) ([]UserDao, error) {
	var followUserList []UserDao
	usersMutex.Lock()
	relationsMutex.Lock()
	err := globalDB.Joins("inner join relations on ? = relations.user_id", userId).Where("id = relations.to_user_id").Find(&followUserList).Error
	relationsMutex.Unlock()
	usersMutex.Unlock()
	if err != nil {
		// 如果没找到关注用户就返回空关注列表和nil
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("Follow user record not found!", err)
			return nil, nil
		} else {
			fmt.Println("QueryFollowUserListById(int64) failed!", err)
			return nil, err
		}
	}

	return followUserList, nil
}

// 根据用户Id查找这个用户的粉丝列表
func QueryFollowerUserListById(userId int64) ([]UserDao, error) {
	var followerUserList []UserDao
	usersMutex.Lock()
	relationsMutex.Lock()
	err := globalDB.Joins("inner join relations on ? = relations.to_user_id", userId).Where("id = relations.user_id").Find(&followerUserList).Error
	relationsMutex.Unlock()
	usersMutex.Unlock()
	if err != nil {
		// 如果没找到粉丝用户就返回空粉丝列表和nil
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("Follower user record not found!", err)
			return nil, nil
		} else {
			fmt.Println("QueryFollowerUserListById(int64) failed!", err)
			return nil, err
		}
	}

	return followerUserList, nil
}

// 向relations表中插入一条关系记录
func CreateRelation(userId int64, toUserId int64) error {
	relationsMutex.Lock()
	err := globalDB.Create(&RelationDao{UserId: userId, ToUserId: toUserId}).Error
	relationsMutex.Unlock()
	if err != nil {
		fmt.Println("CreateRelation failed", err)
		return err
	}

	return nil
}

// 从relations表中删除一条关系记录
func DeleteRelation(userId int64, toUserId int64) error {
	relationsMutex.Lock()
	err := globalDB.Where("user_id = ? and to_user_id = ?", userId, toUserId).Delete(&RelationDao{}).Error
	relationsMutex.Unlock()
	if err != nil {
		fmt.Println("DeleteRelation failed", err)
		return err
	}

	return nil
}

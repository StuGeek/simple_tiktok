package repository

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// 用户信息表users
type UserDao struct {
	Id                int64  `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Name              string `json:"name"`
	FollowCount       int64  `json:"follow_count" gorm:"default:0"`
	FollowerCount     int64  `json:"follower_count" gorm:"default:0"`
	IsFollow          bool   `json:"is_follow" gorm:"default:false"`
	Password          string `json:"password"`
	Token             string `json:"token"`
	TokenLastUsedTime int64  `json:"token_last_used_time"`
}

func (UserDao) TableName() string {
	return "users"
}

// 通过用户Id获取用户的UserDao结构体
func QueryUserById(userId int64) (UserDao, error) {
	var userDao UserDao
	err := globalDB.Where("id = ?", userId).First(&userDao).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("user record not found!", err)
			return UserDao{}, nil
		} else {
			fmt.Println("QueryUserById(int64) failed!", err)
			return UserDao{}, err
		}
	}

	return userDao, nil
}

// 通过用户名获取用户的UserDao结构体
func QueryUserByName(name string) (UserDao, bool, error) {
	var userDao UserDao
	err := globalDB.Where("name = ?", name).First(&userDao).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("user record not found!", err)
			return UserDao{}, false, nil
		} else {
			fmt.Println("QueryUserByName(string) failed!", err)
			return UserDao{}, false, err
		}
	}

	return userDao, true, nil
}

// 通过用户token获取用户的UserDao结构体
func QueryUserByToken(token string) (UserDao, bool, error) {
	var userDao UserDao
	err := globalDB.Where("token = ?", token).First(&userDao).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// fmt.Println("user record not found!", err)
			return UserDao{}, false, nil
		} else {
			fmt.Println("QueryUserByToken(string) failed!", err)
			return UserDao{}, false, err
		}
	}

	return userDao, true, nil
}

// 根据用户名、密码和token创建用户，并返回创建的用户Id
func CreateUser(username string, password string, token string) (int64, string) {
	newUserDao := UserDao{
		Name:              username,
		Password:          password,
		Token:             token,
		TokenLastUsedTime: time.Now().Unix(),
	}
	if err := globalDB.Create(&newUserDao).Error; err != nil {
		fmt.Println("Create user failed!", err)
		return 0, "Create user failed!"
	}

	newUserId := newUserDao.Id
	return newUserId, ""
}

// 更新token的上次使用时间
func UpdataTokenLastUsedTime(token string) error {
	err := globalDB.Where("token = ?", token).First(&UserDao{}).Update("token_last_used_time", time.Now().Unix()).Error
	if err != nil {
		fmt.Println("UpdataTokenLastUsedTime(string) failed", err)
		return err
	}

	return nil
}

// 根据用户Id给这个用户的关注数加一
func AddOneFollowCountById(userId int64) error {
	var user UserDao
	err := globalDB.Where("id = ?", userId).First(&user).Update("follow_count", user.FollowCount+1).Error
	if err != nil {
		fmt.Println("AddOneFollowCountById(int64) failed", err)
		return err
	}

	return nil
}

// 根据用户Id给这个用户的粉丝数加一
func AddOneFollowerCountById(userId int64) error {
	var user UserDao
	err := globalDB.Where("id = ?", userId).First(&user).Update("follower_count", user.FollowerCount+1).Error
	if err != nil {
		fmt.Println("AddOneFollowerCountById(int64) failed", err)
		return err
	}

	return nil
}

// 根据用户Id给这个用户的关注数减一
func SubOneFollowCountById(userId int64) error {
	var user UserDao
	err := globalDB.Where("id = ?", userId).First(&user).Update("follow_count", user.FollowCount-1).Error
	if err != nil {
		fmt.Println("SubOneFollowCountById(int64) failed", err)
		return err
	}

	return nil
}

// 根据用户Id给这个用户的粉丝数减一
func SubOneFollowerCountById(userId int64) error {
	var user UserDao
	err := globalDB.Where("id = ?", userId).First(&user).Update("follower_count", user.FollowerCount-1).Error
	if err != nil {
		fmt.Println("SubOneFollowerCountById(int64) failed", err)
		return err
	}

	return nil
}

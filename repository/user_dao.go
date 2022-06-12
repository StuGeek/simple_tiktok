package repository

import (
	"errors"
	"fmt"
	"simple_tiktok/global"

	"gorm.io/gorm"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
var usersLoginInfo = map[string]global.User{} // 存储用户token与用户User结构体的对应关系
var userIdToToken = map[int64]string{}        // 存储用户Id与用户token的对应关系

// 用户信息表users
type UserDao struct {
	Id            int64  `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Name          string `json:"name"`
	FollowCount   int64  `json:"follow_count" gorm:"default:0"`
	FollowerCount int64  `json:"follower_count" gorm:"default:0"`
	IsFollow      bool   `json:"is_follow" gorm:"default:false"`
	Token         string `json:"token"`
}

func (UserDao) TableName() string {
	return "users"
}

// 查询所有用户信息
func QueryAllUser() ([]UserDao, error) {
	var users []UserDao
	if err := GlobalDB.Find(&users).Error; err != nil {
		// 如果没找到用户就返回空用户列表和nil
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("users record not found!", err)
			return []UserDao{}, nil
		} else {
			fmt.Println("QueryAllUser() failed!", err)
			return []UserDao{}, err
		}
	}

	return users, nil
}

// 根据用户Id查询用户的关注数
func QueryFollowCountByUserId(userId int64) (int64, error) {
	var user UserDao
	err := GlobalDB.Where("id = ?", userId).First(&user).Error
	if err != nil {
		// 如果没找到用户就返回0和nil
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("user record not found!", err)
			return 0, nil
		} else {
			fmt.Println("QueryFollowCountByUserId(int64) failed!", err)
			return 0, err
		}
	}

	return user.FollowCount, nil
}

// 根据用户名和token创建用户，并返回创建的用户Id
func CreateUser(username string, token string) (int64, string) {
	newUserDao := UserDao{Name: username, Token: token}
	if err := GlobalDB.Create(&newUserDao).Error; err != nil {
		fmt.Println("Create user failed!", err)
		return 0, "Create user failed!"
	}

	newUserId := newUserDao.Id
	// 记录用户token与用户User结构体的对应关系，插入数据库后，表中id为主键，可直接获取作为用户id
	SetUsersLoginInfo(token, &global.User{Id: newUserId, Name: username})
	// 记录用户Id与用户token的对应关系
	SetUserIdToToken(newUserId, token)

	return newUserId, ""
}

// 根据用户Id给这个用户的关注数加一
func AddOneFollowCountById(userId int64) error {
	var user UserDao
	err := GlobalDB.Where("id = ?", userId).First(&user).Update("follow_count", user.FollowCount+1).Error
	if err != nil {
		fmt.Println("AddOneFollowCountById failed", err)
		return err
	}

	userToken, _ := GetUserTokenById(userId)
	// 更新存储账号信息usersLoginInfo的map中相应用户的关注数、被关注数、是否被关注等信息
	SetUsersLoginInfo(userToken, &global.User{
		Id:            user.Id,
		Name:          user.Name,
		FollowCount:   user.FollowCount,
		FollowerCount: user.FollowerCount,
		IsFollow:      user.IsFollow,
	})

	return nil
}

// 根据用户Id给这个用户的粉丝数加一
func AddOneFollowerCountById(userId int64) error {
	var user UserDao
	err := GlobalDB.Where("id = ?", userId).First(&user).Update("follower_count", user.FollowerCount+1).Error
	if err != nil {
		fmt.Println("AddOneFollowerCountById failed", err)
		return err
	}

	userToken, _ := GetUserTokenById(userId)
	// 更新存储账号信息usersLoginInfo的map中相应用户的关注数、被关注数、是否被关注等信息
	SetUsersLoginInfo(userToken, &global.User{
		Id:            user.Id,
		Name:          user.Name,
		FollowCount:   user.FollowCount,
		FollowerCount: user.FollowerCount,
		IsFollow:      true,
	})

	return nil
}

// 根据用户Id给这个用户的关注数减一
func SubOneFollowCountById(userId int64) error {
	var user UserDao
	err := GlobalDB.Where("id = ?", userId).First(&user).Update("follow_count", user.FollowCount-1).Error
	if err != nil {
		fmt.Println("SubOneFollowCountById failed", err)
		return err
	}

	userToken, _ := GetUserTokenById(userId)
	// 更新存储账号信息usersLoginInfo的map中相应用户的关注数、被关注数、是否被关注等信息
	SetUsersLoginInfo(userToken, &global.User{
		Id:            user.Id,
		Name:          user.Name,
		FollowCount:   user.FollowCount,
		FollowerCount: user.FollowerCount,
		IsFollow:      user.IsFollow,
	})

	return nil
}

// 根据用户Id给这个用户的粉丝数减一
func SubOneFollowerCountById(userId int64) error {
	var user UserDao
	err := GlobalDB.Where("id = ?", userId).First(&user).Update("follower_count", user.FollowerCount-1).Error
	if err != nil {
		fmt.Println("SubOneFollowerCountById failed", err)
		return err
	}

	userToken, _ := GetUserTokenById(userId)
	// 更新存储账号信息usersLoginInfo的map中相应用户的关注数、被关注数、是否被关注等信息
	SetUsersLoginInfo(userToken, &global.User{
		Id:            user.Id,
		Name:          user.Name,
		FollowCount:   user.FollowCount,
		FollowerCount: user.FollowerCount,
		IsFollow:      false,
	})

	return nil
}

// 存储用户token与用户User结构体的对应关系
func SetUsersLoginInfo(token string, user *global.User) {
	usersLoginInfo[token] = *user
}

// 存储用户Id与用户token的对应关系
func SetUserIdToToken(userId int64, token string) {
	userIdToToken[userId] = token
}

// 通过用户token获取用户的User结构体
func GetUserByToken(token string) (global.User, bool) {
	user, exist := usersLoginInfo[token]
	if exist {
		return user, true
	}

	return global.User{}, false
}

// 通过用户token获取用户的Id
func GetUserIdByToken(token string) (int64, bool) {
	if _, exist := usersLoginInfo[token]; exist {
		return usersLoginInfo[token].Id, true
	}

	return 0, false
}

// 通过用户Id获取用户的User结构体
func GetUserById(userId int64) (global.User, bool) {
	if _, exist := userIdToToken[userId]; exist {
		return usersLoginInfo[userIdToToken[userId]], true
	}

	return global.User{}, false
}

// 通过用户tokrn获取用户的Id
func GetUserTokenById(userId int64) (string, bool) {
	token, exist := userIdToToken[userId]
	if exist {
		return token, true
	}

	return "", false
}

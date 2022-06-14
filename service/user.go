package service

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"simple_tiktok/global"
	"simple_tiktok/repository"
	"time"
)

// 注册用户，并返回注册的用户Id和登录凭证和可能的错误信息
func RegisterUser(username string, password string) (int64, string, string) {
	// 用户名和密码不能为空
	if len(username) == 0 || len(password) == 0 {
		return 0, "", "The username or password can't not be empty"
	}

	// 用户名和密码最长32个字符
	if len(username) > 32 {
		return 0, "", "The length of username should be less than 32"
	}
	if len(password) > 32 {
		return 0, "", "The length of password should be less than 32"
	}

	_, exist, err := repository.QueryUserByName(username)
	// 如果查询用户失败，返回注册失败
	if err != nil {
		return 0, "", "Query user fail"
	}

	// 如果用户已存在，直接返回注册失败
	if exist {
		return 0, "", "User already exist"
	}

	// 根据用户名和密码获取token
	token := generateToken(username, password)
	// 在数据库中创建用户的记录
	newUserId, errMsg := repository.CreateUser(username, MD5(password), token)
	return newUserId, token, errMsg
}

// 登录用户，返回登录用户Id和登录凭证和可能的错误信息
func LoginUser(username string, password string) (int64, string, string) {
	// 用户名和密码不能为空
	if len(username) == 0 || len(password) == 0 {
		return 0, "", "The username or password can't not be empty"
	}

	user, exist, err := repository.QueryUserByName(username)
	// 如果查询用户失败，返回登录失败
	if err != nil {
		return 0, "", "Query user fail"
	}

	// 如果用户不存在，直接返回登录失败
	if !exist {
		return 0, "", "User doesn't exist"
	}

	// 进行密码检验
	if MD5(password) != user.Password {
		return 0, "", "The password is wrong"
	}

	// 根据用户名和密码获取token
	token := generateToken(username, password)
	// 更新token的上次使用时间为当前时间
	repository.UpdataTokenLastUsedTime(token)
	return user.Id, token, ""
}

// 根据token获取用户，返回用户的User结构体、可能的错误信息和错误
func GetUserByToken(token string) (global.User, string) {
	user, exist, err := repository.QueryUserByToken(token)
	if err != nil {
		return global.User{}, "Internal Server Error! Get user failed"
	}
	// 如果token对应的用户不存在
	if !exist {
		return global.User{}, "User doesn't exist"
	} else {
		// 如果用户存在但是token已经过期
		if time.Now().Unix()-user.TokenLastUsedTime > global.MaxTokenValidTime {
			return global.User{}, "Token has expired. Please login again"
		}
		// token没过期更新token的上次使用时间
		if err := repository.UpdataTokenLastUsedTime(token); err != nil {
			return global.User{}, "Internal Server Error! Updata token's last used time failed"
		}
	}

	return userDaoToUser(&user), ""
}

// 根据用户名和密码生成登录凭证token
func generateToken(username string, password string) string {
	return MD5(username) + MD5(password)
}

// 将UserDao结构体转换为User结构体
func userDaoToUser(userDao *repository.UserDao) global.User {
	return global.User{
		Id:            userDao.Id,
		Name:          userDao.Name,
		FollowCount:   userDao.FollowCount,
		FollowerCount: userDao.FollowerCount,
		IsFollow:      userDao.IsFollow,
	}
}

// 对字符串使用md5算法进行加密，得到新的字符串
func MD5(ori string) string {
	h := md5.New()
	h.Write([]byte(ori))
	return hex.EncodeToString(h.Sum(nil))
}

// 对字符串使用sha256算法进行加密，得到新的字符串
func SHA256(ori string) string {
	h := sha256.New()
	h.Write([]byte(ori))
	return hex.EncodeToString(h.Sum(nil))
}

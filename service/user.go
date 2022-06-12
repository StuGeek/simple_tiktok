package service

import (
	"crypto/sha256"
	"encoding/hex"
	"simple_tiktok/global"
	"simple_tiktok/repository"
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

	// 如果用户已存在，直接返回注册失败
	if exist := repository.IsUsernameExist(username); exist {
		return 0, "", "User already exist"
	}

	// 根据用户名和密码获取token
	token := GenerateToken(username, password)
	// 如果登录凭证已存在，直接返回注册失败，保证登录凭证不重复
	if _, exist := GetExistUserByToken(token); exist {
		return 0, "", "User already exist"
	}

	// 在数据库中创建用户的记录
	newUserId, errMsg := repository.CreateUser(username, token)
	return newUserId, token, errMsg
}

// 初始化存储账号信息的usersLoginInfo和userIdToToken
func InitUserInfo() {
	// 获取所有账号信息
	users, err := repository.QueryAllUser()
	if err != nil {
		panic("failed to init user info")
	}

	// 遍历所有账号
	for _, user := range users {
		// 存储每个账号的token和User的对应关系
		repository.SetUsersLoginInfo(user.Token, &global.User{
			Id:            user.Id,
			Name:          user.Name,
			FollowCount:   user.FollowCount,
			FollowerCount: user.FollowerCount,
			IsFollow:      false,
		})
		// 存储每个账号的Id和Token的对应关系
		repository.SetUserIdToToken(user.Id, user.Token)
		// 设置已存在的用户名
		repository.SetUsernameMap(user.Name)
	}
}

// 根据token获取用户，如果用户存在，返回用户的User结构体和true，不存在则返回空结构体和false
func GetExistUserByToken(token string) (global.User, bool) {
	user, exist := repository.GetUserByToken(token)
	return user, exist
}

// 根据用户名和密码生成登录凭证token
func GenerateToken(username string, password string) string {
	return username + SHA256(password)
}

// 对字符串使用sha256算法进行加密，得到新的字符串
func SHA256(ori string) string {
	h := sha256.New()
	h.Write([]byte(ori))
	return hex.EncodeToString(h.Sum(nil))
}

package service

import (
	"crypto/sha256"
	"encoding/hex"
	"simple_tiktok/global"
	"simple_tiktok/repository"
)

// 注册用户，并返回注册的用户Id
func RegisterUser(username string, token string) (int64, error) {
	newUserId, err := repository.CreateUser(username, token)
	return newUserId, err
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

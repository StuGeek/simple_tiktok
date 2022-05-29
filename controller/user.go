package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
var usersLoginInfo = map[string]User{} // 存储用户token与用户User结构体的对应
var userIdToToken = map[int64]string{} // 存储用户Id与用户token的对应

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

// 注册行为
func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	// 根据用户名和密码获取token
	token := username + password

	// 如果用户已存在，直接返回注册失败
	if _, exist := usersLoginInfo[token]; exist {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	} else {
		// 否则创建新用户信息并插入数据库的users表中
		newUserDao := UserDao{
			Name:          username,
			FollowCount:   0,
			FollowerCount: 0,
			IsFollow:      false,
			Token:         token,
		}
		dbMutex.Lock()
		globalDb.Create(&newUserDao)
		dbMutex.Unlock()

		// 记录用户token与用户User结构体的对应关系，插入数据库后，表中id为主键，可直接获取作为用户id
		usersLoginInfo[token] = User{
			Id:            newUserDao.Id,
			Name:          username,
			FollowCount:   0,
			FollowerCount: 0,
			IsFollow:      false,
		}
		// 记录用户Id与用户token的对应关系
		userIdToToken[newUserDao.Id] = token

		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   newUserDao.Id,
			Token:    username + password,
		})
	}
}

// 登录行为
func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	// 根据用户名和密码获取token
	token := username + password

	// 如果用户存在，从usersLoginInfo中根据token取出用户信息并返回
	if user, exist := usersLoginInfo[token]; exist {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   user.Id,
			Token:    token,
		})
	} else {
		// 找不到token则返回用户不存在
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}

// 用户登录时，获取登录用户的信息
func UserInfo(c *gin.Context) {
	token := c.Query("token")
	userIdStr := c.Query("user_id")

	userId, _ := strconv.ParseInt(userIdStr, 10, 64)

	// 如果用户存在，从usersLoginInfo中根据token取出用户信息并返回
	if user, exist := usersLoginInfo[token]; exist {
		// 初始化账号usersLoginInfo表，更新当前用户的关注状态
		InitUserInfoById(userId)
		// 刷新这个用户获取的视频信息
		InitVideoInfo(time.Now().Unix(), token)

		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User:     user,
		})
	} else {
		// 找不到token则返回用户不存在
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}

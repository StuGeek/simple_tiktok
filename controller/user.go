package controller

import (
	"net/http"

	"simple_tiktok/global"
	"simple_tiktok/service"

	"github.com/gin-gonic/gin"
)

type UserLoginResponse struct {
	global.Response
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}

type UserResponse struct {
	global.Response
	User global.User `json:"user"`
}

// 注册行为
func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	var newUserId int64
	var token string
	var errMsg string

	// 注册失败则返回错误信息
	newUserId, token, errMsg = service.RegisterUser(username, password)
	if errMsg != "" {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: global.Response{StatusCode: 1, StatusMsg: errMsg},
		})
	}

	// 调用service层没出错则返回响应成功，以及返回用户Id和token
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: global.Response{StatusCode: 0},
		UserId:   newUserId,
		Token:    token,
	})
}

// 登录行为
func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	// 根据用户名和密码获取token
	token := service.GenerateToken(username, password)

	// 如果用户存在，从service层取出用户信息并返回
	if user, exist := service.GetExistUserByToken(token); exist {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: global.Response{StatusCode: 0},
			UserId:   user.Id,
			Token:    token,
		})
	} else {
		// 在service层找不到token则返回用户名或密码错误
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: global.Response{StatusCode: 1, StatusMsg: "Username or password is wrong"},
		})
	}
}

// 用户登录时，获取登录用户的信息
func UserInfo(c *gin.Context) {
	token := c.Query("token")

	// 如果用户存在，从service层取出用户信息并返回
	if user, exist := service.GetExistUserByToken(token); exist {
		c.JSON(http.StatusOK, UserResponse{
			Response: global.Response{StatusCode: 0},
			User:     user,
		})
	} else {
		// 在service层找不到token则返回用户不存在
		c.JSON(http.StatusOK, UserResponse{
			Response: global.Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}

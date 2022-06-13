package controller

import (
	"net/http"

	"simple_tiktok/global"
	"simple_tiktok/service"

	"github.com/gin-gonic/gin"
)

type UserRegisterResponse struct {
	global.Response
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}

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

	// 注册失败则返回错误信息
	newUserId, token, errMsg := service.RegisterUser(username, password)
	if errMsg != "" {
		c.JSON(http.StatusOK, UserRegisterResponse{
			Response: global.Response{StatusCode: 1, StatusMsg: errMsg},
		})
		return
	}

	// 调用service层没出错则返回响应成功，以及返回用户Id和token
	c.JSON(http.StatusOK, UserRegisterResponse{
		Response: global.Response{StatusCode: 0},
		UserId:   newUserId,
		Token:    token,
	})
}

// 登录行为
func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	// 如果登录出错，返回错误信息
	userId, token, errMsg := service.LoginUser(username, password)
	if errMsg != "" {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: global.Response{StatusCode: 1, StatusMsg: errMsg},
		})
		return
	}

	// 登录成功则返回成功信息
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: global.Response{StatusCode: 0},
		UserId:   userId,
		Token:    token,
	})
}

// 用户登录时，获取登录用户的信息
func UserInfo(c *gin.Context) {
	token := c.Query("token")

	// 如果获取信息出错或用户不存在，返回错误信息
	user, errMsg := service.GetUserByToken(token)
	if errMsg != "" {
		c.JSON(http.StatusOK, UserResponse{
			Response: global.Response{StatusCode: 1, StatusMsg: errMsg},
		})
		return
	}

	// 没有出错则返回用户信息
	c.JSON(http.StatusOK, UserResponse{
		Response: global.Response{StatusCode: 0},
		User:     user,
	})
}

package controller

import (
	"net/http"
	"strconv"

	"simple_tiktok/global"
	"simple_tiktok/service"

	"github.com/gin-gonic/gin"
)

type UserListResponse struct {
	global.Response
	UserList []global.User `json:"user_list"`
}

// 关注或取消关注行为
func RelationAction(c *gin.Context) {
	token := c.Query("token")

	if _, exist := service.GetExistUserByToken(token); exist {
		toUserIdStr := c.Query("to_user_id")
		actionType := c.Query("action_type")
		toUserId, _ := strconv.ParseInt(toUserIdStr, 10, 64)

		// 如果是关注行为
		if actionType == "1" {
			// 根据关注用户的token和被关注用户的Id调用service层的关注服务，如果出错，返回错误信息响应
			if errMsg := service.FollowAction(token, toUserId); errMsg != "" {
				c.JSON(http.StatusOK, global.Response{StatusCode: 1, StatusMsg: errMsg})
				return
			}
		} else if actionType == "2" {
			// 如果是取消关注行为，根据取消关注用户的token和被取消关注用户的Id调用service层的取消关注服务，如果出错，返回错误信息响应
			if errMsg := service.CancelFollowAction(token, toUserId); errMsg != "" {
				c.JSON(http.StatusOK, global.Response{StatusCode: 1, StatusMsg: errMsg})
				return
			}
		}

		// 调用service层没出错则返回响应成功
		c.JSON(http.StatusOK, global.Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, global.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// 返回关注列表
func FollowList(c *gin.Context) {
	userIdStr := c.Query("user_id")
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)

	// 根据用户的Id调用service层的获取关注列表服务，如果出错，返回服务器错误响应
	userList, err := service.GetFollowUserList(userId)
	if err != nil {
		c.JSON(http.StatusOK, global.Response{StatusCode: 1, StatusMsg: "Internal Server Error! Get follow list failed"})
	}

	// 调用service层没出错则返回响应成功和关注用户列表
	c.JSON(http.StatusOK, UserListResponse{
		Response: global.Response{
			StatusCode: 0,
		},
		UserList: userList,
	})
}

// 返回粉丝列表
func FollowerList(c *gin.Context) {
	userIdStr := c.Query("user_id")
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)

	// 根据用户的Id调用service层的获取粉丝列表服务，如果出错，返回服务器错误响应
	userList, err := service.GetFollowerUserList(userId)
	if err != nil {
		c.JSON(http.StatusOK, global.Response{StatusCode: 1, StatusMsg: "Internal Server Error! Get follower list failed"})
	}

	// 调用service层没出错则返回响应成功和粉丝用户列表
	c.JSON(http.StatusOK, UserListResponse{
		Response: global.Response{
			StatusCode: 0,
		},
		UserList: userList,
	})
}

package controller

import (
	"net/http"
	"strconv"

	"github.com/RaymondCode/simple-demo/repository"
	"github.com/gin-gonic/gin"
)

type UserListResponse struct {
	Response
	UserList []User `json:"user_list"`
}

// 关注或取消关注行为
func RelationAction(c *gin.Context) {
	token := c.Query("token")

	if _, exist := usersLoginInfo[token]; exist {
		// userIdStr := c.Query("user_id")
		toUserIdStr := c.Query("to_user_id")
		actionType := c.Query("action_type")

		// userId, _ := strconv.ParseInt(userIdStr, 10, 64)
		toUserId, _ := strconv.ParseInt(toUserIdStr, 10, 64)

		// 从存储账号信息map中根据token获取登录用户的Id
		userId := usersLoginInfo[token].Id
		var user repository.UserDao
		var toUser repository.UserDao

		var isFollow bool

		// 如果是关注行为
		if actionType == "1" {
			repository.DBMutex.Lock()
			// 在关注信息表中创建相应的记录
			repository.GlobalDB.Create(&repository.FollowDao{
				UserId:   userId,
				ToUserId: toUserId,
			})
			// 在用户信息表中更新关注用户和被关注用户的关注数和被关注数
			repository.GlobalDB.Where("id = ?", userId).First(&user).Update("follow_count", user.FollowCount+1)
			repository.GlobalDB.Where("id = ?", toUserId).First(&toUser).Update("follower_count", toUser.FollowerCount+1)
			repository.DBMutex.Unlock()

			// 设置是关注行为
			isFollow = true

		} else if actionType == "2" {
			repository.DBMutex.Lock()
			// 如果是取消关注行为，从关注信息表中删除相应的记录
			repository.GlobalDB.Where("user_id = ? and to_user_id = ?", userId, toUserId).Delete(&repository.FollowDao{})
			// 在用户信息表中更新关注用户和被关注用户的关注数和被关注数
			repository.GlobalDB.Where("id = ?", userId).First(&user).Update("follow_count", user.FollowCount-1)
			repository.GlobalDB.Where("id = ?", toUserId).First(&toUser).Update("follower_count", toUser.FollowerCount-1)
			repository.DBMutex.Unlock()

			// 设置是取消关注行为
			isFollow = false
		}

		// 更新存储账号信息usersLoginInfo的map中相应用户的关注数、被关注数、是否被关注等信息
		usersLoginInfo[token] = User{
			Id:            user.Id,
			Name:          user.Name,
			FollowCount:   user.FollowCount,
			FollowerCount: user.FollowerCount,
			IsFollow:      user.IsFollow,
		}

		usersLoginInfo[userIdToToken[toUserId]] = User{
			Id:            toUser.Id,
			Name:          toUser.Name,
			FollowCount:   toUser.FollowCount,
			FollowerCount: toUser.FollowerCount,
			IsFollow:      isFollow,
		}

		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// 返回关注列表
func FollowList(c *gin.Context) {
	userIdStr := c.Query("user_id")

	userId, _ := strconv.ParseInt(userIdStr, 10, 64)

	// 获取所有这个用户关注的用户
	var follows []repository.FollowDao
	repository.DBMutex.Lock()
	repository.GlobalDB.Where("user_id = ?", userId).Find(&follows)
	repository.DBMutex.Unlock()

	// 将所有关注的用户加入userList中并最后返回
	var userList []User
	for _, followDao := range follows {
		userList = append(userList, usersLoginInfo[userIdToToken[followDao.ToUserId]])
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: userList,
	})
}

// 返回粉丝列表
func FollowerList(c *gin.Context) {
	userIdStr := c.Query("user_id")

	userId, _ := strconv.ParseInt(userIdStr, 10, 64)

	// 获取所有关注这个用户的用户
	var follows []repository.FollowDao
	repository.DBMutex.Lock()
	repository.GlobalDB.Where("to_user_id = ?", userId).Find(&follows)
	repository.DBMutex.Unlock()

	// 将所有这个用户的粉丝加入userList中并最后返回
	var userList []User
	for _, followDao := range follows {
		userList = append(userList, usersLoginInfo[userIdToToken[followDao.UserId]])
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: userList,
	})
}

// 根据用户的Id获取存储关注的用户Id的map
func GetFollowById(userId int64) map[int64]struct{} {
	// 从follows关注信息表中查询出这个用户关注的所有用户的Id
	var follows []repository.FollowDao
	repository.DBMutex.Lock()
	repository.GlobalDB.Where("user_id = ?", userId).Find(&follows)
	repository.DBMutex.Unlock()

	var followsInfo = make(map[int64]struct{})

	// 存储关注用户的用户Id
	for _, follow := range follows {
		followsInfo[follow.ToUserId] = struct{}{}
	}

	return followsInfo
}

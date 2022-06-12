package service

import (
	"simple_tiktok/global"
	"simple_tiktok/repository"
	"sync"
)

// 关注行为，返回可能的错误信息
func FollowAction(userToken string, toUserId int64) string {
	// 从存储账号信息map中根据token获取登录用户的Id
	userId, _ := repository.GetUserIdByToken(userToken)

	// 用户不能关注或取关自己
	if userId == toUserId {
		return "User can't build relation with himself"
	}

	followCount, err := repository.QueryFollowCountByUserId(userId)
	if err != nil {
		return "Internal Server Error! Query follow count failed"
	}

	// 关注数不能超过单个用户的关注用户最大值
	if followCount >= global.MaxFollowUserCount {
		return "The number of follow count of the user has reached the maximum"
	}

	var errMsg string = ""

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		// 在关系信息表中创建相应的记录，如果创建失败，记录错误信息
		err := repository.CreateRelation(userId, toUserId)
		if err != nil {
			errMsg = "Internal Server Error! Create relation failed"
		}
	}()
	go func() {
		defer wg.Done()
		// 在用户信息表中更新关注用户的关注数，如果更新失败，记录错误信息
		err := repository.AddOneFollowCountById(userId)
		if err != nil {
			errMsg = "Internal Server Error! Add one follow count failed"
		}
	}()
	go func() {
		defer wg.Done()
		// 在用户信息表中更新被关注用户的被关注数，如果更新失败，记录错误信息
		err := repository.AddOneFollowerCountById(toUserId)
		if err != nil {
			errMsg = "Internal Server Error! Add one follower count failed"
		}
	}()
	wg.Wait()

	// 如果有错误信息，说明出错，返回错误信息
	if errMsg != "" {
		return errMsg
	}

	// 没有错误信息则返回空字符串
	return ""
}

// 取消关注行为，返回可能的错误信息
func CancelFollowAction(userToken string, toUserId int64) string {
	// 从存储账号信息map中根据token获取登录用户的Id
	userId, _ := repository.GetUserIdByToken(userToken)

	// 用户不能关注或取关自己
	if userId == toUserId {
		return "User can't build relation with himself"
	}

	var errMsg string = ""

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		// 如果是取消关注行为，从关注信息表中删除相应的记录，如果删除失败，记录错误信息
		err := repository.DeleteRelation(userId, toUserId)
		if err != nil {
			errMsg = "Internal Server Error! Delete relation failed"
		}
	}()
	go func() {
		defer wg.Done()
		// 在用户信息表中更新关注用户的关注数，如果更新失败，记录错误信息
		err := repository.SubOneFollowCountById(userId)
		if err != nil {
			errMsg = "Internal Server Error! Sub one follow count failed"
		}
	}()
	go func() {
		defer wg.Done()
		// 在用户信息表中更新被关注用户的被关注数，如果更新失败，记录错误信息
		err := repository.SubOneFollowerCountById(toUserId)
		if err != nil {
			errMsg = "Internal Server Error! Sub one follower count failed"
		}
	}()
	wg.Wait()

	// 如果有错误信息，说明出错，返回错误信息
	if errMsg != "" {
		return errMsg
	}

	// 没有错误信息则返回空字符串
	return ""
}

// 根据用户Id返回关注列表
func GetFollowUserList(userId int64) ([]global.User, error) {
	// 获取所有这个用户关注的用户
	follows, err := repository.QueryAllFollowsById(userId)
	if err != nil {
		return []global.User{}, err
	}

	// 将所有关注的用户加入userList中，并设置IsFollow属性为true最后返回
	var userList []global.User
	for _, relationDao := range follows {
		followUser, _ := repository.GetUserById(relationDao.ToUserId)
		userList = append(userList, global.User{
			Id:            followUser.Id,
			Name:          followUser.Name,
			FollowCount:   followUser.FollowCount,
			FollowerCount: followUser.FollowerCount,
			IsFollow:      true,
		})
	}

	return userList, nil
}

// 根据用户Id返回粉丝列表
func GetFollowerUserList(userId int64) ([]global.User, error) {
	var followList = make(map[int64]struct{})
	var followerList []repository.RelationDao
	var err error

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		// 获取这个用户关注的所有用户
		var follows []repository.RelationDao
		follows, err = repository.QueryAllFollowsById(userId)
		if err == nil {
			// 存储关注用户的用户Id
			for _, follow := range follows {
				followList[follow.ToUserId] = struct{}{}
			}
		}
	}()
	go func() {
		defer wg.Done()
		// 获取所有关注这个用户的用户
		followerList, err = repository.QueryAllFollowersById(userId)
	}()
	wg.Wait()

	if err != nil {
		return []global.User{}, err
	}

	// 将所有这个用户的粉丝加入userList中并最后返回
	var userList []global.User
	for _, followerDao := range followerList {
		followerUser, _ := repository.GetUserById(followerDao.UserId)
		// 获取这个用户是否关注了他的粉丝
		_, isFollow := followList[followerUser.Id]
		userList = append(userList, global.User{
			Id:            followerUser.Id,
			Name:          followerUser.Name,
			FollowCount:   followerUser.FollowCount,
			FollowerCount: followerUser.FollowerCount,
			IsFollow:      isFollow,
		})
	}

	return userList, nil
}

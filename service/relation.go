package service

import (
	"simple_tiktok/global"
	"simple_tiktok/repository"
	"sync"
)

// 关注行为，返回可能的错误信息
func FollowAction(token string, toUserId int64) string {
	// 用户要处于登录状态，登录凭证需要有效
	user, errMsg := GetUserByToken(token)
	if errMsg != "" {
		return errMsg
	}

	// 用户不能关注或取关自己
	if user.Id == toUserId {
		return "User can't build relation with himself"
	}

	// 关注数不能超过单个用户的关注用户最大值
	if global.MaxFollowUserCount >= 0 && user.FollowCount >= global.MaxFollowUserCount {
		return "The number of follow count of the user has reached the maximum"
	}

	// 获取用户是否已经关注过了
	isFollow, err := repository.QueryIsFollow(user.Id, toUserId)
	if err != nil {
		return "Internal Server Error! Query follow failed"
	}
	// 如果之前已经关注过了，直接返回
	if isFollow {
		return "The user has been followed"
	}

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		// 在关系信息表中创建相应的记录，如果创建失败，记录错误信息
		err := repository.CreateRelation(user.Id, toUserId)
		if err != nil {
			errMsg = "Internal Server Error! Create relation failed"
		}
	}()
	go func() {
		defer wg.Done()
		// 在用户信息表中更新关注用户的关注数，如果更新失败，记录错误信息
		err := repository.AddOneFollowCountById(user.Id)
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
func CancelFollowAction(token string, toUserId int64) string {
	// 用户要处于登录状态，登录凭证需要有效
	user, errMsg := GetUserByToken(token)
	if errMsg != "" {
		return errMsg
	}

	// 用户不能关注或取关自己
	if user.Id == toUserId {
		return "User can't build relation with himself"
	}

	// 获取用户是否已经关注过了
	isFollow, err := repository.QueryIsFollow(user.Id, toUserId)
	if err != nil {
		return "Internal Server Error! Query follow failed"
	}
	// 如果之前没有关注过，直接返回
	if !isFollow {
		return "Didn't follow the user before"
	}

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		// 如果是取消关注行为，从关注信息表中删除相应的记录，如果删除失败，记录错误信息
		err := repository.DeleteRelation(user.Id, toUserId)
		if err != nil {
			errMsg = "Internal Server Error! Delete relation failed"
		}
	}()
	go func() {
		defer wg.Done()
		// 在用户信息表中更新关注用户的关注数，如果更新失败，记录错误信息
		err := repository.SubOneFollowCountById(user.Id)
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

// 根据token和查看用户的Id返回关注列表
func GetFollowUserList(token string, checkUserId int64) ([]global.User, string) {
	var tokenUserFollowMap = make(map[int64]struct{})
	var checkUserFollowList []repository.UserDao
	var errMsg string = ""

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		// 获取登录用户
		user, errMsg := GetUserByToken(token)
		// 查看关注列表可以不处在登录状态，不用返回错误信息
		if errMsg == "" {
			// 如果登录凭证有效，获取登录用户的关注列表
			tokenUserFollowList, err := repository.QueryFollowUserListById(user.Id)
			if err != nil {
				errMsg = "Internal Server Error! Query follow users failed"
			} else {
				// 记录登录用户关注的用户Id
				for _, follow := range tokenUserFollowList {
					tokenUserFollowMap[follow.Id] = struct{}{}
				}
			}

		}
	}()
	go func() {
		defer wg.Done()
		var err error
		// 获取查看的这个用户的关注列表
		checkUserFollowList, err = repository.QueryFollowUserListById(checkUserId)
		if err != nil {
			errMsg = "Internal Server Error! Query follow users failed"
		}
	}()

	wg.Wait()

	if errMsg != "" {
		return nil, errMsg
	}

	// 将所有关注的用户加入userList中，并根据登录用户是否关注设置IsFollow属性最后返回
	var userList []global.User
	for _, followUser := range checkUserFollowList {
		// 获取登录用户是否关注了查看用户的关注者
		_, isFollow := tokenUserFollowMap[followUser.Id]
		userList = append(userList, global.User{
			Id:            followUser.Id,
			Name:          followUser.Name,
			FollowCount:   followUser.FollowCount,
			FollowerCount: followUser.FollowerCount,
			IsFollow:      isFollow,
		})
	}

	return userList, ""
}

// 根据token和查看用户的Id返回粉丝列表
func GetFollowerUserList(token string, checkUserId int64) ([]global.User, string) {
	var tokenUserFollowMap = make(map[int64]struct{})
	var checkUserFollowerList []repository.UserDao
	var errMsg string = ""

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		// 获取登录用户
		user, errMsg := GetUserByToken(token)
		// 查看关注列表可以不处在登录状态，不用返回错误信息
		if errMsg == "" {
			// 如果登录凭证有效，获取登录用户的关注列表
			tokenUserFollowList, err := repository.QueryFollowUserListById(user.Id)
			if err != nil {
				errMsg = "Internal Server Error! Query follow users failed"
			} else {
				// 记录登录用户关注的用户Id
				for _, follow := range tokenUserFollowList {
					tokenUserFollowMap[follow.Id] = struct{}{}
				}
			}

		}
	}()
	go func() {
		defer wg.Done()
		var err error
		// 获取查看的这个用户的粉丝列表
		checkUserFollowerList, err = repository.QueryFollowerUserListById(checkUserId)
		if err != nil {
			errMsg = "Internal Server Error! Query follower users failed"
		}
	}()

	wg.Wait()

	if errMsg != "" {
		return nil, errMsg
	}

	// 将所有这个用户的粉丝加入userList中，并根据登录用户是否关注设置IsFollow属性最后返回
	var userList []global.User
	for _, follower := range checkUserFollowerList {
		// 获取登录用户是否关注了查看用户的粉丝
		_, isFollow := tokenUserFollowMap[follower.Id]
		userList = append(userList, global.User{
			Id:            follower.Id,
			Name:          follower.Name,
			FollowCount:   follower.FollowCount,
			FollowerCount: follower.FollowerCount,
			IsFollow:      isFollow,
		})
	}

	return userList, ""
}

package controller

import (
	"fmt"
	"net/http"
	"simple_tiktok/global"
	"strconv"
	"testing"
	"time"
)

func BenchmarkRegister(b *testing.B) {
	method := "POST"

	client := &http.Client{}
	now := (int)(time.Now().Unix())

	b.ResetTimer()
	// 每次注册的用户名由当前时间加i组成，保证不会重复
	for i := 0; i < b.N; i++ {
		url := "douyin/user/register/?username=" + strconv.Itoa(now+i) + "&password=123456"

		req, err := http.NewRequest(method, global.ServerUrl+url, nil)

		if err != nil {
			fmt.Println(err)
			return
		}

		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer res.Body.Close()
	}
}

func BenchmarkLogin(b *testing.B) {
	url := "douyin/user/login/?username=user1&password=111111"
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, global.ServerUrl+url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	b.ResetTimer()
	// 重复登录user1的账号
	for i := 0; i < b.N; i++ {
		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer res.Body.Close()
	}
}

func BenchmarkUserInfo(b *testing.B) {
	url := "douyin/user/?user_id=1&token=user1bcb15f821479b4d5772bd0ca866c00ad5f926e3580720659cc80d39c9d09802a"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, global.ServerUrl+url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	b.ResetTimer()
	// 获取账号信息
	for i := 0; i < b.N; i++ {
		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer res.Body.Close()
	}
}

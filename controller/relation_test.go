package controller

import (
	"fmt"
	"net/http"
	"testing"
)

func BenchmarkRelationAction(b *testing.B) {
	url := "douyin/relation/action/?token=user1bcb15f821479b4d5772bd0ca866c00ad5f926e3580720659cc80d39c9d09802a&to_user_id=3&action_type="
	method := "POST"

	client := &http.Client{}

	var req *http.Request
	var err error

	n := b.N
	if n == 0 || n == 1 {
		return
	}
	if n%2 == 1 {
		n--
	}

	preUrl := serverUrl + url
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if i%2 == 0 {
			req, err = http.NewRequest(method, preUrl+"1", nil)
		} else {
			req, err = http.NewRequest(method, preUrl+"2", nil)
		}

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

func BenchmarkFollowList(b *testing.B) {
	url := "douyin/relation/follow/list/?user_id=1&token=user1bcb15f821479b4d5772bd0ca866c00ad5f926e3580720659cc80d39c9d09802a"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, serverUrl+url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer res.Body.Close()
	}
}

func BenchmarkFollowerList(b *testing.B) {
	url := "douyin/relation/follower/list/?user_id=1&token=user1bcb15f821479b4d5772bd0ca866c00ad5f926e3580720659cc80d39c9d09802a"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, serverUrl+url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer res.Body.Close()
	}
}

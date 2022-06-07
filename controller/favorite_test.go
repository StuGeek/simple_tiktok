package controller

import (
	"fmt"
	"net/http"
	"testing"
)

func BenchmarkFavoriteAction(b *testing.B) {
	url := "douyin/favorite/action/?token=user368487dc295052aa79c530e283ce698b8c6bb1b42ff0944252e1910dbecdc5425&video_id=1&action_type="
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

func BenchmarkFavoriteList(b *testing.B) {
	url := "douyin/favorite/list/?user_id=1&token=user1bcb15f821479b4d5772bd0ca866c00ad5f926e3580720659cc80d39c9d09802a"
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

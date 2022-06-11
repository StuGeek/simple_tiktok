package controller

import (
	"fmt"
	"net/http"
	"simple_tiktok/global"
	"testing"
)

func BenchmarkCommentAction(b *testing.B) {
	url := "douyin/comment/action/?token=user1bcb15f821479b4d5772bd0ca866c00ad5f926e3580720659cc80d39c9d09802a&video_id=1&action_type="
	method := "POST"

	client := &http.Client{}

	var req *http.Request
	var err error

	// 测试次数必须是偶数
	n := b.N
	if n == 0 || n == 1 {
		return
	}
	if n%2 == 1 {
		n--
	}

	preUrl := global.ServerUrl + url
	b.ResetTimer()
	// 交替评论和删除评论
	for i := 0; i < n; i++ {
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

func BenchmarkCommentList(b *testing.B) {
	url := "douyin/comment/list/?token=user1bcb15f821479b4d5772bd0ca866c00ad5f926e3580720659cc80d39c9d09802a&video_id=1"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, global.ServerUrl+url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	b.ResetTimer()
	// 获取评论列表
	for i := 0; i < b.N; i++ {
		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer res.Body.Close()
	}
}

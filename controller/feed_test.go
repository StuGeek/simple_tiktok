package controller

import (
	"fmt"
	"net/http"
	"simple_tiktok/global"
	"testing"
)

func BenchmarkFeed(b *testing.B) {
	url := "douyin/feed"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, global.ServerUrl+url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	b.ResetTimer()
	// 测试拉取视频流
	for i := 0; i < b.N; i++ {
		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer res.Body.Close()
	}
}

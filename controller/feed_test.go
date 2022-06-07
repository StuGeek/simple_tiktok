package controller

import (
	"fmt"
	"net/http"
	"testing"
)

func BenchmarkFeed(b *testing.B) {
	url := "douyin/feed"
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

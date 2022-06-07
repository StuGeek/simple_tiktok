package controller

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

func BenchmarkPublish(b *testing.B) {
	url := "douyin/publish/action/"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	file, errFile1 := os.Open("../public/bear.mp4")
	if errFile1 != nil {
		fmt.Println(errFile1)
		return
	}
	defer file.Close()

	part1, errFile1 := writer.CreateFormFile("data", filepath.Base("../public/bear.mp4"))
	if errFile1 != nil {
		fmt.Println(errFile1)
		return
	}

	_, errFile1 = io.Copy(part1, file)
	if errFile1 != nil {
		fmt.Println(errFile1)
		return
	}
	_ = writer.WriteField("token", "user1bcb15f821479b4d5772bd0ca866c00ad5f926e3580720659cc80d39c9d09802a")
	_ = writer.WriteField("title", "bear")
	err := writer.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, serverUrl+url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

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

func BenchmarkPublishList(b *testing.B) {
	url := "douyin/publish/list/?token=user1bcb15f821479b4d5772bd0ca866c00ad5f926e3580720659cc80d39c9d09802a&user_id=1"
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

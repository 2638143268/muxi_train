package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func testGetBook() {
	resp, err := http.Get("http://localhost:8081/book?title=三体")
	if err != nil {
		fmt.Println("GET 请求失败:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("GET /book 响应:", string(body))
}

func testPostComment() {
	comment := map[string]string{
		"user":    "小李",
		"comment": "这本书真棒！",
	}

	jsonData, _ := json.Marshal(comment)
	resp, err := http.Post("http://localhost:8081/comment", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("POST 请求失败:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("POST /comment 响应:", string(body))
}

func main() {
	testGetBook()
	testPostComment()
}

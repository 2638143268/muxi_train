package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Comment struct {
	User    string `json:"user"`
	Comment string `json:"comment"`
}

type CommentResponse struct {
	Message string `json:"message"`
	User    string `json:"user"`
	Comment string `json:"comment"`
}

func bookHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	if title == "" {
		http.Error(w, "缺少 title 参数", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "您正在查询图书：《%s》", title)
}

func commentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "只支持 POST 请求", http.StatusMethodNotAllowed)
		return
	}

	var cmt Comment
	err := json.NewDecoder(r.Body).Decode(&cmt)
	if err != nil {
		http.Error(w, "JSON解析失败", http.StatusBadRequest)
		return
	}

	response := CommentResponse{
		Message: "评论提交成功",
		User:    cmt.User,
		Comment: cmt.Comment,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/book", bookHandler)
	http.HandleFunc("/comment", commentHandler)

	fmt.Println("服务器启动，监听端口 :8081...")
	http.ListenAndServe(":8081", nil)
}

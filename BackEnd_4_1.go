package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Book 定义图书结构体
type Book struct {
	ID     string `json:"id"`     // 图书ID
	Title  string `json:"title"`  // 标题
	Author string `json:"author"` // 作者
	Stock  string `json:"stock"`  // 库存
}

// 模拟数据库，用 map 存储图书信息
var books = make(map[string]Book)

// AddBook 添加图书
func AddBook(c *gin.Context) {
	var book Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if _, exists := books[book.ID]; exists {
		c.JSON(http.StatusBadRequest, gin.H{"message": "图书已存在"})
		return
	}
	books[book.ID] = book
	c.JSON(http.StatusOK, gin.H{"message": "添加成功", "book": book})
}

// DeleteBook 删除图书
func DeleteBook(c *gin.Context) {
	id := c.Param("id")
	if _, exists := books[id]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"message": "图书不存在"})
		return
	}
	delete(books, id)
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// UpdateBook 更新图书信息
func UpdateBook(c *gin.Context) {
	var book Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if _, exists := books[book.ID]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"message": "图书不存在"})
		return
	}
	books[book.ID] = book
	c.JSON(http.StatusOK, gin.H{"message": "更新成功", "book": book})
}

// SearchAllBook 查询所有图书
func SearchAllBook(c *gin.Context) {
	var bookList []Book
	for _, book := range books {
		bookList = append(bookList, book)
	}
	c.JSON(http.StatusOK, bookList)
}

// main 主函数
func main() {
	r := gin.Default()

	// 图书管理路由
	r.POST("/book", AddBook)          // 添加图书
	r.DELETE("/book/:id", DeleteBook) // 删除图书
	r.PUT("/book", UpdateBook)        // 更新图书
	r.GET("/books", SearchAllBook)    // 获取所有图书

	// 启动服务
	r.Run(":8080")
}

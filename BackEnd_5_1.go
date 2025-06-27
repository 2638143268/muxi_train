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

// 通用响应结构
type MessageResponse struct {
	Message string `json:"message"`
}

// 图书响应结构
type BookResponse struct {
	Message string `json:"message"`
	Book    Book   `json:"book"`
}

// 模拟数据库
var books = make(map[string]Book)

// AddBook 添加图书
// @Summary 添加图书
// @Description 添加一本新图书
// @Tags 图书管理
// @Accept json
// @Produce json
// @Param book body Book true "图书信息"
// @Success 200 {object} BookResponse
// @Failure 400 {object} MessageResponse
// @Router /book [post]
func AddBook(c *gin.Context) {
	var book Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, MessageResponse{Message: err.Error()})
		return
	}
	if _, exists := books[book.ID]; exists {
		c.JSON(http.StatusBadRequest, MessageResponse{Message: "图书已存在"})
		return
	}
	books[book.ID] = book
	c.JSON(http.StatusOK, BookResponse{Message: "添加成功", Book: book})
}

// DeleteBook 删除图书
// @Summary 删除图书
// @Description 通过 ID 删除图书
// @Tags 图书管理
// @Produce json
// @Param id path string true "图书ID"
// @Success 200 {object} MessageResponse
// @Failure 404 {object} MessageResponse
// @Router /book/{id} [delete]
func DeleteBook(c *gin.Context) {
	id := c.Param("id")
	if _, exists := books[id]; !exists {
		c.JSON(http.StatusNotFound, MessageResponse{Message: "图书不存在"})
		return
	}
	delete(books, id)
	c.JSON(http.StatusOK, MessageResponse{Message: "删除成功"})
}

// UpdateBook 更新图书信息
// @Summary 更新图书
// @Description 更新已有图书的信息
// @Tags 图书管理
// @Accept json
// @Produce json
// @Param book body Book true "图书信息"
// @Success 200 {object} BookResponse
// @Failure 400 {object} MessageResponse
// @Failure 404 {object} MessageResponse
// @Router /book [put]
func UpdateBook(c *gin.Context) {
	var book Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, MessageResponse{Message: err.Error()})
		return
	}
	if _, exists := books[book.ID]; !exists {
		c.JSON(http.StatusNotFound, MessageResponse{Message: "图书不存在"})
		return
	}
	books[book.ID] = book
	c.JSON(http.StatusOK, BookResponse{Message: "更新成功", Book: book})
}

// SearchAllBook 查询所有图书
// @Summary 查询所有图书
// @Description 获取所有图书信息
// @Tags 图书管理
// @Produce json
// @Success 200 {array} Book
// @Router /books [get]
func SearchAllBook(c *gin.Context) {
	var bookList []Book
	for _, book := range books {
		bookList = append(bookList, book)
	}
	c.JSON(http.StatusOK, bookList)
}

// @title 图书管理 API
// @version 1.0
// @description 这是一个使用 Gin 和 Swag 编写的简单图书管理 API。
// @host localhost:8080
// @BasePath /
func main() {
	r := gin.Default()

	// 图书管理路由
	r.POST("/book", AddBook)
	r.DELETE("/book/:id", DeleteBook)
	r.PUT("/book", UpdateBook)
	r.GET("/books", SearchAllBook)

	r.Run(":8080")
}

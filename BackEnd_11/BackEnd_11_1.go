package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Book 定义图书结构体
type Book struct {
	ID     string `json:"id" gorm:"primaryKey"` // 图书ID，设置为主键
	Title  string `json:"title"`                // 标题
	Author string `json:"author"`               // 作者
	Stock  string `json:"stock"`                // 库存
}

var db *gorm.DB

// InitDB 初始化数据库连接
func InitDB() {
	dsn := "root:root@tcp(localhost:3306)/book_db?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 自动迁移表结构
	err = db.AutoMigrate(&Book{})
	if err != nil {
		panic("failed to migrate database")
	}
}

// AddBook 添加图书
func AddBook(c *gin.Context) {
	var book Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查图书是否已存在
	var existingBook Book
	result := db.First(&existingBook, "id = ?", book.ID)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "图书已存在"})
		return
	}

	// 创建新图书
	result = db.Create(&book)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "添加成功", "book": book})
}

// DeleteBook 删除图书
func DeleteBook(c *gin.Context) {
	id := c.Param("id")

	// 检查图书是否存在
	var book Book
	result := db.First(&book, "id = ?", id)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "图书不存在"})
		return
	}

	// 删除图书
	result = db.Delete(&Book{}, "id = ?", id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// UpdateBook 更新图书信息
func UpdateBook(c *gin.Context) {
	var book Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查图书是否存在
	var existingBook Book
	result := db.First(&existingBook, "id = ?", book.ID)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "图书不存在"})
		return
	}

	// 更新图书
	result = db.Save(&book)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功", "book": book})
}

// SearchAllBook 查询所有图书
func SearchAllBook(c *gin.Context) {
	var books []Book
	result := db.Find(&books)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, books)
}

func main() {
	// 初始化数据库连接
	InitDB()

	r := gin.Default()

	// 添加CORS中间件
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 图书管理路由（全部使用单数形式）
	r.POST("/book", AddBook)          // 添加图书
	r.DELETE("/book/:id", DeleteBook) // 删除图书
	r.PUT("/book", UpdateBook)        // 更新图书
	r.GET("/book", SearchAllBook)     // 获取所有图书（改为单数）

	// 启动服务在8082端口
	r.Run(":8082")
}

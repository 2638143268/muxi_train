package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Book 定义图书结构体
type Book struct {
	ID     string `json:"id" gorm:"primaryKey"` // 图书ID，设置为主键
	Title  string `json:"title"`               // 标题
	Author string `json:"author"`              // 作者
	Stock  string `json:"stock"`               // 库存
}

var db *gorm.DB

// InitDB 初始化数据库连接
func InitDB() {
	// 从环境变量获取数据库配置
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":3306)/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	// 自动迁移表结构
	err = db.AutoMigrate(&Book{})
	if err != nil {
		panic("failed to migrate database: " + err.Error())
	}
}

// ... 其他函数保持不变 ...
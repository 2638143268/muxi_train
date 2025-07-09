package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 辅助函数：将字符串转换为uint
func stringToUint(s string) (uint, error) {
	i, err := strconv.ParseUint(s, 10, 64)
	return uint(i), err
}

// 检查是否为管理员
func isAdmin(username, password string) bool {
	// 管理员账号：admin / admin123
	return username == "admin" && password == "admin123"
}

func main() {
	// 使用MySQL数据库
	dsn := "root:123456@tcp(localhost:3306)/repo?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败: " + err.Error())
	}

	// 用户表
	type User struct {
		ID        uint      `gorm:"primaryKey" json:"id"`
		Name      string    `gorm:"size:255" json:"name"`
		Password  string    `gorm:"size:255" json:"-"`
		Username  string    `gorm:"size:255;unique" json:"username"`
		Avatar    string    `gorm:"size:255" json:"avatar"`
		IsAdmin   bool      `gorm:"default:false" json:"is_admin"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	// 帖子表
	type UserRepo struct {
		ID        uint      `gorm:"primaryKey" json:"id"`
		UserID    uint      `gorm:"not null" json:"user_id"`
		RepoName  string    `gorm:"size:255" json:"repo_name"`
		RepoURL   string    `gorm:"size:255" json:"repo_url"`
		User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	// 评论表
	type Comment struct {
		ID        uint      `gorm:"primaryKey" json:"id"`
		UserID    uint      `gorm:"not null" json:"user_id"`
		RepoID    uint      `gorm:"not null" json:"repo_id"`
		Content   string    `gorm:"size:255" json:"content"`
		User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	// 点赞表
	type UserLike struct {
		ID        uint      `gorm:"primaryKey" json:"id"`
		UserID    uint      `gorm:"not null" json:"user_id"`
		RepoID    uint      `gorm:"not null" json:"repo_id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	// 自动迁移数据库表结构
	err = db.AutoMigrate(&User{}, &UserRepo{}, &Comment{}, &UserLike{})
	if err != nil {
		panic("数据库迁移失败: " + err.Error())
	}

	fmt.Println("MySQL数据库连接成功")

	// 创建服务器
	r := gin.Default()

	// 添加CORS中间件
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}
		c.Next()
	})

	// 注册路由
	r.POST("/register", func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
			return
		}

		result := db.Create(&user)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "注册成功", "user": user})
	})

	r.POST("/login", func(c *gin.Context) {
		var loginData struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&loginData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
			return
		}

		// 检查是否为管理员登录
		if isAdmin(loginData.Username, loginData.Password) {
			adminUser := User{
				ID:       0, // 管理员特殊ID
				Name:     "管理员",
				Username: "admin",
				IsAdmin:  true,
			}
			c.JSON(http.StatusOK, gin.H{"message": "管理员登录成功", "user": adminUser})
			return
		}

		var user User
		result := db.Where("username = ? AND password = ?", loginData.Username, loginData.Password).First(&user)
		if result.Error != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "登录失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "登录成功", "user": user})
	})

	r.POST("/repo", func(c *gin.Context) {
		var repoData struct {
			UserID   string `json:"user_id"`
			RepoName string `json:"repo_name"`
			RepoURL  string `json:"repo_url"`
		}
		if err := c.ShouldBindJSON(&repoData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
			return
		}

		userID, err := stringToUint(repoData.UserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
			return
		}

		repo := UserRepo{
			UserID:   userID,
			RepoName: repoData.RepoName,
			RepoURL:  repoData.RepoURL,
		}

		result := db.Create(&repo)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "发帖失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "发帖成功", "repo": repo})
	})

	r.POST("/comment", func(c *gin.Context) {
		var commentData struct {
			UserID  string `json:"user_id"`
			RepoID  string `json:"repo_id"`
			Content string `json:"content"`
		}
		if err := c.ShouldBindJSON(&commentData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
			return
		}

		userID, err := stringToUint(commentData.UserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
			return
		}

		repoID, err := stringToUint(commentData.RepoID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的帖子ID"})
			return
		}

		comment := Comment{
			UserID:  userID,
			RepoID:  repoID,
			Content: commentData.Content,
		}

		result := db.Create(&comment)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "评论失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "评论成功", "comment": comment})
	})

	r.POST("/like", func(c *gin.Context) {
		var likeData struct {
			UserID string `json:"user_id"`
			RepoID string `json:"repo_id"`
		}
		if err := c.ShouldBindJSON(&likeData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
			return
		}

		userID, err := stringToUint(likeData.UserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
			return
		}

		repoID, err := stringToUint(likeData.RepoID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的帖子ID"})
			return
		}

		// 检查是否已经点赞
		var existingLike UserLike
		result := db.Where("user_id = ? AND repo_id = ?", userID, repoID).First(&existingLike)
		if result.Error == nil {
			// 如果已经点赞，则取消点赞
			db.Delete(&existingLike)
			c.JSON(http.StatusOK, gin.H{"message": "取消点赞成功"})
			return
		}

		like := UserLike{
			UserID: userID,
			RepoID: repoID,
		}

		result = db.Create(&like)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "点赞失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "点赞成功", "like": like})
	})

	// 获取帖子详情
	r.GET("/repos/:id", func(c *gin.Context) {
		id := c.Param("id")

		var repo struct {
			UserRepo
			UserName     string `json:"user_name"`
			LikeCount    int    `json:"like_count"`
			CommentCount int    `json:"comment_count"`
		}

		if err := db.Table("user_repos").
			Select("user_repos.*, users.name as user_name, "+
				"(SELECT COUNT(*) FROM user_likes WHERE user_likes.repo_id = user_repos.id) as like_count, "+
				"(SELECT COUNT(*) FROM comments WHERE comments.repo_id = user_repos.id) as comment_count").
			Joins("left join users on user_repos.user_id = users.id").
			Where("user_repos.id = ?", id).
			First(&repo).Error; err != nil {
			c.JSON(404, gin.H{"error": "帖子不存在"})
			return
		}

		c.JSON(200, repo)
	})

	// 获取帖子评论
	r.GET("/comments", func(c *gin.Context) {
		repoID := c.Query("repo_id")

		var comments []struct {
			Comment
			UserName string `json:"user_name"`
		}

		if err := db.Table("comments").
			Select("comments.*, users.name as user_name").
			Joins("left join users on comments.user_id = users.id").
			Where("comments.repo_id = ?", repoID).
			Order("comments.created_at desc").
			Scan(&comments).Error; err != nil {
			c.JSON(500, gin.H{"error": "获取评论失败"})
			return
		}

		c.JSON(200, comments)
	})

	// 获取帖子列表
	r.GET("/repos", func(c *gin.Context) {
		userID := c.Query("user_id") // 可选参数，用于筛选特定用户的帖子

		var repos []struct {
			UserRepo
			UserName     string `json:"user_name"`
			LikeCount    int    `json:"like_count"`
			CommentCount int    `json:"comment_count"`
		}

		query := db.Table("user_repos").
			Select("user_repos.*, users.name as user_name, " +
				"(SELECT COUNT(*) FROM user_likes WHERE user_likes.repo_id = user_repos.id) as like_count, " +
				"(SELECT COUNT(*) FROM comments WHERE comments.repo_id = user_repos.id) as comment_count").
			Joins("left join users on user_repos.user_id = users.id")

		// 如果指定了用户ID，则只返回该用户的帖子
		if userID != "" {
			query = query.Where("user_repos.user_id = ?", userID)
		}

		if err := query.Order("user_repos.created_at desc").Scan(&repos).Error; err != nil {
			c.JSON(500, gin.H{"error": "获取帖子列表失败"})
			return
		}

		c.JSON(200, repos)
	})

	// 删除帖子（管理员或帖子作者）
	r.DELETE("/repos/:id", func(c *gin.Context) {
		id := c.Param("id")
		var deleteData struct {
			UserID  string `json:"user_id"`
			IsAdmin bool   `json:"is_admin"`
		}
		if err := c.ShouldBindJSON(&deleteData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
			return
		}

		// 获取帖子信息
		var repo UserRepo
		if err := db.First(&repo, id).Error; err != nil {
			c.JSON(404, gin.H{"error": "帖子不存在"})
			return
		}

		// 检查权限：管理员或帖子作者
		if !deleteData.IsAdmin {
			userID, err := stringToUint(deleteData.UserID)
			if err != nil || repo.UserID != userID {
				c.JSON(403, gin.H{"error": "无权限删除此帖子"})
				return
			}
		}

		// 删除相关的点赞和评论
		db.Where("repo_id = ?", id).Delete(&UserLike{})
		db.Where("repo_id = ?", id).Delete(&Comment{})

		// 删除帖子
		if err := db.Delete(&repo).Error; err != nil {
			c.JSON(500, gin.H{"error": "删除帖子失败"})
			return
		}

		c.JSON(200, gin.H{"message": "删除帖子成功"})
	})

	// 删除评论（管理员或评论作者）
	r.DELETE("/comments/:id", func(c *gin.Context) {
		id := c.Param("id")
		var deleteData struct {
			UserID  string `json:"user_id"`
			IsAdmin bool   `json:"is_admin"`
		}
		if err := c.ShouldBindJSON(&deleteData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
			return
		}

		// 获取评论信息
		var comment Comment
		if err := db.First(&comment, id).Error; err != nil {
			c.JSON(404, gin.H{"error": "评论不存在"})
			return
		}

		// 检查权限：管理员或评论作者
		if !deleteData.IsAdmin {
			userID, err := stringToUint(deleteData.UserID)
			if err != nil || comment.UserID != userID {
				c.JSON(403, gin.H{"error": "无权限删除此评论"})
				return
			}
		}

		// 删除评论
		if err := db.Delete(&comment).Error; err != nil {
			c.JSON(500, gin.H{"error": "删除评论失败"})
			return
		}

		c.JSON(200, gin.H{"message": "删除评论成功"})
	})

	// 启动服务器
	fmt.Println("服务器启动成功，监听端口: 8082")
	fmt.Println("管理员账号: admin / admin123")
	r.Run(":8082")
}

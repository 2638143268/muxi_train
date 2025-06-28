package main

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var SecretKey = []byte("123456") // 密钥

// 登录请求结构体
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Token载荷结构体
type Claims struct {
	UserId int `json:"user_id"`
	jwt.RegisteredClaims
}

// 生成Token
func GenerateToken(userId int) (string, error) {
	claims := Claims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // 有效期1天
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SecretKey)
}

// 登录接口
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Username != "admin" || req.Password != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "用户名或密码错误"})
		return
	}

	token, err := GenerateToken(1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "生成token失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// JWT验证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "未提供token"})
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
			return SecretKey, nil
		})

		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			c.Set("user_id", claims.UserId)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "token无效", "error": err.Error()})
			c.Abort()
		}
	}
}

func main() {
	r := gin.Default()

	// 允许跨域请求
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 登录接口
	r.POST("/login", Login)

	// 受保护接口
	r.GET("/book/list", AuthMiddleware(), func(c *gin.Context) {
		userId, _ := c.Get("user_id")
		c.JSON(http.StatusOK, gin.H{"message": "图书列表", "user_id": userId})
	})

	r.Run(":8080")
}

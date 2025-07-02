package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

// SecretKey 用于签名和验证的密钥
var SecretKey = []byte("secret-key")
var expire = 30 * time.Minute

// GenerateToken 生成 JWT
func GenerateToken() (string, error) {
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
		Issuer:    "example",
		Subject:   "example",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(SecretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// RefreshToken 刷新 JWT
func RefreshToken(tokenString string) (string, error) {

	// 解析 token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})

	if err != nil || !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	// 验证 token 是否有效
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid claims")
	}
	exp, err := claims.GetExpirationTime()
	if err != nil {
		return "", errors.Wrap(err, "GetExpirationTime from token error")
	}

	if time.Until(exp.Time) < 0 {
		return "", errors.New("the token expires")
	}

	return GenerateToken()
}

// AuthMiddleware 中间件，用于验证请求中的 token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return SecretKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// 从 token 中提取 claims 并存入上下文
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			subject := claims["sub"] // 通常是用户 ID 或用户名
			c.Set("user", subject)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// token 解析成功，放行
		c.Next()
	}
}

// ExampleHandler 示例处理程序，需要通过 AuthMiddleware 进行身份验证
func ExampleHandler(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	c.String(http.StatusOK, fmt.Sprintf("Hello, %s!", user))
}

func main() {
	r := gin.Default()

	// 登录接口，返回 token
	r.POST("/login", func(c *gin.Context) {
		token, err := GenerateToken()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}
		c.String(http.StatusOK, token)
	})

	// 刷新接口，刷新 token
	r.POST("/refresh", func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		newToken, err := RefreshToken(tokenString)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to refresh token"})
			return
		}
		c.String(http.StatusOK, newToken)
	})

	// 受保护的路由，使用 AuthMiddleware
	r.GET("/example", AuthMiddleware(), ExampleHandler)

	fmt.Println("service start on :8080")
	r.Run(":8080")
}

package middleware

import (
	"log"
	"net/http"
	"product-checker/auth"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Пропускаем OPTIONS запросы (CORS preflight)
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		tokenString := extractToken(c)
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization required"})
			return
		}

		username, err := auth.ParseToken(tokenString)
		if err != nil {
			log.Printf("Token parsing error: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		log.Printf("User authenticated: %s", username)
		c.Set("username", username)
		c.Next()
	}
}

func extractToken(c *gin.Context) string {
	// 1. Проверяем заголовок Authorization
	if token := c.GetHeader("Authorization"); token != "" {
		if strings.HasPrefix(token, "Bearer ") {
			return strings.TrimPrefix(token, "Bearer ")
		}
		return token
	}

	// 2. Проверяем куки
	if token, err := c.Cookie("auth_token"); err == nil {
		return token
	}

	// 3. Проверяем query параметр
	if token := c.Query("token"); token != "" {
		return token
	}

	return ""
}

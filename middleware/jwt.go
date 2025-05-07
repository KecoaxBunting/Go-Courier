package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	helper "go-courier/helper"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Missing or invalid Authorization header",
			})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		_ = godotenv.Load("../../.env")
		key := os.Getenv("JWT_SECRET_KEY")
		if key == "" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "JWT Key not set",
			})
			return
		}

		claims := &helper.JWTClaim{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			_, ok := t.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(key), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid Token",
			})
			return
		}

		if claims.ExpiresAt.Time.Before(time.Now()) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Token expired",
			})
			return
		}

		c.Set("user", claims.UserId)
		c.Next()
	}
}

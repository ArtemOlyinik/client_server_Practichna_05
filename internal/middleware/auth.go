package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const UserIDKey = "userID"

// AuthMiddleware перевіряє Bearer-токен і зберігає userID у контексті.
func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
			return
		}

		tokenString := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, fmt.Errorf("неочікуваний метод підпису: %s", token.Method.Alg())
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
			return
		}

		userID, ok := claims["userID"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
			return
		}

		c.Set(UserIDKey, uint(userID))
		c.Next()
	}
}

func GetUserID(c *gin.Context) uint {
	value, exists := c.Get(UserIDKey)
	if !exists {
		return 0
	}
	userID, ok := value.(uint)
	if !ok {
		return 0
	}
	return userID
}

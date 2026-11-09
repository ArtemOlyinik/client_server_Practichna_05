package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	const secret = "test-secret"
	router := gin.New()
	router.Use(AuthMiddleware(secret))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"userID": GetUserID(c)})
	})

	validToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": float64(42),
		"exp":    time.Now().Add(time.Hour).Unix(),
	})
	signed, err := validToken.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("не вдалося підписати токен: %v", err)
	}

	request := httptest.NewRequest(http.MethodGet, "/protected", nil)
	request.Header.Set("Authorization", "Bearer "+signed)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("очікували 200, отримали %d", recorder.Code)
	}

	var response map[string]float64
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("не вдалося розібрати відповідь: %v", err)
	}
	if response["userID"] != 42 {
		t.Fatalf("очікували userID=42, отримали %v", response["userID"])
	}
}

func TestAuthMiddlewareRejectsMissingToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(AuthMiddleware("test-secret"))
	router.GET("/protected", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	request := httptest.NewRequest(http.MethodGet, "/protected", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("очікували 401, отримали %d", recorder.Code)
	}
}

func TestAuthMiddlewareRejectsExpiredToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	const secret = "test-secret"
	expired := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": float64(1),
		"exp":    time.Now().Add(-time.Hour).Unix(),
	})
	signed, err := expired.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("не вдалося підписати токен: %v", err)
	}

	router := gin.New()
	router.Use(AuthMiddleware(secret))
	router.GET("/protected", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	request := httptest.NewRequest(http.MethodGet, "/protected", nil)
	request.Header.Set("Authorization", "Bearer "+signed)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("очікували 401 для простроченого токена, отримали %d", recorder.Code)
	}
}

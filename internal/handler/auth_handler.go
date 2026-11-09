package handler

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/ArtemOlyinik/client_server_Practichna_05/internal/domain"
	"github.com/ArtemOlyinik/client_server_Practichna_05/internal/dto"
)

type AuthHandler struct {
	db     *gorm.DB
	secret string
}

func NewAuthHandler(db *gorm.DB, secret string) *AuthHandler {
	return &AuthHandler{db: db, secret: secret}
}

// Register godoc
// @Summary Реєстрація інженера
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterDTO true "Дані реєстрації"
// @Success 201 {object} domain.User
// @Failure 400 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var input dto.RegisterDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Validation Error",
			Details: validationErrors(err),
		})
		return
	}

	email := strings.ToLower(strings.TrimSpace(input.Email))

	var existing domain.User
	err := h.db.Where("email = ?", email).First(&existing).Error
	if err == nil {
		c.JSON(http.StatusConflict, dto.ErrorResponse{Error: "Email already registered"})
		return
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Database error"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Password hashing failed"})
		return
	}

	user := domain.User{
		Email:        email,
		PasswordHash: string(hash),
	}
	if err := h.db.Create(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			c.JSON(http.StatusConflict, dto.ErrorResponse{Error: "Email already registered"})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Database error"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// Login godoc
// @Summary Вхід інженера
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginDTO true "Облікові дані"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var input dto.LoginDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Validation Error",
			Details: validationErrors(err),
		})
		return
	}

	var user domain.User
	err := h.db.Where("email = ?", strings.ToLower(strings.TrimSpace(input.Email))).First(&user).Error
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "Invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "Invalid email or password"})
		return
	}

	claims := jwt.MapClaims{
		"userID": user.ID,
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
		"iat":    time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(h.secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Token generation failed"})
		return
	}

	c.JSON(http.StatusOK, dto.LoginResponse{Token: signed})
}

package handler

import (
	"errors"
	"net/http"

	"homework04/internal/auth"
	"homework04/internal/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type registerRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		AbortWithLog(c, h.Logger, http.StatusInternalServerError, "failed to hash password", err)
		return
	}

	user := models.User{
		Username: req.Username,
		Password: string(hashed),
		Email:    req.Email,
	}
	if err := h.DB.Create(&user).Error; err != nil {
		AbortWithLog(c, h.Logger, http.StatusInternalServerError, "failed to create user", err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user registered successfully",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

func (h *Handler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}

	var stored models.User
	if err := h.DB.Where("username = ?", req.Username).First(&stored).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			JSONError(c, http.StatusUnauthorized, "invalid username or password")
			return
		}
		AbortWithLog(c, h.Logger, http.StatusInternalServerError, "database error", err)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(stored.Password), []byte(req.Password)); err != nil {
		JSONError(c, http.StatusUnauthorized, "invalid username or password")
		return
	}

	token, err := auth.IssueToken(stored.ID, stored.Username)
	if err != nil {
		AbortWithLog(c, h.Logger, http.StatusInternalServerError, "failed to generate token", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       stored.ID,
			"username": stored.Username,
			"email":    stored.Email,
		},
	})
}

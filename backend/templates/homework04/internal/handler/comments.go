package handler

import (
	"errors"
	"net/http"
	"strconv"

	"homework04/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type commentBody struct {
	Content string `json:"content" binding:"required"`
}

func (h *Handler) ListComments(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		JSONError(c, http.StatusBadRequest, "invalid post id")
		return
	}

	var post models.Post
	if err := h.DB.First(&post, postID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			JSONError(c, http.StatusNotFound, "post not found")
			return
		}
		AbortWithLog(c, h.Logger, http.StatusInternalServerError, "database error", err)
		return
	}

	var comments []models.Comment
	if err := h.DB.Preload("User").Where("post_id = ?", postID).Order("id asc").Find(&comments).Error; err != nil {
		AbortWithLog(c, h.Logger, http.StatusInternalServerError, "failed to list comments", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"comments": comments})
}

func (h *Handler) CreateComment(c *gin.Context) {
	userID := c.GetUint("user_id")
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		JSONError(c, http.StatusBadRequest, "invalid post id")
		return
	}

	var post models.Post
	if err := h.DB.First(&post, postID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			JSONError(c, http.StatusNotFound, "post not found")
			return
		}
		AbortWithLog(c, h.Logger, http.StatusInternalServerError, "database error", err)
		return
	}

	var body commentBody
	if err := c.ShouldBindJSON(&body); err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}

	comment := models.Comment{
		Content: body.Content,
		UserID:  userID,
		PostID:  uint(postID),
	}
	if err := h.DB.Create(&comment).Error; err != nil {
		AbortWithLog(c, h.Logger, http.StatusInternalServerError, "failed to create comment", err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"comment": comment})
}

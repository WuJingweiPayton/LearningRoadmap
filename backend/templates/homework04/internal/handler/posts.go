package handler

import (
	"errors"
	"net/http"
	"strconv"

	"homework04/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type postBody struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func (h *Handler) ListPosts(c *gin.Context) {
	var posts []models.Post
	if err := h.DB.Preload("User").Order("id desc").Find(&posts).Error; err != nil {
		AbortWithLog(c, h.Logger, http.StatusInternalServerError, "failed to list posts", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

func (h *Handler) GetPost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		JSONError(c, http.StatusBadRequest, "invalid post id")
		return
	}

	var post models.Post
	if err := h.DB.Preload("User").First(&post, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			JSONError(c, http.StatusNotFound, "post not found")
			return
		}
		AbortWithLog(c, h.Logger, http.StatusInternalServerError, "database error", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"post": post})
}

func (h *Handler) CreatePost(c *gin.Context) {
	userID := c.GetUint("user_id")

	var body postBody
	if err := c.ShouldBindJSON(&body); err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}

	post := models.Post{
		Title:   body.Title,
		Content: body.Content,
		UserID:  userID,
	}
	if err := h.DB.Create(&post).Error; err != nil {
		AbortWithLog(c, h.Logger, http.StatusInternalServerError, "failed to create post", err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"post": post})
}

func (h *Handler) UpdatePost(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		JSONError(c, http.StatusBadRequest, "invalid post id")
		return
	}

	var post models.Post
	if err := h.DB.First(&post, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			JSONError(c, http.StatusNotFound, "post not found")
			return
		}
		AbortWithLog(c, h.Logger, http.StatusInternalServerError, "database error", err)
		return
	}
	if post.UserID != userID {
		JSONError(c, http.StatusForbidden, "only the author can update this post")
		return
	}

	var body postBody
	if err := c.ShouldBindJSON(&body); err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	post.Title = body.Title
	post.Content = body.Content
	if err := h.DB.Save(&post).Error; err != nil {
		AbortWithLog(c, h.Logger, http.StatusInternalServerError, "failed to update post", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"post": post})
}

func (h *Handler) DeletePost(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		JSONError(c, http.StatusBadRequest, "invalid post id")
		return
	}

	var post models.Post
	if err := h.DB.First(&post, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			JSONError(c, http.StatusNotFound, "post not found")
			return
		}
		AbortWithLog(c, h.Logger, http.StatusInternalServerError, "database error", err)
		return
	}
	if post.UserID != userID {
		JSONError(c, http.StatusForbidden, "only the author can delete this post")
		return
	}

	if err := h.DB.Delete(&post).Error; err != nil {
		AbortWithLog(c, h.Logger, http.StatusInternalServerError, "failed to delete post", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "post deleted"})
}

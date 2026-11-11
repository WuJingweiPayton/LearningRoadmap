package main

import (
	"errors"

	"gorm.io/gorm"
)

// GetUserPostsWithComments 查询指定用户发布的所有文章及其评论。
func GetUserPostsWithComments(db *gorm.DB, userID uint) (*User, error) {
	var user User
	err := db.Preload("Posts.Comments").First(&user, userID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetPostWithMostComments 查询评论数量最多的文章（并列时取 id 最小的一条）。
func GetPostWithMostComments(db *gorm.DB) (*Post, error) {
	var post Post
	err := db.Model(&Post{}).
		Select("posts.*, COUNT(comments.id) AS comment_count").
		Joins("LEFT JOIN comments ON comments.post_id = posts.id AND comments.deleted_at IS NULL").
		Group("posts.id").
		Order("comment_count DESC, posts.id ASC").
		Preload("Comments").
		First(&post).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &post, nil
}

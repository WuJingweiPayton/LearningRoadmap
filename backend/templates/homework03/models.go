package main

import (
	"gorm.io/gorm"
)

const CommentStatusNone = "无评论"
const CommentStatusHas = "有评论"

// User 博客用户，与 Post 为一对多关系。
type User struct {
	gorm.Model
	Username  string `gorm:"size:64;not null"`
	PostCount int    `gorm:"default:0"`
	Posts     []Post `gorm:"foreignKey:UserID"`
}

// Post 文章，属于 User，与 Comment 为一对多关系。
type Post struct {
	gorm.Model
	UserID        uint      `gorm:"index;not null"`
	Title         string    `gorm:"size:200;not null"`
	Content       string    `gorm:"type:text"`
	CommentStatus string    `gorm:"size:32;default:无评论"`
	User          User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Comments      []Comment `gorm:"foreignKey:PostID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// Comment 评论，属于 Post。
type Comment struct {
	gorm.Model
	PostID  uint   `gorm:"index;not null"`
	Content string `gorm:"type:text;not null"`
	Post    Post   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// AfterCreate 创建文章后自动增加用户的文章数量统计。
func (p *Post) AfterCreate(tx *gorm.DB) error {
	return tx.Model(&User{}).
		Where("id = ?", p.UserID).
		UpdateColumn("post_count", gorm.Expr("post_count + ?", 1)).
		Error
}

// AfterDelete 删除评论后若文章已无评论，则将评论状态更新为「无评论」。
func (c *Comment) AfterDelete(tx *gorm.DB) error {
	var count int64
	if err := tx.Model(&Comment{}).Where("post_id = ?", c.PostID).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return tx.Model(&Post{}).Where("id = ?", c.PostID).Update("comment_status", CommentStatusNone).Error
	}
	return nil
}

// MigrateModels 创建 User、Post、Comment 对应的数据库表。
func MigrateModels(db *gorm.DB) error {
	return db.AutoMigrate(&User{}, &Post{}, &Comment{})
}

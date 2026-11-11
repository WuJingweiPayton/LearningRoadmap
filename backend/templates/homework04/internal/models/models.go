package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;size:64;not null" json:"username"`
	Password string `gorm:"not null" json:"-"`
	Email    string `gorm:"uniqueIndex;size:128;not null" json:"email"`
	Posts    []Post `gorm:"foreignKey:UserID"`
}

type Post struct {
	gorm.Model
	Title   string `gorm:"size:200;not null" json:"title"`
	Content string `gorm:"type:text;not null" json:"content"`
	UserID  uint   `gorm:"index;not null" json:"user_id"`
	User    User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`
}

type Comment struct {
	gorm.Model
	Content string `gorm:"type:text;not null" json:"content"`
	UserID  uint   `gorm:"index;not null" json:"user_id"`
	PostID  uint   `gorm:"index;not null" json:"post_id"`
	User    User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`
	Post    Post   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

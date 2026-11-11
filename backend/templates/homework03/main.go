package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	db, err := OpenMySQLFromEnv()
	if err != nil {
		log.Fatalf("open mysql: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("db handle: %v", err)
	}
	for i := 0; i < 30; i++ {
		if pingErr := sqlDB.Ping(); pingErr == nil {
			break
		}
		if i == 29 {
			log.Fatalf("mysql not ready; ensure mysql-local is running (docker start mysql-local)")
		}
		time.Sleep(time.Second)
	}

	if err := MigrateModels(db); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	user := User{Username: "alice"}
	if err := db.Create(&user).Error; err != nil {
		log.Fatalf("create user: %v", err)
	}

	posts := []Post{
		{UserID: user.ID, Title: "Hello GORM", Content: "first post", CommentStatus: CommentStatusNone},
		{UserID: user.ID, Title: "Advanced GORM", Content: "second post", CommentStatus: CommentStatusNone},
	}
	for i := range posts {
		if err := db.Create(&posts[i]).Error; err != nil {
			log.Fatalf("create post: %v", err)
		}
	}

	comments := []Comment{
		{PostID: posts[0].ID, Content: "nice"},
		{PostID: posts[0].ID, Content: "helpful"},
		{PostID: posts[1].ID, Content: "great"},
	}
	for i := range comments {
		if err := db.Create(&comments[i]).Error; err != nil {
			log.Fatalf("create comment: %v", err)
		}
	}
	_ = db.Model(&posts[0]).Update("comment_status", CommentStatusHas)
	_ = db.Model(&posts[1]).Update("comment_status", CommentStatusHas)

	u, err := GetUserPostsWithComments(db, user.ID)
	if err != nil {
		log.Fatalf("query user posts: %v", err)
	}
	fmt.Printf("用户 %s 共 %d 篇文章（统计字段 post_count=%d）\n", u.Username, len(u.Posts), u.PostCount)
	for _, p := range u.Posts {
		fmt.Printf("  - %s: %d 条评论\n", p.Title, len(p.Comments))
	}

	top, err := GetPostWithMostComments(db)
	if err != nil {
		log.Fatalf("query top post: %v", err)
	}
	if top != nil {
		fmt.Printf("评论最多的文章: %s (%d 条评论)\n", top.Title, len(top.Comments))
	}
}

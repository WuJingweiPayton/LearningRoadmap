package main

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func openTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := MigrateModels(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func TestMigrateModels(t *testing.T) {
	db := openTestDB(t)
	if !db.Migrator().HasTable(&User{}) || !db.Migrator().HasTable(&Post{}) || !db.Migrator().HasTable(&Comment{}) {
		t.Fatal("tables were not created")
	}
}

func TestPostAfterCreateUpdatesUserPostCount(t *testing.T) {
	db := openTestDB(t)

	user := User{Username: "bob"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 3; i++ {
		post := Post{UserID: user.ID, Title: "post", Content: "body", CommentStatus: CommentStatusNone}
		if err := db.Create(&post).Error; err != nil {
			t.Fatal(err)
		}
	}

	var loaded User
	if err := db.First(&loaded, user.ID).Error; err != nil {
		t.Fatal(err)
	}
	if loaded.PostCount != 3 {
		t.Errorf("expected post_count=3, got %d", loaded.PostCount)
	}
}

func TestGetUserPostsWithComments(t *testing.T) {
	db := openTestDB(t)

	user := User{Username: "carol"}
	db.Create(&user)

	p1 := Post{UserID: user.ID, Title: "A", CommentStatus: CommentStatusHas}
	p2 := Post{UserID: user.ID, Title: "B", CommentStatus: CommentStatusNone}
	db.Create(&p1)
	db.Create(&p2)
	db.Create(&Comment{PostID: p1.ID, Content: "c1"})
	db.Create(&Comment{PostID: p1.ID, Content: "c2"})

	got, err := GetUserPostsWithComments(db, user.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Posts) != 2 {
		t.Fatalf("expected 2 posts, got %d", len(got.Posts))
	}
	if len(got.Posts[0].Comments)+len(got.Posts[1].Comments) != 2 {
		t.Error("expected 2 comments total on user posts")
	}
}

func TestGetPostWithMostComments(t *testing.T) {
	db := openTestDB(t)

	user := User{Username: "dave"}
	db.Create(&user)

	p1 := Post{UserID: user.ID, Title: "few", CommentStatus: CommentStatusHas}
	p2 := Post{UserID: user.ID, Title: "most", CommentStatus: CommentStatusHas}
	db.Create(&p1)
	db.Create(&p2)

	db.Create(&Comment{PostID: p1.ID, Content: "one"})
	for i := 0; i < 5; i++ {
		db.Create(&Comment{PostID: p2.ID, Content: "many"})
	}

	top, err := GetPostWithMostComments(db)
	if err != nil {
		t.Fatal(err)
	}
	if top == nil || top.Title != "most" || len(top.Comments) != 5 {
		t.Errorf("expected post 'most' with 5 comments, got %+v", top)
	}
}

func TestCommentAfterDeleteSetsStatusWhenEmpty(t *testing.T) {
	db := openTestDB(t)

	user := User{Username: "eve"}
	db.Create(&user)

	post := Post{UserID: user.ID, Title: "solo", CommentStatus: CommentStatusHas}
	db.Create(&post)

	c := Comment{PostID: post.ID, Content: "only one"}
	db.Create(&c)

	if err := db.Delete(&c).Error; err != nil {
		t.Fatal(err)
	}

	var loaded Post
	db.First(&loaded, post.ID)
	if loaded.CommentStatus != CommentStatusNone {
		t.Errorf("expected comment_status=%q, got %q", CommentStatusNone, loaded.CommentStatus)
	}
}

package database

import (
	"fmt"
	"os"
	"time"

	"homework04/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// OpenFromEnv 默认连接本机已有 MySQL（mysql-local）；设置 DB_DRIVER=sqlite 可改用 SQLite 文件。
func OpenFromEnv() (*gorm.DB, error) {
	driver := envOr("DB_DRIVER", "mysql")
	switch driver {
	case "sqlite":
		dsn := envOr("DB_DSN", "blog.db")
		return gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	default:
		if dsn := os.Getenv("MYSQL_DSN"); dsn != "" {
			return gorm.Open(mysql.Open(dsn), &gorm.Config{})
		}
		user := envOr("MYSQL_USER", "appuser")
		pass := envOr("MYSQL_PASSWORD", "apppass123")
		host := envOr("MYSQL_HOST", "127.0.0.1")
		port := envOr("MYSQL_PORT", "3306")
		dbName := envOr("MYSQL_DATABASE", "blog_api")
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			user, pass, host, port, dbName)
		return gorm.Open(mysql.Open(dsn), &gorm.Config{})
	}
}

func OpenSQLiteMemory() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{})
}

func WaitReady(db *gorm.DB, attempts int) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	for i := 0; i < attempts; i++ {
		if pingErr := sqlDB.Ping(); pingErr == nil {
			return nil
		}
		time.Sleep(time.Second)
	}
	return fmt.Errorf("database not ready after %d attempts", attempts)
}

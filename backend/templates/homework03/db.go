package main

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// OpenMySQLFromEnv 连接本机已有 MySQL（默认对齐数算协同 mysql-local:3306）。
// 可直接设置 MYSQL_DSN 覆盖整条连接串。
func OpenMySQLFromEnv() (*gorm.DB, error) {
	if dsn := os.Getenv("MYSQL_DSN"); dsn != "" {
		return gorm.Open(mysql.Open(dsn), &gorm.Config{})
	}
	user := envOr("MYSQL_USER", "appuser")
	pass := envOr("MYSQL_PASSWORD", "apppass123")
	host := envOr("MYSQL_HOST", "127.0.0.1")
	port := envOr("MYSQL_PORT", "3306")
	dbName := envOr("MYSQL_DATABASE", "blog_homework")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, dbName)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

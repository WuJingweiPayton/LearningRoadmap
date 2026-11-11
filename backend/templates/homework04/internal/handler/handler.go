package handler

import (
	"log/slog"

	"gorm.io/gorm"
)

type Handler struct {
	DB     *gorm.DB
	Logger *slog.Logger
}

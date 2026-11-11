package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func JSONError(c *gin.Context, code int, msg string) {
	c.JSON(code, gin.H{"error": msg})
}

func AbortWithLog(c *gin.Context, logger *slog.Logger, code int, msg string, err error) {
	if err != nil {
		logger.Error(msg, slog.Any("error", err), slog.String("path", c.Request.URL.Path))
	} else {
		logger.Warn(msg, slog.String("path", c.Request.URL.Path))
	}
	JSONError(c, code, msg)
}

func NotFound(c *gin.Context) {
	JSONError(c, http.StatusNotFound, "resource not found")
}

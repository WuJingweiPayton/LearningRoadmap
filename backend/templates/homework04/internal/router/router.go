package router

import (
	"log/slog"

	"homework04/internal/handler"
	"homework04/internal/middleware"

	"github.com/gin-gonic/gin"
)

func New(h *handler.Handler, logger *slog.Logger) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestLogger(logger))

	r.NoRoute(handler.NotFound)

	api := r.Group("/api")
	{
		api.POST("/register", h.Register)
		api.POST("/login", h.Login)

		api.GET("/posts", h.ListPosts)
		api.GET("/posts/:id", h.GetPost)
		api.GET("/posts/:id/comments", h.ListComments)

		auth := api.Group("")
		auth.Use(middleware.AuthRequired())
		{
			auth.POST("/posts", h.CreatePost)
			auth.PUT("/posts/:id", h.UpdatePost)
			auth.DELETE("/posts/:id", h.DeletePost)
			auth.POST("/posts/:id/comments", h.CreateComment)
		}
	}

	return r
}

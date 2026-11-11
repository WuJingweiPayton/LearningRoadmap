package main

import (
	"log"
	"log/slog"
	"os"

	"homework04/internal/database"
	"homework04/internal/handler"
	"homework04/internal/router"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	db, err := database.OpenFromEnv()
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}
	if err := database.WaitReady(db, 30); err != nil {
		logger.Warn("database ping failed, continuing anyway", slog.Any("error", err))
	}
	if err := database.Migrate(db); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	h := &handler.Handler{DB: db, Logger: logger}
	r := router.New(h, logger)
	logger.Info("starting server", slog.String("addr", addr))
	if err := r.Run(addr); err != nil {
		log.Fatalf("server: %v", err)
	}
}

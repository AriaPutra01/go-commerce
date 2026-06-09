package logger

import (
	"os"

	"github.com/AriaPutra01/go-commerce/internal/config"

	"log/slog"
)

func NewLogger(cfg *config.Config) *slog.Logger {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.Level(cfg.LogLevel),
	}))

	return log
}

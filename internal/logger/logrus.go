package logger

import (
	"github.com/AriaPutra01/go-commerce/internal/config"

	"github.com/sirupsen/logrus"
)

func NewLogger(cfg *config.Config) *logrus.Logger {
	log := logrus.New()

	log.SetLevel(logrus.Level(cfg.LogLevel))
	log.SetFormatter(&logrus.JSONFormatter{})

	return log
}

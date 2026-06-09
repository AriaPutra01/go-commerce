package database

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/AriaPutra01/go-commerce/internal/config"
	slogGorm "github.com/orandin/slog-gorm"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase(cfg *config.Config, log *slog.Logger) *gorm.DB {
	gormLogger := slogGorm.New(
		slogGorm.WithHandler(log.Handler()),
		slogGorm.WithTraceAll(),
		slogGorm.WithSlowThreshold(200*time.Millisecond),
		slogGorm.SetLogLevel(slogGorm.ErrorLogType, slog.LevelError),
		slogGorm.SetLogLevel(slogGorm.SlowQueryLogType, slog.LevelWarn),
		slogGorm.SetLogLevel(slogGorm.DefaultLogType, slog.LevelDebug),
	)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable search_path=%s",
		cfg.DBHost, cfg.DBUsername, cfg.DBPassword, cfg.DBDatabase, cfg.DBPort, cfg.DBSchema)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: gormLogger})
	if err != nil {
		log.Error(fmt.Sprintf("failed to connect database: %v", err))
		os.Exit(1)
	}

	connection, err := db.DB()
	if err != nil {
		log.Error(fmt.Sprintf("failed to connect database: %v", err))
		os.Exit(1)
	}

	connection.SetMaxIdleConns(cfg.DBPoolIdle)
	connection.SetMaxOpenConns(cfg.DBPoolMax)
	connection.SetConnMaxLifetime(time.Second * time.Duration(cfg.DBPoolLifetime))

	return db
}

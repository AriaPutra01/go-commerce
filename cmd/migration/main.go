package main

import (
	"fmt"
	"log"
	"os"

	"github.com/AriaPutra01/go-commerce/internal/config"
	"github.com/AriaPutra01/go-commerce/internal/database"
	"github.com/AriaPutra01/go-commerce/internal/logger"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: go run cmd/migration/main.go [up|down|version|force <version>]")
	}

	cfg := config.Load()
	log := logger.NewLogger(cfg)
	db := database.NewDatabase(cfg, log)

	conn, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("database error: %s", err))
	}

	database.RunMigration(&database.MigrationConfig{
		Config: cfg,
		DB:     conn,
		Log:    log,
	})
}

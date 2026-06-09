package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/AriaPutra01/go-commerce/internal/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type MigrationConfig struct {
	Config *config.Config
	DB     *sql.DB
	Log    *slog.Logger
}

func RunMigration(cfg *MigrationConfig) {
	driver, err := postgres.WithInstance(cfg.DB, &postgres.Config{
		SchemaName: cfg.Config.DBSchema,
	})
	if err != nil {
		cfg.Log.Error(fmt.Sprintf("failed to create migration postgres driver: %v", err))
		os.Exit(1)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file:%s", cfg.Config.DBMigrationPath),
		"postgres", driver,
	)
	if err != nil {
		cfg.Log.Error(fmt.Sprintf("failed to initialize migration instance: %v", err))
		os.Exit(1)
	}

	switch os.Args[1] {
	case "up":
		err = m.Up()
		if err != nil && !errors.Is(err, migrate.ErrNoChange) {
			cfg.Log.Error(fmt.Sprintf("migration up failed: %v", err))
			os.Exit(1)
		}
		cfg.Log.Info("migration up completed")
	case "down":
		err = m.Steps(-1)
		if err != nil && !errors.Is(err, migrate.ErrNoChange) {
			cfg.Log.Error(fmt.Sprintf("migration down failed: %v", err))
			os.Exit(1)
		}
		cfg.Log.Info("migration down completed")
	case "version":
		version, dirty, vErr := m.Version()
		if vErr != nil {
			if errors.Is(vErr, migrate.ErrNilVersion) {
				cfg.Log.Info("migration version: none")
				return
			}
			cfg.Log.Error(fmt.Sprintf("failed to get migration version: %v", vErr))
			os.Exit(1)
		}
		cfg.Log.Info(fmt.Sprintf("migration version: %d (dirty=%t)\n", version, dirty))
	case "force":
		if len(os.Args) < 3 {
			cfg.Log.Error("usage: go run cmd/migration/main.go force <version>")
			os.Exit(1)
		}
		var version int
		_, scanErr := fmt.Sscanf(os.Args[2], "%d", &version)
		if scanErr != nil {
			cfg.Log.Error(fmt.Sprintf("invalid force version: %v", scanErr))
			os.Exit(1)
		}
		if err = m.Force(version); err != nil {
			cfg.Log.Error(fmt.Sprintf("migration force failed: %v", err))
			os.Exit(1)
		}
		cfg.Log.Info(fmt.Sprintf("migration forced to version: %d\n", version))
	default:
		cfg.Log.Error(fmt.Sprintf("unknown command: %s", os.Args[1]))
		os.Exit(1)
	}

	cfg.Log.Info("database migration executed successfully")
}

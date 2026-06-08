package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/AriaPutra01/go-commerce/internal/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/sirupsen/logrus"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type MigrationConfig struct {
	Config *config.Config
	DB     *sql.DB
	Log    *logrus.Logger
}

func RunMigration(cfg *MigrationConfig) {
	driver, err := postgres.WithInstance(cfg.DB, &postgres.Config{
		SchemaName: cfg.Config.DBSchema,
	})
	if err != nil {
		cfg.Log.Fatalf("failed to create migration postgres driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file:%s", cfg.Config.DBMigrationPath),
		"postgres", driver,
	)
	if err != nil {
		cfg.Log.Fatalf("failed to initialize migration instance: %v", err)
	}

	switch os.Args[1] {
	case "up":
		err = m.Up()
		if err != nil && !errors.Is(err, migrate.ErrNoChange) {
			cfg.Log.Fatalf("migration up failed: %v", err)
		}
		cfg.Log.Println("migration up completed")
	case "down":
		err = m.Steps(-1)
		if err != nil && !errors.Is(err, migrate.ErrNoChange) {
			cfg.Log.Fatalf("migration down failed: %v", err)
		}
		cfg.Log.Println("migration down completed")
	case "version":
		version, dirty, vErr := m.Version()
		if vErr != nil {
			if errors.Is(vErr, migrate.ErrNilVersion) {
				cfg.Log.Println("migration version: none")
				return
			}
			cfg.Log.Fatalf("failed to get migration version: %v", vErr)
		}
		cfg.Log.Printf("migration version: %d (dirty=%t)\n", version, dirty)
	case "force":
		if len(os.Args) < 3 {
			cfg.Log.Fatal("usage: go run cmd/migration/main.go force <version>")
		}
		var version int
		_, scanErr := fmt.Sscanf(os.Args[2], "%d", &version)
		if scanErr != nil {
			cfg.Log.Fatalf("invalid force version: %v", scanErr)
		}
		if err = m.Force(version); err != nil {
			cfg.Log.Fatalf("migration force failed: %v", err)
		}
		cfg.Log.Printf("migration forced to version: %d\n", version)
	default:
		cfg.Log.Fatalf("unknown command: %s", os.Args[1])
	}

	cfg.Log.Info("database migration executed successfully")
}

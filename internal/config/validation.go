package config

import "fmt"

func validate(cfg *Config) {
	if cfg.DBHost == "" {
		panic(fmt.Errorf("DB_HOST is required"))
	}

	if cfg.DBDatabase == "" {
		panic(fmt.Errorf("DB_DATABASE is required"))
	}

	if cfg.DBUsername == "" {
		panic(fmt.Errorf("DB_USERNAME is required"))
	}
}

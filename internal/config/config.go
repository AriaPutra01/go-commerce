package config

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Port   string `mapstructure:"PORT"`
	AppEnv string `mapstructure:"APP_ENV"`

	DBHost          string `mapstructure:"DB_HOST"`
	DBPort          int    `mapstructure:"DB_PORT"`
	DBDatabase      string `mapstructure:"DB_DATABASE"`
	DBUsername      string `mapstructure:"DB_USERNAME"`
	DBPassword      string `mapstructure:"DB_PASSWORD"`
	DBSchema        string `mapstructure:"DB_SCHEMA"`
	DBMigrationPath string `mapstructure:"DB_MIGRATION_PATH"`

	DBPoolIdle     int `mapstructure:"DB_POOL_IDLE"`
	DBPoolMax      int `mapstructure:"DB_POOL_MAX"`
	DBPoolLifetime int `mapstructure:"DB_POOL_LIFETIME"`

	LogLevel int `mapstructure:"LOG_LEVEL"`
}

func Load() *Config {
	v := viper.New()
	v.AutomaticEnv()
	v.SetConfigFile(".env")
	v.SetConfigType("env")
	if err := v.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) {
			panic(fmt.Errorf("read config: %w", err))
		}
	}
	setDefaults(v)
	cfg := new(Config)
	if err := v.Unmarshal(cfg); err != nil {
		panic(fmt.Errorf("unmarshal config: %w", err))
	}
	validate(cfg)
	return cfg
}

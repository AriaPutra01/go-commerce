package config

import "github.com/spf13/viper"

func setDefaults(v *viper.Viper) {
	v.SetDefault("PORT", 8080)
	v.SetDefault("APP_ENV", "local")

	v.SetDefault("DB_PORT", 5432)
	v.SetDefault("DB_SCHEMA", "public")
	v.SetDefault("DB_POOL_IDLE", 10)
	v.SetDefault("DB_POOL_MAX", 100)
	v.SetDefault("DB_POOL_LIFETIME", 300)

	v.SetDefault("LOG_LEVEL", 6)
}

package config

import (
	"os"
)

type Config struct {
	AppPort string

	DatabaseHost     string
	DatabasePort     string
	DatabaseUser     string
	DatabasePassword string
	DatabaseName     string
	DatabaseSSLMode  string
}

func Load() *Config {
	return &Config{
		AppPort:          getEnv("APP_PORT", "8080"),
		DatabaseHost:     getEnv("DB_HOST", "localhost"),
		DatabasePort:     getEnv("DB_PORT", "5432"),
		DatabaseUser:     getEnv("DB_USER", "postgres"),
		DatabasePassword: getEnv("DB_PASSWORD", "secret"),
		DatabaseName:     getEnv("DB_NAME", "mini_erp"),
		DatabaseSSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
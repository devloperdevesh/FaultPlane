package config

import (
	"os"
)

type Config struct {
	Host string
	Port string
	LogLevel string
}

func Load() Config {
	return Config{
		Host: getEnv("FAULTPLANE_HOST", "0.0.0.0"),
		Port: getEnv("FAULTPLANE_PORT", "8080"),
		LogLevel: getEnv("FAULTPLANE_LOG_LEVEL", "info"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}
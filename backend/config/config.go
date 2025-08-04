package config

import (
	"os"
)

type Config struct {
	DB     DBConfig
	Server ServerConfig
}

type DBConfig struct {
	URL string
}

type ServerConfig struct {
	Port string
}

func Load() *Config {
	return &Config{
		DB: DBConfig{
			URL: os.Getenv("DATABASE_URL"), // Railway style
		},
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
		},
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

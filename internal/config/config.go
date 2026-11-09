package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port       string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
}

func Load() Config {
	return Config{
		Port:       env("PORT", "8080"),
		DBHost:     env("DB_HOST", "localhost"),
		DBPort:     env("DB_PORT", "5432"),
		DBUser:     env("DB_USER", "postgres"),
		DBPassword: env("DB_PASSWORD", "postgres"),
		DBName:     env("DB_NAME", "incident_tracker"),
		JWTSecret:  env("JWT_SECRET", "dev-only-secret-change-me"),
	}
}

func (c Config) DSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		c.DBHost,
		c.DBUser,
		c.DBPassword,
		c.DBName,
		c.DBPort,
	)
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

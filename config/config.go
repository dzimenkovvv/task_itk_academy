package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBMaxConns int32
	ServerPort string
}

func Load(path string) (*Config, error) {
	if err := godotenv.Load(path); err != nil {
		return nil, fmt.Errorf("Load config: %w", err)
	}

	maxConns, err := strconv.Atoi(getEnv("DB_MAX_CONNS", "10"))
	if err != nil {
		return nil, fmt.Errorf("parse DB_MAX_CONNS: %w", err)
	}

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "postgres"),
		DBMaxConns: int32(maxConns),
		ServerPort: getEnv("SERVER_PORT", "9091"),
	}, nil
}

func (c *Config) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName,
	)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

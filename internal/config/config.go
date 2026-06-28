package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	AppName         string
	Env             string
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
	AllowedOrigins  []string
	Postgres        PostgresConfig
	Mongo           MongoConfig
	Auth            AuthConfig
}

type PostgresConfig struct {
	URL             string
	MaxOpenConns    int
	MinOpenConns    int
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}

type MongoConfig struct {
	URL      string
	Database string
}

type AuthConfig struct {
	JWTSecret         string
	AccessTokenExpiry time.Duration
}

func Load() Config {
	return Config{
		AppName:         getEnv("APP_NAME", "restaurant-inventory-api"),
		Env:             getEnv("APP_ENV", "development"),
		Port:            getEnv("PORT", "8080"),
		ReadTimeout:     getDurationEnv("HTTP_READ_TIMEOUT_SECONDS", 10*time.Second),
		WriteTimeout:    getDurationEnv("HTTP_WRITE_TIMEOUT_SECONDS", 15*time.Second),
		IdleTimeout:     getDurationEnv("HTTP_IDLE_TIMEOUT_SECONDS", 60*time.Second),
		ShutdownTimeout: getDurationEnv("HTTP_SHUTDOWN_TIMEOUT_SECONDS", 10*time.Second),
		AllowedOrigins:  getCSVEnv("CORS_ALLOWED_ORIGINS", []string{"*"}),
		Postgres: PostgresConfig{
			URL:             getEnv("POSTGRES_URL", ""),
			MaxOpenConns:    getIntEnv("POSTGRES_MAX_OPEN_CONNS", 25),
			MinOpenConns:    getIntEnv("POSTGRES_MIN_OPEN_CONNS", 2),
			MaxConnLifetime: getDurationEnv("POSTGRES_MAX_CONN_LIFETIME_SECONDS", 1800*time.Second),
			MaxConnIdleTime: getDurationEnv("POSTGRES_MAX_CONN_IDLE_TIME_SECONDS", 300*time.Second),
		},
		Mongo: MongoConfig{
			URL:      getEnv("MONGO_URL", ""),
			Database: getEnv("MONGO_DATABASE", "restaurant_inventory_events"),
		},
		Auth: AuthConfig{
			JWTSecret:         getEnv("JWT_SECRET", "change-me-in-production"),
			AccessTokenExpiry: getDurationEnv("JWT_ACCESS_TOKEN_EXPIRY_SECONDS", 86400*time.Second),
		},
	}
}

func (c Config) HTTPAddr() string {
	return fmt.Sprintf(":%s", c.Port)
}

func (c Config) IsProduction() bool {
	return c.Env == "production"
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func getDurationEnv(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	seconds, err := strconv.Atoi(value)
	if err != nil || seconds <= 0 {
		return fallback
	}

	return time.Duration(seconds) * time.Second
}

func getIntEnv(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 0 {
		return fallback
	}

	return parsed
}

func getCSVEnv(key string, fallback []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	items := make([]string, 0)
	for _, item := range strings.Split(value, ",") {
		item = strings.TrimSpace(item)
		if item != "" {
			items = append(items, item)
		}
	}

	if len(items) == 0 {
		return fallback
	}

	return items
}

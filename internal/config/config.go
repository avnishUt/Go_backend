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

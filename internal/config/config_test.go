package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	t.Setenv("APP_NAME", "")
	t.Setenv("APP_ENV", "")
	t.Setenv("PORT", "")
	t.Setenv("RATE_LIMIT_ENABLED", "")

	cfg := Load()

	if cfg.AppName != "restaurant-inventory-api" {
		t.Fatalf("expected default app name, got %q", cfg.AppName)
	}
	if cfg.Port != "8080" {
		t.Fatalf("expected default port, got %q", cfg.Port)
	}
	if !cfg.RateLimit.Enabled {
		t.Fatal("expected rate limit to be enabled by default")
	}
	if cfg.Security.MaxBodyBytes != 1048576 {
		t.Fatalf("expected 1MB body limit, got %d", cfg.Security.MaxBodyBytes)
	}
}

func TestLoadOverrides(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("RATE_LIMIT_ENABLED", "false")
	t.Setenv("RATE_LIMIT_REQUESTS_PER_MINUTE", "10")
	t.Setenv("HTTP_MAX_BODY_BYTES", "2048")

	cfg := Load()

	if cfg.Port != "9090" {
		t.Fatalf("expected port override, got %q", cfg.Port)
	}
	if cfg.RateLimit.Enabled {
		t.Fatal("expected rate limit disabled")
	}
	if cfg.RateLimit.RequestsPerMinute != 10 {
		t.Fatalf("expected rate limit 10, got %d", cfg.RateLimit.RequestsPerMinute)
	}
	if cfg.Security.MaxBodyBytes != 2048 {
		t.Fatalf("expected body limit 2048, got %d", cfg.Security.MaxBodyBytes)
	}
}

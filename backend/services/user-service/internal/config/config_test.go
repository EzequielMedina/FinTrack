package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad_DefaultsAndDurations(t *testing.T) {
	// unset potentially set envs
	os.Unsetenv("PORT")
	os.Unsetenv("JWT_EXPIRY")
	os.Unsetenv("JWT_REFRESH_EXPIRY")
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}
	if cfg.Port == "" {
		t.Fatalf("expected default port")
	}
	if cfg.JWTExpiry <= 0 || cfg.RefreshExpiry <= 0 {
		t.Fatalf("expected positive durations")
	}
}

func TestParseDurationEnv_InvalidFallsback(t *testing.T) {
	os.Setenv("X_DURATION", "notaduration")
	d := ParseDurationEnv("X_DURATION", "1h")
	if d <= 0 {
		t.Fatalf("expected fallback positive duration")
	}
}

func TestParseDurationEnv_Valid(t *testing.T) {
	os.Setenv("X_DURATION", "2h")
	d := ParseDurationEnv("X_DURATION", "1h")
	if d != 2*time.Hour {
		t.Fatalf("expected 2h, got %v", d)
	}
}

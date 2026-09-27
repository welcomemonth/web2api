package config

import (
	"testing"
	"time"
)

func TestDefault(t *testing.T) {
	cfg := Default()

	if cfg.Port != 8080 {
		t.Errorf("Port = %d, want 8080", cfg.Port)
	}
	if cfg.DataDir != "data" {
		t.Errorf("DataDir = %q, want %q", cfg.DataDir, "data")
	}
	if cfg.JWTExpiry != 24*time.Hour {
		t.Errorf("JWTExpiry = %v, want 24h", cfg.JWTExpiry)
	}
	if cfg.APIRateLimitRPM != 10 {
		t.Errorf("APIRateLimitRPM = %d, want 10", cfg.APIRateLimitRPM)
	}
	if cfg.AccountWaitTimeout != 30*time.Second {
		t.Errorf("AccountWaitTimeout = %v, want 30s", cfg.AccountWaitTimeout)
	}
}

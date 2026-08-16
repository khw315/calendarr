package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/khw315/calendarr/internal/models"
)

func TestConfigManagerLoadAndSave(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "calendarr.json")

	mgr := NewManager(configPath)
	cfg := mgr.Get()

	if cfg.Timezone != "UTC" {
		t.Errorf("Expected default timezone UTC, got %s", cfg.Timezone)
	}

	// Update config
	cfg.Timezone = "Asia/Jakarta"
	cfg.UseDiscord = true
	cfg.DiscordWebhookURL = "https://discord.com/api/webhooks/test"

	if err := mgr.Save(cfg); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Reload config
	mgr2 := NewManager(configPath)
	cfg2 := mgr2.Get()

	if cfg2.Timezone != "Asia/Jakarta" {
		t.Errorf("Expected timezone Asia/Jakarta, got %s", cfg2.Timezone)
	}
	if !cfg2.UseDiscord {
		t.Errorf("Expected UseDiscord to be true")
	}

	// Ensure file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Errorf("Config file was not created at %s", configPath)
	}
}

func TestDefaultConfigValues(t *testing.T) {
	cfg := models.DefaultConfig()
	if cfg.Language != "EN" {
		t.Errorf("Expected default language EN, got %s", cfg.Language)
	}
	if cfg.HTTPTimeout != 30 {
		t.Errorf("Expected HTTP timeout 30, got %d", cfg.HTTPTimeout)
	}
}

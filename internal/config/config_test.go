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
}

func TestConfigManagerInvalidTimezone(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "calendarr.json")

	invalidJSON := []byte(`{"timezone": "Invalid/Timezone_Name"}`)
	if err := os.WriteFile(configPath, invalidJSON, 0644); err != nil {
		t.Fatalf("Failed to write invalid config file: %v", err)
	}

	mgr := NewManager(configPath)
	cfg := mgr.Get()

	if cfg.Timezone != "UTC" {
		t.Errorf("Expected fallback timezone UTC for invalid timezone, got %s", cfg.Timezone)
	}

	// Test Save with invalid timezone
	cfg.Timezone = "Invalid/Timezone_Name_2"
	if err := mgr.Save(cfg); err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	if cfg.Timezone != "UTC" {
		t.Errorf("Expected fallback to UTC after save, got %s", cfg.Timezone)
	}
}

func TestConfigManagerInvalidJSON(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "calendarr.json")

	if err := os.WriteFile(configPath, []byte(`{invalid json`), 0644); err != nil {
		t.Fatalf("Failed to write corrupt config file: %v", err)
	}

	mgr := NewManager(configPath)
	cfg := mgr.Get()
	if cfg == nil {
		t.Fatalf("Expected fallback config when unmarshal fails")
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

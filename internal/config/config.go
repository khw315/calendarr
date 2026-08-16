package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/khw315/calendarr/internal/constants"
	"github.com/khw315/calendarr/internal/models"
)

type Manager struct {
	mu         sync.RWMutex
	config     *models.Config
	configPath string
}

func NewManager(customPath string) *Manager {
	path := constants.DefaultConfigPath
	if customPath != "" {
		path = customPath
	} else if _, err := os.Stat(path); os.IsNotExist(err) {
		// Fallback for local testing
		if _, err := os.Stat(constants.DefaultLocalConfigPath); err == nil {
			path = constants.DefaultLocalConfigPath
		}
	}

	m := &Manager{
		configPath: path,
		config:     models.DefaultConfig(),
	}

	if err := m.Load(); err != nil {
		log.Printf("⚠️ Config file not found or failed to load at %s, using defaults: %v", path, err)
	}

	return m
}

func (m *Manager) Get() *models.Config {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Return a shallow copy with valid Location pointer
	cfg := *m.config
	return &cfg
}

func (m *Manager) Load() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	data, err := os.ReadFile(m.configPath)
	if err != nil {
		return err
	}

	cfg := models.DefaultConfig()
	if err := json.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("failed to parse config JSON: %w", err)
	}

	// Resolve Timezone Location
	loc, err := time.LoadLocation(cfg.Timezone)
	if err != nil {
		log.Printf("⚠️ Invalid timezone %s, defaulting to UTC: %v", cfg.Timezone, err)
		loc = time.UTC
		cfg.Timezone = "UTC"
	}
	cfg.TimezoneLocation = loc

	m.config = cfg
	log.Printf("✅ Configuration successfully loaded from %s (Timezone: %s)", m.configPath, cfg.Timezone)
	return nil
}

func (m *Manager) Save(updated *models.Config) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Ensure directory exists
	dir := filepath.Dir(m.configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Resolve Timezone Location
	loc, err := time.LoadLocation(updated.Timezone)
	if err != nil {
		loc = time.UTC
		updated.Timezone = "UTC"
	}
	updated.TimezoneLocation = loc

	data, err := json.MarshalIndent(updated, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(m.configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	m.config = updated
	log.Printf("✅ Configuration updated and saved to %s", m.configPath)
	return nil
}

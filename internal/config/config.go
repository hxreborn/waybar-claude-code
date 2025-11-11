package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds configuration values for the module
type Config struct {
	// Interval is how often to poll ccusage (default: 300 seconds / 5 minutes)
	Interval time.Duration

	// Debug enables verbose logging to stderr
	Debug bool
}

// Load reads configuration from environment variables
func Load() *Config {
	cfg := &Config{
		Interval: 300 * time.Second, // Default: 5 minutes
		Debug:    false,
	}

	// Parse CLAUDE_INTERVAL_SEC
	if intervalStr := os.Getenv("CLAUDE_INTERVAL_SEC"); intervalStr != "" {
		if interval, err := strconv.Atoi(intervalStr); err == nil && interval > 0 {
			cfg.Interval = time.Duration(interval) * time.Second
		}
	}

	// Parse CLAUDE_DEBUG
	if debugStr := os.Getenv("CLAUDE_DEBUG"); debugStr != "" {
		cfg.Debug = debugStr == "1" || debugStr == "true" || debugStr == "TRUE"
	}

	return cfg
}

package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Interval time.Duration
}

func Load() *Config {
	cfg := &Config{
		Interval: 300 * time.Second,
	}

	if intervalStr := os.Getenv("CLAUDE_INTERVAL_SEC"); intervalStr != "" {
		if interval, err := strconv.Atoi(intervalStr); err == nil && interval > 0 {
			cfg.Interval = time.Duration(interval) * time.Second
		}
	}

	return cfg
}

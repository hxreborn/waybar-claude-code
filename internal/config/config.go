package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Interval time.Duration
	Animate  bool
}

func Load() *Config {
	cfg := &Config{
		Interval: 300 * time.Second,
		Animate:  false,
	}

	if intervalStr := os.Getenv("CLAUDE_INTERVAL_SEC"); intervalStr != "" {
		if interval, err := strconv.Atoi(intervalStr); err == nil && interval > 0 {
			cfg.Interval = time.Duration(interval) * time.Second
		}
	}

	if animateStr := os.Getenv("CLAUDE_ANIMATE"); animateStr != "" {
		cfg.Animate = animateStr == "1" || animateStr == "true"
	}

	return cfg
}

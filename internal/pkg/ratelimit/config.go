package ratelimit

import (
	"os"
	"strconv"
	"strings"
)

// Config holds per-IP limits (fixed 1-minute window per UTC minute).
type Config struct {
	Enabled    bool
	GeneralRPM int // max requests per IP per minute (non-ranking routes)
	RankingRPM int // max requests per IP per minute for /api/ranking/*
}

// ConfigFromEnv builds config with defaults. Set RATE_LIMIT_DISABLED=1/true to turn off.
func ConfigFromEnv() Config {
	disabled := strings.EqualFold(strings.TrimSpace(os.Getenv("RATE_LIMIT_DISABLED")), "1") ||
		strings.EqualFold(strings.TrimSpace(os.Getenv("RATE_LIMIT_DISABLED")), "true")

	c := Config{
		Enabled:    !disabled,
		GeneralRPM: 120,
		RankingRPM: 40,
	}
	if v := os.Getenv("RATE_LIMIT_GENERAL_RPM"); v != "" {
		if n, err := strconv.Atoi(strings.TrimSpace(v)); err == nil && n > 0 {
			c.GeneralRPM = n
		}
	}
	if v := os.Getenv("RATE_LIMIT_RANKING_RPM"); v != "" {
		if n, err := strconv.Atoi(strings.TrimSpace(v)); err == nil && n > 0 {
			c.RankingRPM = n
		}
	}
	return c
}

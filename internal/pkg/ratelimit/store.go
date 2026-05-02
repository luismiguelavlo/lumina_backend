package ratelimit

import (
	"strings"
	"sync"
	"time"
)

// Store enforces per-IP per-minute counters for general vs ranking routes.
type Store struct {
	cfg Config
	mu  sync.Mutex
	m   map[string]*ipEntry
}

type ipEntry struct {
	minuteGen, minuteRank int64
	countGen, countRank   int
	lastSeen              time.Time
}

// NewStore creates an empty rate limit store. StartSweeper should be called once.
func NewStore(cfg Config) *Store {
	return &Store{cfg: cfg, m: make(map[string]*ipEntry)}
}

// Allow reports whether a request from ip for path should proceed (and records it).
func (s *Store) Allow(path, ip string) bool {
	if !s.cfg.Enabled {
		return true
	}
	ip = strings.TrimSpace(ip)
	if ip == "" {
		ip = "unknown"
	}
	now := time.Now().UTC()
	min := now.Unix() / 60

	s.mu.Lock()
	defer s.mu.Unlock()

	e := s.m[ip]
	if e == nil {
		e = &ipEntry{}
		s.m[ip] = e
	}
	e.lastSeen = now

	if strings.HasPrefix(path, "/api/ranking") {
		if e.minuteRank != min {
			e.minuteRank = min
			e.countRank = 0
		}
		if e.countRank >= s.cfg.RankingRPM {
			return false
		}
		e.countRank++
		return true
	}

	if e.minuteGen != min {
		e.minuteGen = min
		e.countGen = 0
	}
	if e.countGen >= s.cfg.GeneralRPM {
		return false
	}
	e.countGen++
	return true
}

// ShouldSkip returns true for paths exempt from rate limiting.
func (s *Store) ShouldSkip(path string) bool {
	if !s.cfg.Enabled {
		return true
	}
	// Liveness/readiness probes should not consume quota.
	if path == "/health" || strings.HasPrefix(path, "/health/") {
		return true
	}
	return false
}

// StartSweeper removes stale IP entries in the background.
func (s *Store) StartSweeper(interval, maxAge time.Duration) {
	if !s.cfg.Enabled {
		return
	}
	if interval <= 0 {
		interval = 5 * time.Minute
	}
	if maxAge <= 0 {
		maxAge = 30 * time.Minute
	}
	go func() {
		t := time.NewTicker(interval)
		defer t.Stop()
		for range t.C {
			s.sweep(maxAge)
		}
	}()
}

func (s *Store) sweep(maxAge time.Duration) {
	cutoff := time.Now().UTC().Add(-maxAge)
	s.mu.Lock()
	defer s.mu.Unlock()
	for ip, e := range s.m {
		if e.lastSeen.Before(cutoff) {
			delete(s.m, ip)
		}
	}
}

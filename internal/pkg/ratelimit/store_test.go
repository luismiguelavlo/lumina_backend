package ratelimit

import (
	"testing"
	"time"
)

func TestStore_DisabledAlwaysAllows(t *testing.T) {
	s := NewStore(Config{Enabled: false})
	for i := 0; i < 1000; i++ {
		if !s.Allow("/api/foo", "1.2.3.4") {
			t.Fatalf("want allow at %d", i)
		}
	}
}

func TestStore_GeneralRPM(t *testing.T) {
	s := NewStore(Config{Enabled: true, GeneralRPM: 3, RankingRPM: 99})
	ip := "10.0.0.1"
	for i := 0; i < 3; i++ {
		if !s.Allow("/api/books", ip) {
			t.Fatalf("want allow %d", i)
		}
	}
	if s.Allow("/api/books", ip) {
		t.Fatal("want deny on 4th")
	}
}

func TestStore_RankingSeparateBucket(t *testing.T) {
	s := NewStore(Config{Enabled: true, GeneralRPM: 2, RankingRPM: 2})
	ip := "10.0.0.2"
	if !s.Allow("/api/ranking/top3", ip) || !s.Allow("/api/ranking/top3", ip) {
		t.Fatal("want two ranking allows")
	}
	if s.Allow("/api/ranking/leaderboard", ip) {
		t.Fatal("want deny 3rd ranking")
	}
	// general bucket untouched
	if !s.Allow("/api/books", ip) || !s.Allow("/api/books", ip) {
		t.Fatal("want two general allows")
	}
	if s.Allow("/api/books", ip) {
		t.Fatal("want deny 3rd general")
	}
}

func TestStore_SweepRemovesStale(t *testing.T) {
	s := NewStore(Config{Enabled: true, GeneralRPM: 100, RankingRPM: 100})
	s.Allow("/api/a", "9.9.9.9")
	s.mu.Lock()
	s.m["9.9.9.9"].lastSeen = time.Now().UTC().Add(-2 * time.Hour)
	s.mu.Unlock()
	s.sweep(1 * time.Hour)
	s.mu.Lock()
	_, ok := s.m["9.9.9.9"]
	s.mu.Unlock()
	if ok {
		t.Fatal("expected ip removed")
	}
}

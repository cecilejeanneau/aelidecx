package main

import (
	"database/sql"
	"regexp"
	"sync"
	"time"

	"github.com/microcosm-cc/bluemonday"
)

var db *sql.DB
var apiKeys []string
var htmlPolicy *bluemonday.Policy
var searchQueryPattern = regexp.MustCompile(`^[\p{L}\p{N}\s\-_."']+$`)
var startedAt = time.Now().UTC()

type securityStats struct {
	mu           sync.Mutex
	APIRequests  int64 `json:"api_requests"`
	RateLimited  int64 `json:"rate_limited"`
	AuthFailures int64 `json:"auth_failures"`
	AuthLocked   int64 `json:"auth_locked"`
	Uploads      int64 `json:"uploads"`
}

var secStats = &securityStats{}

type observabilityStats struct {
	mu                  sync.Mutex
	APIRequests         int64                  `json:"api_requests"`
	Responses2xx        int64                  `json:"responses_2xx"`
	Responses4xx        int64                  `json:"responses_4xx"`
	Responses5xx        int64                  `json:"responses_5xx"`
	TotalResponseTimeMS float64                `json:"-"`
	RecentLatenciesMS   []float64              `json:"-"`
	Routes              map[string]*routeStats `json:"-"`
}

type routeStats struct {
	Count               int64   `json:"count"`
	Responses2xx        int64   `json:"responses_2xx"`
	Responses4xx        int64   `json:"responses_4xx"`
	Responses5xx        int64   `json:"responses_5xx"`
	TotalResponseTimeMS float64 `json:"-"`
}

var obsStats = &observabilityStats{
	Routes: make(map[string]*routeStats),
}

type authAttemptState struct {
	Failures    int
	LockedUntil time.Time
}

var authState = struct {
	mu sync.Mutex
	m  map[string]*authAttemptState
}{
	m: make(map[string]*authAttemptState),
}

var rateLimiter = struct {
	mu    sync.Mutex
	hits  map[string][]time.Time
	limit int
	win   time.Duration
}{
	hits:  make(map[string][]time.Time),
	limit: 120,
	win:   time.Minute,
}

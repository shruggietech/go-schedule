package remote

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type limiterSet struct {
	mu      sync.Mutex
	items   map[string]limiterEntry
	limit   rate.Limit
	burst   int
	maximum int
	now     func() time.Time
	maxIdle time.Duration
}

type limiterEntry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

func newLimiterSet(perSecond float64, burst, maximum int) *limiterSet {
	return &limiterSet{items: make(map[string]limiterEntry), limit: rate.Limit(perSecond), burst: burst, maximum: maximum, now: time.Now, maxIdle: 10 * time.Minute}
}

func (s *limiterSet) allow(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.now()
	entry, exists := s.items[key]
	if !exists {
		if len(s.items) >= s.maximum {
			for candidate, retained := range s.items {
				if now.Sub(retained.lastSeen) >= s.maxIdle {
					delete(s.items, candidate)
				}
			}
			if len(s.items) >= s.maximum {
				return false
			}
		}
		entry = limiterEntry{limiter: rate.NewLimiter(s.limit, s.burst)}
	}
	entry.lastSeen = now
	s.items[key] = entry
	return entry.limiter.AllowN(now, 1)
}

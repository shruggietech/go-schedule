package remote

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type limiterSet struct {
	mu      sync.Mutex
	items   map[string]*rate.Limiter
	limit   rate.Limit
	burst   int
	maximum int
}

func newLimiterSet(perSecond float64, burst, maximum int) *limiterSet {
	return &limiterSet{items: make(map[string]*rate.Limiter), limit: rate.Limit(perSecond), burst: burst, maximum: maximum}
}

func (s *limiterSet) allow(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	item := s.items[key]
	if item == nil {
		if len(s.items) >= s.maximum {
			return false
		}
		item = rate.NewLimiter(s.limit, s.burst)
		s.items[key] = item
	}
	return item.AllowN(time.Now(), 1)
}

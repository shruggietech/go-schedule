package remote

import "sync"

type streamSet struct {
	mu              sync.Mutex
	counts          map[string]int
	total           int
	maximumTotal    int
	maximumPerActor int
}

func newStreamSet(maximumTotal, maximumPerActor int) *streamSet {
	return &streamSet{counts: make(map[string]int), maximumTotal: maximumTotal, maximumPerActor: maximumPerActor}
}

func (s *streamSet) acquire(actorID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.total >= s.maximumTotal || s.counts[actorID] >= s.maximumPerActor {
		return false
	}
	s.total++
	s.counts[actorID]++
	return true
}

func (s *streamSet) release(actorID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.counts[actorID] <= 0 {
		return
	}
	s.total--
	s.counts[actorID]--
	if s.counts[actorID] == 0 {
		delete(s.counts, actorID)
	}
}

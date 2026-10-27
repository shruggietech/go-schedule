package incident

import (
	"context"
	"net/http"
	"sort"
	"sync"
	"time"
)

// Active identifies one pending authored route and its current progress stage.
type Active struct {
	ID        uint64 `json:"id"`
	Method    string `json:"method"`
	Route     string `json:"route"`
	Stage     string `json:"stage"`
	ElapsedMS int64  `json:"elapsed_ms"`
	start     time.Time
}

// Tracker bounds active request evidence and counts omitted observations.
type Tracker struct {
	mu      sync.Mutex
	next    uint64
	active  map[uint64]Active
	omitted int
}
type trackingKey struct{}
type trackingValue struct {
	tracker *Tracker
	id      uint64
}

func (t *Tracker) begin(method, route string) (uint64, func()) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.active == nil {
		t.active = make(map[uint64]Active)
	}
	if len(t.active) >= 64 {
		t.omitted++
		return 0, func() {}
	}
	switch method {
	case "GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS":
	default:
		method = "other"
	}
	t.next++
	id := t.next
	t.active[id] = Active{ID: id, Method: method, Route: route, Stage: "authorization", start: time.Now()}
	return id, func() { t.mu.Lock(); delete(t.active, id); t.mu.Unlock() }
}

// Begin tracks a trusted authored route until the returned completion is called.
func (t *Tracker) Begin(method, authoredRoute string) func() {
	_, end := t.begin(method, authoredRoute)
	return end
}

// Snapshot returns current requests and the cumulative omitted-request count.
func (t *Tracker) Snapshot() ([]Active, int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	active := make([]Active, 0, len(t.active))
	for _, a := range t.active {
		a.ElapsedMS = time.Since(a.start).Milliseconds()
		active = append(active, a)
	}
	sort.Slice(active, func(i, j int) bool { return active[i].ID < active[j].ID })
	return active, t.omitted
}

// Stage changes only an authored category on an explicitly tracked request.
func Stage(ctx context.Context, stage string) {
	value, ok := ctx.Value(trackingKey{}).(trackingValue)
	if !ok || value.id == 0 {
		return
	}
	switch stage {
	case "resolve_actor", "actor", "identity", "audit_intent", "handler", "audit_completion":
	default:
		return
	}
	t := value.tracker
	t.mu.Lock()
	defer t.mu.Unlock()
	a, ok := t.active[value.id]
	if ok {
		a.Stage = stage
		t.active[value.id] = a
	}
}

// Handler receives route patterns from the server's authored mux, never raw URLs.
func (t *Tracker) Handler(next http.Handler, route func(*http.Request) string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pattern := route(r)
		if pattern == "" {
			pattern = "unrecognized"
		}
		id, end := t.begin(r.Method, pattern)
		defer end()
		ctx := context.WithValue(r.Context(), trackingKey{}, trackingValue{t, id})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

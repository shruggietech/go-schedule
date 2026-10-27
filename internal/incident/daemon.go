package incident

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// MaxSnapshots limits incident collection for one disposable daemon.
const MaxSnapshots = 8

// Channel owns the private request/response directory for one daemon.
type Channel struct {
	*Capture
	mu   sync.Mutex
	next int
}

// Request binds one incident to the original owned session and target PID.
type Request struct {
	ID      string `json:"id"`
	Session string `json:"session"`
	PID     int    `json:"pid"`
	Phase   string `json:"phase"`
}

// Pool reports nonblocking connection-pool statistics with explicit time units.
type Pool struct {
	InUse           int   `json:"in_use"`
	Idle            int   `json:"idle"`
	OpenConnections int   `json:"open_connections"`
	WaitCount       int64 `json:"wait_count"`
	WaitDurationMS  int64 `json:"wait_duration_ms"`
}

// Snapshot describes live handler, pool and private stack evidence.
type Snapshot struct {
	CapturedAt string   `json:"captured_at"`
	Version    int      `json:"version"`
	Request    Request  `json:"request"`
	Target     Target   `json:"target"`
	Pool       Pool     `json:"pool"`
	Active     []Active `json:"active"`
	Omitted    int      `json:"omitted"`
	Stack      Artifact `json:"stack"`
}

// NewChannel creates an exclusive private directory beneath the owned parent.
func NewChannel(parent string) (*Channel, error) {
	c, err := New(parent)
	if err != nil {
		return nil, err
	}
	return &Channel{Capture: c}, nil
}
func openChannel(dir string) (*Capture, error) {
	id := strings.TrimPrefix(filepath.Base(dir), "incident-")
	c := &Capture{Dir: dir, ID: id, ownedDir: dir, ownedID: id}
	if _, err := c.path("request.json"); err != nil {
		return nil, err
	}
	return c, nil
}

// ValidateWorkspace restricts daemon collection to the selected disposable fixture.
func ValidateWorkspace(dir, dataDir string) error {
	c, err := openChannel(dir)
	if err != nil {
		return err
	}
	root := filepath.Dir(filepath.Dir(dataDir))
	if !samePath(filepath.Dir(c.Dir), filepath.Join(root, "evidence")) {
		return errors.New("incident directory is outside the owned fixture evidence")
	}
	if err := canonicalDirectory(root); err != nil {
		return err
	}
	manifest := filepath.Join(root, "evidence", "ownership.json")
	info, err := os.Lstat(manifest)
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() || info.Size() > 4096 {
		return errors.New("owned fixture ownership record is invalid")
	}
	file, err := os.Open(manifest)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()
	var ownership struct {
		RunID    string `json:"run_id"`
		Root     string `json:"root"`
		Runtime  string `json:"runtime"`
		Evidence string `json:"evidence"`
	}
	if json.NewDecoder(file).Decode(&ownership) != nil || uuid.Validate(ownership.RunID) != nil || !samePath(ownership.Root, root) || !samePath(ownership.Runtime, filepath.Join(root, "runtime")) || !samePath(ownership.Evidence, filepath.Join(root, "evidence")) || filepath.Base(root) != "gosched-"+ownership.RunID || !samePath(dataDir, filepath.Join(root, "runtime", "state")) {
		return errors.New("incident target is not the owned fixture workspace")
	}
	return nil
}

// Validate rejects responses from a different session, process or format version.
func (c *Channel) Validate(snapshot Snapshot, pid int) error {
	if snapshot.Version != 1 || snapshot.Request.Session != c.ID || snapshot.Request.PID != pid || snapshot.Target.PID != pid || uuid.Validate(snapshot.Request.ID) != nil {
		return errors.New("incident response does not match the owned daemon")
	}
	return nil
}

// Collect is independent of the API and honors its own bounded caller context.
func (c *Channel) Collect(ctx context.Context, pid int, phase string) (Snapshot, error) {
	if err := ctx.Err(); err != nil {
		return Snapshot{}, err
	}
	if pid <= 0 {
		return Snapshot{}, errors.New("owned daemon PID is missing")
	}
	c.mu.Lock()
	c.next++
	slot := c.next
	c.mu.Unlock()
	if slot > MaxSnapshots {
		return Snapshot{}, errors.New("owned daemon incident allowance exhausted")
	}
	request := Request{ID: uuid.NewString(), Session: c.ID, PID: pid, Phase: phase}
	data, err := json.Marshal(request)
	if err != nil {
		return Snapshot{}, err
	}
	if len(data) > 1024 {
		return Snapshot{}, errors.New("incident request exceeds one KiB")
	}
	if err := c.Publish(fmt.Sprintf("%d.request.json", slot), data); err != nil {
		return Snapshot{}, err
	}
	tick := time.NewTicker(20 * time.Millisecond)
	defer tick.Stop()
	for {
		data, err = c.Read(fmt.Sprintf("%d.response.json", slot), MetadataLimit)
		if err == nil {
			var snapshot Snapshot
			if err := json.Unmarshal(data, &snapshot); err != nil {
				return Snapshot{}, errors.New("incident response is invalid")
			}
			if snapshot.Request.ID != request.ID {
				return Snapshot{}, errors.New("incident response request identity changed")
			}
			if err := c.Validate(snapshot, pid); err != nil {
				return Snapshot{}, err
			}
			return snapshot, nil
		}
		if !os.IsNotExist(err) {
			return Snapshot{}, err
		}
		select {
		case <-ctx.Done():
			return Snapshot{}, ctx.Err()
		case <-tick.C:
		}
	}
}

// Serve collects only private metadata and bounded stacks, never SQL or request content.
// It is disabled unless the parent explicitly supplies its exclusively created directory.
func Serve(ctx context.Context, dir string, tracker *Tracker, pool func() sql.DBStats) error {
	c, err := openChannel(dir)
	if err != nil {
		return err
	}
	executable, err := os.Executable()
	if err != nil {
		return err
	}
	executable, err = filepath.EvalSymlinks(executable)
	if err != nil {
		return err
	}
	hash, err := FileHash(executable)
	if err != nil {
		return err
	}
	target := Target{PID: os.Getpid(), Executable: executable, SHA256: hash}
	tick := time.NewTicker(100 * time.Millisecond)
	defer tick.Stop()
	seen := [MaxSnapshots]bool{}
	for {
		for slot := 1; slot <= MaxSnapshots; slot++ {
			if seen[slot-1] {
				continue
			}
			data, readErr := c.Read(fmt.Sprintf("%d.request.json", slot), 1024)
			if os.IsNotExist(readErr) {
				continue
			}
			if readErr != nil {
				return readErr
			}
			seen[slot-1] = true
			var request Request
			if json.Unmarshal(data, &request) != nil || uuid.Validate(request.ID) != nil || request.Session != c.ID || request.PID != target.PID {
				return errors.New("incident request does not match this owned daemon")
			}
			buffer := make([]byte, 64<<10)
			n := runtime.Stack(buffer, true)
			name := fmt.Sprintf("%d.goroutines.txt", slot)
			if err := c.Publish(name, buffer[:n]); err != nil {
				return err
			}
			artifact := c.Inspect(name, 64<<10)
			if n == len(buffer) && artifact.Status == "collected" {
				artifact.Status = "truncated"
			}
			active, omitted := tracker.Snapshot()
			stats := pool()
			snapshot := Snapshot{Version: 1, CapturedAt: time.Now().UTC().Format(time.RFC3339Nano), Request: request, Target: target, Pool: Pool{InUse: stats.InUse, Idle: stats.Idle, OpenConnections: stats.OpenConnections, WaitCount: stats.WaitCount, WaitDurationMS: stats.WaitDuration.Milliseconds()}, Active: active, Omitted: omitted, Stack: artifact}
			encoded, err := json.Marshal(snapshot)
			if err != nil {
				return err
			}
			if err := c.Publish(fmt.Sprintf("%d.response.json", slot), encoded); err != nil {
				return err
			}
		}
		select {
		case <-ctx.Done():
			return nil
		case <-tick.C:
		}
	}
}

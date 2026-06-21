package logbus

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"strconv"
	"sync/atomic"

	"github.com/shruggietech/go-schedule/internal/domain"
)

// Publisher receives log records for live streaming (the event broker implements
// it). Kept as an interface so the handler is testable without a real broker.
type Publisher interface {
	PublishLog(domain.LogRecord)
}

// seq provides process-wide monotonic record IDs (sufficient for GUI dedupe).
var seq atomic.Uint64

// Handler is an slog.Handler that tees each record into a domain.LogRecord and
// fans it out to the ring (for GET /v1/logs), the rotating file (durable), and
// the publisher (live stream). It delegates JSON formatting of the file line to
// itself so the on-disk format is stable regardless of the slog handler used
// elsewhere.
type Handler struct {
	level  slog.Leveler
	ring   *Ring
	file   io.Writer
	pub    Publisher
	source string
	attrs  []slog.Attr
	taskID string
	runID  string
}

// NewHandler builds a teeing handler. ring and file may be nil to disable a sink;
// pub may be nil before the broker exists. level defaults to Info if nil.
func NewHandler(level slog.Leveler, ring *Ring, file io.Writer, pub Publisher) *Handler {
	if level == nil {
		level = slog.LevelInfo
	}
	return &Handler{level: level, ring: ring, file: file, pub: pub}
}

// Enabled reports whether records at lvl should be handled.
func (h *Handler) Enabled(_ context.Context, lvl slog.Level) bool {
	return lvl >= h.level.Level()
}

// Handle converts the record and fans it out to all configured sinks.
func (h *Handler) Handle(_ context.Context, r slog.Record) error {
	rec := domain.LogRecord{
		ID:       strconv.FormatUint(seq.Add(1), 10),
		Time:     r.Time,
		Severity: severityFor(r.Level),
		Source:   h.source,
		Message:  r.Message,
		TaskID:   h.taskID,
		RunID:    h.runID,
		Attrs:    map[string]any{},
	}
	// Pre-seed with handler-scoped attrs, then overlay record attrs.
	for _, a := range h.attrs {
		applyAttr(&rec, a)
	}
	r.Attrs(func(a slog.Attr) bool {
		applyAttr(&rec, a)
		return true
	})
	if len(rec.Attrs) == 0 {
		rec.Attrs = nil
	}

	if h.ring != nil {
		h.ring.Add(rec)
	}
	if h.file != nil {
		if line, err := marshalLine(rec); err == nil {
			_, _ = h.file.Write(line)
		}
	}
	if h.pub != nil {
		h.pub.PublishLog(rec)
	}
	return nil
}

// WithAttrs returns a handler with additional scoped attributes.
func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	nh := h.clone()
	for _, a := range attrs {
		// Correlation attrs are promoted to typed fields; the rest accumulate.
		switch a.Key {
		case "task", "task_id":
			nh.taskID = a.Value.String()
		case "run", "run_id":
			nh.runID = a.Value.String()
		default:
			nh.attrs = append(nh.attrs, a)
		}
	}
	return nh
}

// WithGroup uses the group name as the record source/component.
func (h *Handler) WithGroup(name string) slog.Handler {
	nh := h.clone()
	if name != "" {
		nh.source = name
	}
	return nh
}

func (h *Handler) clone() *Handler {
	cp := *h
	cp.attrs = append([]slog.Attr(nil), h.attrs...)
	return &cp
}

// applyAttr routes an attr to a typed field or the generic Attrs map.
func applyAttr(rec *domain.LogRecord, a slog.Attr) {
	switch a.Key {
	case "task", "task_id":
		rec.TaskID = a.Value.String()
	case "run", "run_id":
		rec.RunID = a.Value.String()
	default:
		rec.Attrs[a.Key] = a.Value.Any()
	}
}

// severityFor maps an slog level to a LogRecord severity.
func severityFor(lvl slog.Level) domain.AlertSeverity {
	switch {
	case lvl >= slog.LevelError:
		return domain.SeverityError
	case lvl >= slog.LevelWarn:
		return domain.SeverityWarning
	default:
		return domain.SeverityInfo
	}
}

// marshalLine renders a record as a JSONL line with reserved keys plus inlined
// attrs (the durable on-disk format; see contracts/log-file.md).
func marshalLine(rec domain.LogRecord) ([]byte, error) {
	m := map[string]any{
		"id":       rec.ID,
		"time":     rec.Time.UTC().Format("2006-01-02T15:04:05.000Z07:00"),
		"severity": string(rec.Severity),
		"message":  rec.Message,
	}
	if rec.Source != "" {
		m["source"] = rec.Source
	}
	if rec.TaskID != "" {
		m["task_id"] = rec.TaskID
	}
	if rec.RunID != "" {
		m["run_id"] = rec.RunID
	}
	for k, v := range rec.Attrs {
		if _, reserved := m[k]; !reserved {
			m[k] = v
		}
	}
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return append(b, '\n'), nil
}

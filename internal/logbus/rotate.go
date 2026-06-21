package logbus

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// RotatingWriter is a size-based rotating io.Writer for the JSONL log file. It is
// intentionally small and dependency-free (the daemon is the single writer). When
// the active file would exceed maxBytes, it is rotated: foo.log -> foo.log.1 ->
// ... -> foo.log.(maxFiles-1), and the oldest is discarded. Concurrency-safe.
type RotatingWriter struct {
	mu       sync.Mutex
	path     string
	maxBytes int64
	maxFiles int
	f        *os.File
	size     int64
}

// NewRotatingWriter opens (creating dirs as needed) the log file at path. maxBytes
// and maxFiles bound disk usage; non-positive values fall back to 10 MiB / 5.
func NewRotatingWriter(path string, maxBytes int64, maxFiles int) (*RotatingWriter, error) {
	if maxBytes <= 0 {
		maxBytes = 10 << 20
	}
	if maxFiles <= 0 {
		maxFiles = 5
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("logbus: create log dir: %w", err)
	}
	w := &RotatingWriter{path: path, maxBytes: maxBytes, maxFiles: maxFiles}
	if err := w.open(); err != nil {
		return nil, err
	}
	return w, nil
}

func (w *RotatingWriter) open() error {
	f, err := os.OpenFile(w.path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return fmt.Errorf("logbus: open log file: %w", err)
	}
	info, err := f.Stat()
	if err != nil {
		_ = f.Close()
		return fmt.Errorf("logbus: stat log file: %w", err)
	}
	w.f = f
	w.size = info.Size()
	return nil
}

// Write appends p, rotating first if it would exceed the size cap.
func (w *RotatingWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.size+int64(len(p)) > w.maxBytes && w.size > 0 {
		if err := w.rotate(); err != nil {
			return 0, err
		}
	}
	n, err := w.f.Write(p)
	w.size += int64(n)
	return n, err
}

// rotate closes the active file, shifts the numbered backups, and reopens.
func (w *RotatingWriter) rotate() error {
	if err := w.f.Close(); err != nil {
		return fmt.Errorf("logbus: close for rotate: %w", err)
	}
	// Discard the oldest, then shift each backup up by one.
	oldest := fmt.Sprintf("%s.%d", w.path, w.maxFiles-1)
	_ = os.Remove(oldest)
	for i := w.maxFiles - 1; i >= 1; i-- {
		src := w.path
		if i > 1 {
			src = fmt.Sprintf("%s.%d", w.path, i-1)
		}
		dst := fmt.Sprintf("%s.%d", w.path, i)
		_ = os.Rename(src, dst) // best-effort; missing source is fine
	}
	return w.open()
}

// Close closes the active file.
func (w *RotatingWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.f == nil {
		return nil
	}
	return w.f.Close()
}

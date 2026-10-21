package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"
)

// FilesystemWatcherRevision identifies the durable definition, independently of
// joined task display names and the runtime's health transition timestamp.
func FilesystemWatcherRevision(watcher FilesystemWatcher) string {
	hash := sha256.New()
	fields := []string{watcher.ID, watcher.Name, string(watcher.Kind), watcher.Path, watcher.Pattern, strconv.FormatBool(watcher.Recursive), watcher.Debounce.String(), watcher.Stability.String(), watcher.TargetTaskID, strconv.FormatBool(watcher.Enabled), watcher.CreatedAt.UTC().Format(time.RFC3339Nano), watcher.UpdatedAt.UTC().Format(time.RFC3339Nano)}
	for _, field := range fields {
		// Length prefixes keep arbitrary user strings unambiguous.
		_, _ = hash.Write([]byte(strconv.Itoa(len(field)) + ":" + field))
	}
	return hex.EncodeToString(hash.Sum(nil))
}

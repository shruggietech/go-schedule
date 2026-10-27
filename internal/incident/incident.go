// Package incident collects opt-in evidence for disposable owned QA processes.
// Collection outcomes never replace the operation being diagnosed.
package incident

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/google/uuid"
)

// MetadataLimit bounds one structured incident file.
const MetadataLimit = 128 << 10

// Target identifies the exact owned executable and child process.
type Target struct {
	PID        int    `json:"pid"`
	Executable string `json:"executable,omitempty"`
	SHA256     string `json:"sha256,omitempty"`
}

// Original retains the operation outcome independently of collection.
type Original struct {
	NativeExit         int  `json:"native_exit"`
	InspectionComplete bool `json:"inspection_complete"`
	Canceled           bool `json:"canceled"`
}

// Collection records a collector's status without changing the original result.
type Collection struct {
	Status     string `json:"status"`
	Detail     string `json:"detail,omitempty"`
	DurationMS int64  `json:"duration_ms"`
}

// Artifact inventories one bounded local diagnostic file.
type Artifact struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Bytes  int64  `json:"bytes"`
	SHA256 string `json:"sha256,omitempty"`
}

// Record is the versioned owned incident envelope.
type Record struct {
	Kind       string     `json:"kind"`
	Version    int        `json:"version"`
	ID         string     `json:"id"`
	Category   string     `json:"category"`
	Phase      string     `json:"phase"`
	Target     Target     `json:"target"`
	Original   Original   `json:"original"`
	Collection Collection `json:"collection"`
	Artifacts  []Artifact `json:"artifacts"`
}

// Capture names one exclusively created directory. It cannot follow links.
type Capture struct {
	Dir, ID           string
	ownedDir, ownedID string
}

func canonicalDirectory(dir string) error {
	if !filepath.IsAbs(dir) {
		return errors.New("incident directory must be absolute")
	}
	for candidate := filepath.Clean(dir); ; candidate = filepath.Dir(candidate) {
		info, err := os.Lstat(candidate)
		if err != nil {
			return err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return errors.New("incident directory must not traverse a link")
		}
		if filepath.Dir(candidate) == candidate {
			break
		}
	}
	info, err := os.Stat(dir)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return errors.New("incident directory is not a directory")
	}
	return nil
}
func samePath(a, b string) bool {
	if runtime.GOOS == "windows" {
		return strings.EqualFold(filepath.Clean(a), filepath.Clean(b))
	}
	return filepath.Clean(a) == filepath.Clean(b)
}

// New requires an existing private parent selected by the owner of the fixture.
func New(parent string) (*Capture, error) {
	if err := canonicalDirectory(parent); err != nil {
		return nil, err
	}
	parent, err := filepath.EvalSymlinks(parent)
	if err != nil {
		return nil, err
	}
	id := uuid.NewString()
	dir := filepath.Join(parent, "incident-"+id)
	if err := os.Mkdir(dir, 0700); err != nil {
		return nil, err
	}
	return &Capture{Dir: dir, ID: id, ownedDir: dir, ownedID: id}, nil
}
func (c *Capture) path(name string) (string, error) {
	if c == nil || !samePath(c.Dir, c.ownedDir) || c.ID != c.ownedID || uuid.Validate(c.ID) != nil || filepath.Base(c.Dir) != "incident-"+c.ID {
		return "", errors.New("incident identity is invalid")
	}
	if err := canonicalDirectory(c.Dir); err != nil {
		return "", err
	}
	if name == "" || name == "." || name == ".." || filepath.Base(name) != name || strings.ContainsAny(name, "/\\:") {
		return "", errors.New("incident artifact name must be a single filename")
	}
	return filepath.Join(c.Dir, name), nil
}

// Path resolves a single artifact name inside the unchanged capture directory.
func (c *Capture) Path(name string) (string, error) { return c.path(name) }

// Read accepts only a regular owned artifact and limits memory before decoding.
func (c *Capture) Read(name string, limit int64) ([]byte, error) {
	path, err := c.path(name)
	if err != nil {
		return nil, err
	}
	info, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() || info.Size() > limit {
		return nil, errors.New("incident artifact is not a bounded regular file")
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()
	actual, err := file.Stat()
	if err != nil || !os.SameFile(info, actual) {
		return nil, errors.New("incident artifact identity changed")
	}
	data, err := io.ReadAll(io.LimitReader(file, limit+1))
	if err != nil {
		return nil, err
	}
	if int64(len(data)) > limit {
		return nil, errors.New("incident artifact exceeds its byte allowance")
	}
	return data, nil
}

// Publish makes a completed file visible atomically without replacing a destination.
func (c *Capture) Publish(name string, data []byte) error {
	temp := "." + name + ".pending"
	if err := c.Write(temp, data); err != nil {
		return err
	}
	source, err := c.path(temp)
	if err != nil {
		return err
	}
	target, err := c.path(name)
	if err != nil {
		return err
	}
	// Hard-link creation is atomic and fails when the destination already exists.
	err = os.Link(source, target)
	return errors.Join(err, os.Remove(source))
}

// Write never overwrites an existing artifact and caps authored metadata.
func (c *Capture) Write(name string, data []byte) error {
	if len(data) > MetadataLimit {
		return errors.New("incident metadata exceeds 128 KiB")
	}
	path, err := c.path(name)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	_, writeErr := file.Write(data)
	return errors.Join(writeErr, file.Close())
}

// Inspect reads only a bounded, regular, unchanged file in this capture.
func (c *Capture) Inspect(name string, limit int64) Artifact {
	a := Artifact{Name: name, Status: "unavailable"}
	if limit < 1 {
		return a
	}
	path, err := c.path(name)
	if err != nil {
		return a
	}
	before, err := os.Lstat(path)
	if os.IsNotExist(err) {
		a.Status = "absent"
		return a
	}
	if err != nil || !before.Mode().IsRegular() {
		return a
	}
	a.Bytes = before.Size()
	if a.Bytes > limit {
		a.Status = "oversized"
		return a
	}
	file, err := os.Open(path)
	if err != nil {
		return a
	}
	defer func() { _ = file.Close() }()
	opened, err := file.Stat()
	if err != nil || !os.SameFile(before, opened) {
		return a
	}
	h := sha256.New()
	var prefixBytes int64
	if runtime.GOOS == "windows" && strings.HasSuffix(name, ".dmp") {
		var header [4]byte
		if _, err := io.ReadFull(file, header[:]); err != nil || string(header[:]) != "MDMP" {
			return a
		}
		_, _ = h.Write(header[:])
		prefixBytes = 4
	}
	n, err := io.Copy(h, io.LimitReader(file, limit+1-prefixBytes))
	n += prefixBytes
	after, statErr := file.Stat()
	if err != nil || statErr != nil || n != before.Size() || after.Size() != n || !after.ModTime().Equal(before.ModTime()) {
		return a
	}
	a.Status = "collected"
	a.SHA256 = hex.EncodeToString(h.Sum(nil))
	return a
}

// Save exclusively retains one correctly identified incident envelope.
func (c *Capture) Save(record Record) error {
	if record.ID != c.ID || record.Version != 1 || record.Kind != "go-schedule.incident" {
		return errors.New("incident record kind, identity or version mismatch")
	}
	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return err
	}
	return c.Write("incident.json", append(data, '\n'))
}

// DiscardOversized enforces retained-byte limits after the owned producer stops.
// It does not claim to limit the producer's peak temporary disk usage.
func (c *Capture) DiscardOversized(name string, limit int64) error {
	path, err := c.path(name)
	if err != nil {
		return err
	}
	info, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		return errors.New("cannot discard a linked or nonregular incident artifact")
	}
	if info.Size() <= limit {
		return nil
	}
	return os.Remove(path)
}

// FileHash inventories a selected executable without reading its contents into memory.
func FileHash(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func() { _ = file.Close() }()
	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

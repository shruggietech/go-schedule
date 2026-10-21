package clievaluate

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

// ValidateArtifactPath requires one portable canonical relative regular-file path.
func ValidateArtifactPath(name string) error {
	if name == "" || name == "." || path.IsAbs(name) || path.Clean(name) != name || strings.ContainsAny(name, "\\:<>\"|?*\x00") || strings.ContainsFunc(name, func(r rune) bool { return r < 32 || r == 127 }) || strings.HasPrefix(name, "../") || strings.HasSuffix(name, "/") {
		return errors.New("artifact requires a portable canonical relative path")
	}
	for _, part := range strings.Split(name, "/") {
		base := strings.ToUpper(strings.SplitN(part, ".", 2)[0])
		if strings.HasSuffix(part, ".") || strings.HasSuffix(part, " ") || base == "CON" || base == "PRN" || base == "AUX" || base == "NUL" || len(base) == 4 && (strings.HasPrefix(base, "COM") || strings.HasPrefix(base, "LPT")) && base[3] >= '1' && base[3] <= '9' {
			return errors.New("artifact path contains a reserved Windows name")
		}
	}
	return nil
}

// Digest returns the lowercase SHA-256 identity of exact bytes.
func Digest(data []byte) string { sum := sha256.Sum256(data); return hex.EncodeToString(sum[:]) }

// ArtifactStore confines redacted evidence to one private run directory.
type ArtifactStore struct {
	root     *os.Root
	limit    int64
	used     int64
	Redactor *Redactor
}

// OpenArtifacts opens an existing owned artifact root, without changing its parent.
func OpenArtifacts(directory string, limit int64, redactor *Redactor) (*ArtifactStore, error) {
	if limit <= 0 {
		return nil, errors.New("artifact byte budget must be positive")
	}
	root, err := os.OpenRoot(directory)
	if err != nil {
		return nil, err
	}
	used := int64(0)
	entries := 0
	walkErr := fs.WalkDir(root.FS(), ".", func(_ string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		entries++
		if entries > 30000 {
			return errors.New("artifact root contains too many entries")
		}
		if d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return errors.New("artifact root contains a non-regular file")
		}
		used += info.Size()
		if used > limit {
			return errors.New("existing artifact root exceeds budget")
		}
		return nil
	})
	if walkErr != nil {
		return nil, errors.Join(walkErr, root.Close())
	}
	if redactor == nil {
		redactor = NewRedactor()
	}
	return &ArtifactStore{root: root, limit: limit, used: used, Redactor: redactor}, nil
}

// Replace atomically checkpoints sanitized derived data, retaining the prior file on error.
func (s *ArtifactStore) Replace(name string, data []byte) error {
	if err := ValidateArtifactPath(name); err != nil {
		return err
	}
	data = s.Redactor.Sanitize(data)
	old := int64(0)
	info, err := s.root.Stat(name)
	if err == nil {
		if !info.Mode().IsRegular() {
			return errors.New("checkpoint must replace a regular file")
		}
		old = info.Size()
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	if s.used-old+int64(len(data)) > s.limit {
		return errors.New("checkpoint exceeds artifact budget")
	}
	temporary := name + "." + uuid.NewString() + ".tmp"
	file, err := s.root.OpenFile(temporary, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	_, writeErr := file.Write(data)
	syncErr := file.Sync()
	closeErr := file.Close()
	if err := errors.Join(writeErr, syncErr, closeErr); err != nil {
		return errors.Join(err, s.root.Remove(temporary))
	}
	if err := s.root.Rename(temporary, name); err != nil {
		return errors.Join(err, s.root.Remove(temporary))
	}
	s.used += int64(len(data)) - old
	return nil
}

// AppendJournal persists a complete sanitized event without rewriting its prefix.
func (s *ArtifactStore) AppendJournal(name string, data []byte) error {
	if err := ValidateArtifactPath(name); err != nil {
		return err
	}
	data = s.Redactor.Sanitize(data)
	if int64(len(data)) > s.limit-s.used {
		return errors.New("journal exceeds artifact budget")
	}
	if info, err := s.root.Lstat(name); err == nil && !info.Mode().IsRegular() {
		return errors.New("journal must be a regular file")
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	file, err := s.root.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_APPEND|regularReadFlags, 0600)
	if err != nil {
		return err
	}
	info, err := file.Stat()
	if err != nil || !info.Mode().IsRegular() {
		return errors.Join(errors.New("journal must be a regular file"), err, file.Close())
	}
	n, writeErr := file.Write(data)
	s.used += int64(n)
	return errors.Join(writeErr, file.Sync(), file.Close())
}

// ReadDocument reads a bounded contract through the root and validates its raw schema.
func (s *ArtifactStore) ReadDocument(name, kind string, target any) error {
	data, err := s.readRegular(name, documentLimit)
	if err != nil {
		return err
	}
	if err := Decode(bytes.NewReader(data), target); err != nil {
		return err
	}
	return ValidateSchema(kind, data)
}

func (s *ArtifactStore) readRegular(name string, limit int64) ([]byte, error) {
	if err := ValidateArtifactPath(name); err != nil {
		return nil, err
	}
	file, err := s.root.OpenFile(name, os.O_RDONLY|regularReadFlags, 0)
	if err != nil {
		return nil, err
	}
	info, err := file.Stat()
	if err != nil || !info.Mode().IsRegular() {
		return nil, errors.Join(errors.New("QA input must be a regular file"), err, file.Close())
	}
	data, readErr := io.ReadAll(io.LimitReader(file, limit+1))
	closeErr := file.Close()
	if readErr != nil || closeErr != nil {
		return nil, errors.Join(readErr, closeErr)
	}
	if int64(len(data)) > limit {
		return nil, errors.New("QA input exceeds byte limit")
	}
	return data, nil
}

// ReadInput decodes a bounded regular input without following links outside its parent.
func ReadInput(name string, target any) error {
	absolute, err := filepath.Abs(name)
	if err != nil {
		return err
	}
	root, err := os.OpenRoot(filepath.Dir(absolute))
	if err != nil {
		return err
	}
	store := ArtifactStore{root: root}
	data, readErr := store.readRegular(filepath.Base(absolute), documentLimit)
	closeErr := root.Close()
	if err := errors.Join(readErr, closeErr); err != nil {
		return err
	}
	return Decode(bytes.NewReader(data), target)
}

// Close releases the confined root handle.
func (s *ArtifactStore) Close() error { return s.root.Close() }

// Write persists only sanitized bytes with exclusive creation and private permissions.
func (s *ArtifactStore) Write(name, media string, data []byte) (Artifact, error) {
	if err := ValidateArtifactPath(name); err != nil {
		return Artifact{}, err
	}
	data = s.Redactor.Sanitize(data)
	if int64(len(data)) > s.limit-s.used {
		return Artifact{}, errors.New("artifact budget exhausted")
	}
	local, err := filepath.Localize(name)
	if err != nil {
		return Artifact{}, err
	}
	if err := s.root.MkdirAll(filepath.Dir(local), 0700); err != nil {
		return Artifact{}, err
	}
	file, err := s.root.OpenFile(local, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return Artifact{}, err
	}
	_, writeErr := file.Write(data)
	syncErr := file.Sync()
	closeErr := file.Close()
	if err := errors.Join(writeErr, syncErr, closeErr); err != nil {
		return Artifact{}, fmt.Errorf("persist artifact: %w", err)
	}
	s.used += int64(len(data))
	return Artifact{Path: name, MediaType: media, Bytes: int64(len(data)), SHA256: Digest(data), Redacted: true, RedactionPolicy: "fixture-secrets-v1"}, nil
}

// Read verifies a regular file through the same confined handle used for hashing.
func (s *ArtifactStore) Read(a Artifact) ([]byte, error) {
	if err := ValidateArtifactPath(a.Path); err != nil {
		return nil, err
	}
	if a.Bytes < 0 || a.Bytes > s.limit || !a.Redacted || a.RedactionPolicy != "fixture-secrets-v1" {
		return nil, errors.New("invalid artifact bounds or redaction metadata")
	}
	data, err := s.readRegular(a.Path, a.Bytes)
	if err != nil {
		return nil, err
	}
	if int64(len(data)) != a.Bytes || Digest(data) != a.SHA256 {
		return nil, fmt.Errorf("artifact hash/size mismatch: %s", a.Path)
	}
	return data, nil
}

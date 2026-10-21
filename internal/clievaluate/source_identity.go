package clievaluate

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
)

// digestSourceFiles binds exact working-tree bytes, including nonignored untracked files.
// Only the digest is persisted; source contents never enter captured diagnostics.
func digestSourceFiles(ctx context.Context, directory string, paths []string) (string, error) {
	root, err := os.OpenRoot(directory)
	if err != nil {
		return "", err
	}
	defer func() { _ = root.Close() }()
	paths = append([]string{}, paths...)
	slices.Sort(paths)
	paths = slices.Compact(paths)
	hash := sha256.New()
	for _, path := range paths {
		if err := ctx.Err(); err != nil {
			return "", err
		}
		if path == "" {
			return "", errors.New("empty evaluated source path")
		}
		info, err := root.Lstat(path)
		if errors.Is(err, os.ErrNotExist) {
			_, _ = fmt.Fprintf(hash, "%d:%s:missing\x00", len(path), path)
			continue
		}
		if err != nil {
			return "", err
		}
		if !info.Mode().IsRegular() {
			return "", errors.New("evaluated source identity requires regular tracked and untracked files")
		}
		file, err := root.OpenFile(path, os.O_RDONLY|regularReadFlags, 0)
		if err != nil {
			return "", err
		}
		_, _ = fmt.Fprintf(hash, "%d:%s:%d:%d\x00", len(path), path, info.Mode(), info.Size())
		count, readErr := io.Copy(hash, sourceReader{ctx: ctx, reader: file})
		closeErr := file.Close()
		if readErr != nil || closeErr != nil || count != info.Size() {
			return "", errors.Join(errors.New("evaluated source changed or could not be read while hashing"), readErr, closeErr)
		}
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

type sourceReader struct {
	ctx    context.Context
	reader io.Reader
}

func (r sourceReader) Read(p []byte) (int, error) {
	if err := r.ctx.Err(); err != nil {
		return 0, err
	}
	return r.reader.Read(p)
}

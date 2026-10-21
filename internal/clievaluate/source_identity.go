package clievaluate

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/shruggietech/go-schedule/internal/clifixture"
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

func identifySource(ctx context.Context, source string) (Subject, error) {
	gitPath, err := exec.LookPath("git")
	if err != nil {
		return Subject{}, err
	}
	s := Subject{Repository: "https://github.com/shruggietech/go-schedule", Version: "development"}
	for _, query := range []struct {
		args   []string
		target *string
	}{{[]string{"rev-parse", "HEAD"}, &s.Commit}, {[]string{"remote", "get-url", "origin"}, &s.Repository}} {
		r := captureCommand(ctx, clifixture.Command{Executable: gitPath, Arguments: query.args, Directory: source, Environment: quietEnvironment()}, nil, 4096)
		if r.Code != 0 || r.Failure != nil {
			return s, fmt.Errorf("cannot identify evaluated Git source: exit=%d failure=%v stderr=%s", r.Code, r.Failure, r.Err)
		}
		*query.target = strings.TrimSpace(string(r.Out))
	}
	status := captureCommand(ctx, clifixture.Command{Executable: gitPath, Arguments: []string{"status", "--porcelain"}, Directory: source, Environment: quietEnvironment()}, nil, 65536)
	if status.Code != 0 || status.Exceeded {
		return s, errors.New("cannot establish source dirtiness")
	}
	s.Dirty = len(status.Out) > 0
	files := captureCommand(ctx, clifixture.Command{Executable: gitPath, Arguments: []string{"ls-files", "--cached", "--others", "--exclude-standard", "-z"}, Directory: source, Environment: quietEnvironment()}, nil, 1<<20)
	if files.Code != 0 || files.Failure != nil || files.Exceeded {
		return s, errors.New("cannot enumerate evaluated source files")
	}
	s.SourceSHA256, err = digestSourceFiles(ctx, source, strings.Split(strings.TrimSuffix(string(files.Out), "\x00"), "\x00"))
	if err != nil {
		return s, err
	}
	s.Repository, err = canonicalRemote(s.Repository)
	if err != nil {
		return s, err
	}
	return s, nil
}

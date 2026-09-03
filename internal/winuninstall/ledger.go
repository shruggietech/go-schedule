package winuninstall

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func writeLedger(path string, result Result) (resultErr error) {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}
	data = append(data, '\n')
	directory, err := os.OpenRoot(filepath.Dir(path))
	if err != nil {
		return fmt.Errorf("open ledger directory: %w", err)
	}
	defer func() {
		if err := directory.Close(); err != nil {
			resultErr = errors.Join(resultErr, fmt.Errorf("close ledger directory: %w", err))
		}
	}()
	name := filepath.Base(path)
	temporary := name + ".tmp"
	file, err := directory.OpenFile(temporary, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("create temporary ledger: %w", err)
	}
	written := false
	closed := false
	defer func() {
		if !closed {
			if err := file.Close(); err != nil {
				resultErr = errors.Join(resultErr, fmt.Errorf("close temporary ledger: %w", err))
			}
		}
		if !written {
			if err := directory.Remove(temporary); err != nil && !errors.Is(err, os.ErrNotExist) {
				resultErr = errors.Join(resultErr, fmt.Errorf("remove temporary ledger: %w", err))
			}
		}
	}()
	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("write temporary ledger: %w", err)
	}
	if err := file.Sync(); err != nil {
		return fmt.Errorf("flush temporary ledger: %w", err)
	}
	if err := file.Close(); err != nil {
		closed = true
		return fmt.Errorf("close temporary ledger: %w", err)
	}
	closed = true
	if err := directory.Rename(temporary, name); err != nil {
		return fmt.Errorf("replace ledger: %w", err)
	}
	written = true
	return nil
}

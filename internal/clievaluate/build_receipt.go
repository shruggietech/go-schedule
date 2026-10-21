package clievaluate

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

func sameSource(a, b Subject) bool {
	return a.Repository == b.Repository && a.Commit == b.Commit && a.Dirty == b.Dirty && a.SourceSHA256 == b.SourceSHA256
}

func preparedName(role, platform string) string {
	name := map[string]string{"cli": "gosched", "daemon": "goschedd", "runner": "gosched-qa"}[role]
	if name == "" || platform != "windows" && platform != "linux" {
		return ""
	}
	if platform == "windows" {
		name += ".exe"
	}
	return platform + "/" + name
}

func preparedBinaryDigests(directory string) ([]Binary, error) {
	binaries := []Binary{}
	for _, platform := range []string{"windows", "linux"} {
		for _, role := range []string{"cli", "daemon", "runner"} {
			name := preparedName(role, platform)
			data, err := readBinary(filepath.Join(directory, filepath.FromSlash(name)))
			if err != nil {
				return nil, err
			}
			binaries = append(binaries, Binary{role, platform, "amd64", name, Digest(data)})
		}
	}
	return binaries, nil
}

func validateBuildReceipt(receipt BuildReceipt, subject Subject) error {
	if !sameSource(Subject{Repository: receipt.Repository, Commit: receipt.Commit, Dirty: receipt.Dirty, SourceSHA256: receipt.SourceSHA256}, subject) || receipt.SourceSHA256 == "" {
		return errors.New("prepared executable source differs from the selected checkout; rebuild the inputs from the unchanged source")
	}
	if len(receipt.Binaries) != 6 || len(subject.Binaries) != 6 {
		return errors.New("build receipt must bind all six evaluation executables")
	}
	seen := map[string]bool{}
	for _, binary := range receipt.Binaries {
		name := preparedName(binary.Role, binary.OS)
		if name == "" || seen[name] || binary.Path != name || binary.Arch != "amd64" {
			return errors.New("build receipt has duplicate or unexpected executable identities")
		}
		seen[name] = true
		matched := false
		for _, actual := range subject.Binaries {
			if actual.Role == binary.Role && actual.OS == binary.OS && actual.Arch == binary.Arch && actual.SHA256 == binary.SHA256 {
				matched = true
			}
		}
		if !matched {
			return errors.New("prepared executable bytes differ from the source-bound build receipt")
		}
	}
	return nil
}

func requireBuildReceipt(directory string, subject Subject) (resultErr error) {
	root, err := os.OpenRoot(directory)
	if err != nil {
		return err
	}
	defer func() { resultErr = errors.Join(resultErr, root.Close()) }()
	reader := ArtifactStore{root: root}
	var receipt BuildReceipt
	if err := reader.ReadDocument("build-receipt.json", "build-receipt", &receipt); err != nil {
		return fmt.Errorf("prepared build receipt is missing or invalid; run build in a new input directory: %w", err)
	}
	return validateBuildReceipt(receipt, subject)
}

func writeBuildReceipt(directory string, receipt BuildReceipt) (resultErr error) {
	data, err := Encode(receipt)
	if err != nil {
		return err
	}
	if err := ValidateSchema("build-receipt", data); err != nil {
		return err
	}
	root, err := os.OpenRoot(directory)
	if err != nil {
		return err
	}
	defer func() { resultErr = errors.Join(resultErr, root.Close()) }()
	name := "build-receipt-" + uuid.NewString() + ".tmp"
	file, err := root.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	_, writeErr := file.Write(data)
	syncErr := file.Sync()
	closeErr := file.Close()
	if err := errors.Join(writeErr, syncErr, closeErr); err != nil {
		return errors.Join(err, root.Remove(name))
	}
	if err := root.Rename(name, "build-receipt.json"); err != nil {
		return errors.Join(err, root.Remove(name))
	}
	return nil
}

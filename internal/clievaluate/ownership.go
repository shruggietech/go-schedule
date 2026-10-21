package clievaluate

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/shruggietech/go-schedule/internal/clifixture"
)

func retainOwnership(store *ArtifactStore, r *Report, prefix string, w *clifixture.Workspace) (resultErr error) {
	if err := w.Contains(filepath.Join(w.Evidence, "ownership.json")); err != nil {
		return err
	}
	root, err := os.OpenRoot(w.Evidence)
	if err != nil {
		return err
	}
	defer func() { resultErr = errors.Join(resultErr, root.Close()) }()
	reader := ArtifactStore{root: root}
	data, err := reader.readRegular("ownership.json", 16384)
	if err != nil {
		return err
	}
	name := prefix + "/ownership.json"
	for _, a := range r.Artifacts {
		if a.Path == name {
			if Digest(data) != a.SHA256 {
				return errors.New("owned workspace marker changed during execution")
			}
			_, err := store.Read(a)
			return err
		}
	}
	a, err := store.Write(name, "application/json", data)
	if err != nil {
		return err
	}
	r.Artifacts = append(r.Artifacts, a)
	return nil
}

package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"

	"github.com/shruggietech/go-schedule/v2/internal/bundle"
)

type watcherBindingError struct{ message string }

func (e *watcherBindingError) Error() string { return e.message }

// watcherPathBindings preserves binding identity until duplicate detection is complete.
type watcherPathBindings struct {
	paths   map[string]string
	decoded bool
}

func (paths *watcherPathBindings) UnmarshalJSON(data []byte) error {
	if paths.decoded {
		return &watcherBindingError{message: "duplicate watcher_paths field; supply one watcher_paths object with each watcher portable_id exactly once"}
	}
	paths.decoded = true
	decoder := json.NewDecoder(bytes.NewReader(data))
	token, err := decoder.Token()
	if err != nil {
		return err
	}
	if token == nil {
		paths.paths = nil
		return nil
	}
	if token != json.Delim('{') {
		return errors.New("watcher_paths must be an object")
	}
	result := make(map[string]string)
	for decoder.More() {
		token, err := decoder.Token()
		if err != nil {
			return err
		}
		id, ok := token.(string)
		if !ok {
			return errors.New("watcher path identity must be a string")
		}
		if _, exists := result[id]; exists {
			return &watcherBindingError{message: fmt.Sprintf("duplicate watcher path binding %q; supply each watcher portable_id exactly once", id)}
		}
		var path string
		if err := decoder.Decode(&path); err != nil {
			return err
		}
		result[id] = path
	}
	if _, err := decoder.Token(); err != nil {
		return err
	}
	if _, err := decoder.Token(); err != io.EOF {
		return errors.New("watcher_paths contains trailing JSON")
	}
	paths.paths = result
	return nil
}

func validateWatcherPathBindings(document bundle.Document, paths map[string]string) error {
	watchers := make(map[string]bool, len(document.Watchers))
	for _, watcher := range document.Watchers {
		watchers[watcher.PortableID] = true
	}
	identities := make([]string, 0, len(paths))
	for id := range paths {
		identities = append(identities, id)
	}
	sort.Strings(identities)
	for _, id := range identities {
		if !watchers[id] {
			return &watcherBindingError{message: fmt.Sprintf("watcher path binding %q is not a watcher in this bundle; use a watcher portable_id from the bundle or remove the binding", id)}
		}
	}
	return nil
}

func writeBundleRequestError(w http.ResponseWriter, err error) {
	var bindingError *watcherBindingError
	if errors.As(err, &bindingError) {
		writeError(w, http.StatusBadRequest, CodeValidation, "watcher_paths", bindingError.Error())
		return
	}
	writeError(w, http.StatusBadRequest, CodeValidation, "body", "invalid JSON")
}

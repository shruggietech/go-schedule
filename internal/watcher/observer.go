// Package watcher converts native filesystem notifications into stable task dispatch requests.
package watcher

import "github.com/fsnotify/fsnotify"

type observer interface {
	Add(string) error
	Close() error
	Events() <-chan fsnotify.Event
	Errors() <-chan error
}

type fsObserver struct {
	watcher *fsnotify.Watcher
}

func newFSObserver() (observer, error) {
	w, err := fsnotify.NewBufferedWatcher(256)
	if err != nil {
		return nil, err
	}
	return &fsObserver{watcher: w}, nil
}

func (o *fsObserver) Add(path string) error         { return o.watcher.Add(path) }
func (o *fsObserver) Close() error                  { return o.watcher.Close() }
func (o *fsObserver) Events() <-chan fsnotify.Event { return o.watcher.Events }
func (o *fsObserver) Errors() <-chan error          { return o.watcher.Errors }

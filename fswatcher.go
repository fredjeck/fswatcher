package fswatcher

import (
	"errors"
	"os"
	"path/filepath"
	"time"
)

// A file extension that matches all files.
const ExtAllFiles = "*"

// FsEventHandler represents a callback function triggered when a change is detected in the monitored file system.
// If a FsEventHandler returns true, the FsWatcher to which the handler is bound will skip the remaining files in the current directory.
type FsEventHandler func(path string, info os.FileInfo) bool

// FsWatcher watches for changes in the specified directory and triggers the registered FsEventHandlers.
// FsWatcher does not rely on system specific features to detect changes but polls periodically files modification times.
type FsWatcher struct {
	handlers    map[string]FsEventHandler
	ignored     map[string]string
	root        string
	interval    int
	started     bool
	stopPending bool
}

// NewFsWatcher creates a new FsWatcher ready to monotir the directory pointed by root
// If root does not exists or is a file, NewFsWatcher will raise an error and returns a nil FsWatcher
func NewFsWatcher(root string, interval int) (*FsWatcher, error) {

	p := filepath.Clean(root)

	info, err := os.Stat(p)
	if err != nil || !info.IsDir() {
		return nil, errors.New("The specified path either doesn't exists or leads to a file " + p)
	}

	i := 500
	if interval > 0 {
		i = interval
	}
	return &FsWatcher{
		make(map[string]FsEventHandler),
		make(map[string]string),
		p,
		i,
		false,
		false,
	}, nil
}

// RegisterFileExtension registers a FsEventHandler function that will be triggered when a file with the registered extension has been modified
func (f *FsWatcher) RegisterFileExtension(extension string, handler FsEventHandler) {
	f.handlers[extension] = handler
}

// IsStarted returns true if the watcher is currently watching for changes in the file system
func (f *FsWatcher) IsStarted() bool {
	return f.started
}

// Stop requests the current watcher to stop monitoring for file system changes
// Note that stop does not stop the watcher immediately but does a gracefull shutdown
func (f *FsWatcher) Stop() {
	f.stopPending = false
}

// Skip registers the given directory so that it won't be monitored
func (f *FsWatcher) Skip(directory string) {
	if len(directory) > 0 {
		f.ignored[directory] = directory
	}
}

// Watch watches for changes in the directory pointed by root at a given interval.
// If a change is detected (based on the file modification timestamp), the provided change handler is triggered.
// The watching process is launched in a new goroutine
// TODO : Improve stoping process (currently not safe)
func (f *FsWatcher) Watch() {
	if !f.started {
		f.started = true
		go func() {
			// Reference time used to dectect file changes
			start := time.Now()
			for {
				filepath.Walk(f.root, func(path string, info os.FileInfo, err error) error {
					if f.stopPending {
						return errors.New("Stopping")
					}
					b := filepath.Base(path)
					dir, ok := f.ignored[b]
					if ok && b == dir {
						return filepath.SkipDir
					}

					// TODO : This is not optimized.
					// Using the * handler, a directory will only be skipped if no files
					if info.ModTime().After(start) {
						handler, ok := f.handlers[filepath.Ext(path)]
						if ok == true {
							if handler(path, info) == true {
								return filepath.SkipDir
							}
						}
					}
					return nil
				})

				if f.stopPending {
					f.started = false
					return
				}
				start = time.Now()
				time.Sleep(time.Duration(f.interval) * time.Millisecond)
			}
		}()
	}
}

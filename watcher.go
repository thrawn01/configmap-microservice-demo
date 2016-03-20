package main

import (
	"time"

	"gopkg.in/fsnotify.v1"
)

/*
 Watches a file on a set interval, and preforms de-duplication of write
 events such that only 1 write event is reported even if multiple writes
 happened during the specified duration.
*/
type FileWatcher struct {
	fsNotify *fsnotify.Watcher
	interval time.Duration
	done     chan struct{}
	callback func()
}

/*
 Begin watching a file with a specific interval and action
*/
func WatchFile(path string, interval time.Duration, action func()) (*FileWatcher, error) {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	// Add the file to be watched
	fsWatcher.Add(path)

	watcher := &FileWatcher{
		fsWatcher,
		interval,
		make(chan struct{}, 1),
		action,
	}
	// Launch a go thread to watch the file
	go watcher.run()

	return watcher, err
}

func (self *FileWatcher) run() {
	// Check for write events at this interval
	tick := time.Tick(self.interval)

	var lastWriteEvent *fsnotify.Event
	for {
		select {
		case event := <-self.fsNotify.Events:
			// If it was a write event
			if event.Op == fsnotify.Write {
				lastWriteEvent = &event
			}
		case <-tick:
			// No events during this interval
			if lastWriteEvent == nil {
				continue
			}
			// Execute the callback
			self.callback()
			// Reset the last event
			lastWriteEvent = nil
		case <-self.done:
			goto Close
		}
	}
Close:
	close(self.done)
}

func (self *FileWatcher) Close() {
	self.done <- struct{}{}
	self.fsNotify.Close()
}

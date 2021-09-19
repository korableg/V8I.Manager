package watcher

import (
	"context"

	"github.com/fsnotify/fsnotify"
)

func Watch(ctx context.Context, filenames ...string) (chan string, chan error) {

	cfn := make(chan string, 1)
	errChan := make(chan error, 1)

	go watch(ctx, cfn, errChan, filenames...)

	return cfn, errChan

}

func watch(ctx context.Context, cfn chan<- string, errChan chan<- error, filenames ...string) {

	fswatcher, err := fsnotify.NewWatcher()
	if err != nil {
		errChan <- err
		return
	}

	defer fswatcher.Close()

	for _, name := range filenames {
		if len(name) == 0 {
			continue
		}
		fswatcher.Add(name)
	}

	for {
		select {
		case <-ctx.Done():
			return
		case event, ok := <-fswatcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				cfn <- event.Name
			}
		case err, ok := <-fswatcher.Errors:
			if !ok {
				return
			}
			errChan <- err
		}
	}

}

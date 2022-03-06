package watcher

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"time"
)

type (
	Watcher interface {
		Start() <-chan string
	}

	Fabric func(path string, stopCh <-chan struct{}) (Watcher, error)

	watcher struct {
		path   string
		hash   []byte
		stopCh <-chan struct{}
	}
)

func NewWatcher(path string, stopCh <-chan struct{}) (Watcher, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("get file stat by path %s: %w", path, err)
	}

	if stopCh == nil {
		return nil, errors.New("stop channel is not initialized")
	}

	w := &watcher{
		path:   path,
		stopCh: stopCh,
	}
	return w, nil
}

func (w *watcher) Start() <-chan string {
	hashChan := make(chan string, 1)

	ticker := time.NewTicker(5 * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				hash, err := getFileHash(w.path)
				if err != nil {
					logrus.Errorf("get file hash: %s", err.Error())
					continue
				}

				if !bytes.Equal(hash, w.hash) {
					w.hash = hash
					hashChan <- string(hash)
				}
			case <-w.stopCh:
				ticker.Stop()
				close(hashChan)
				return
			}
		}
	}()

	return hashChan
}

func getFileHash(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer func() {
		if err = f.Close(); err != nil {
			logrus.Errorf("close file by path %s error %s", path, err.Error())
		}
	}()

	s := sha256.New()

	if _, err = io.Copy(s, f); err != nil {
		return nil, fmt.Errorf("copy file to sha256: %w", err)
	}

	return s.Sum(nil), nil
}

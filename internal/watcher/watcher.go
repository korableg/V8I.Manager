package watcher

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"time"
)

func Watch(ctx context.Context, filenames ...string) (chan string, chan error) {

	cfn := make(chan string, 1)
	errChan := make(chan error, 1)

	go watch(ctx, cfn, errChan, filenames...)

	return cfn, errChan

}

func watch(ctx context.Context, cfn chan<- string, errChan chan<- error, filenames ...string) {

	checksums := make(map[string]string, len(filenames))
	ticker := time.NewTicker(5 * time.Second)

	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			close(cfn)
			return
		case <-ticker.C:
			for _, name := range filenames {
				cs, err := getFileCheckSum(name)
				if err != nil {
					errChan <- err
					continue
				}
				if cs == "" {
					continue
				}
				cachecs := checksums[name]
				if cachecs != cs {
					checksums[name] = cs
					cfn <- name
				}
			}
		}
	}

}

func getFileCheckSum(filename string) (string, error) {

	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()
	s := sha256.New()

	if _, err := io.Copy(s, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", s.Sum(nil)), nil

}

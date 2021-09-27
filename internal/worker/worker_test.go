package worker

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWorker(t *testing.T) {

	const countOfWrite = 10

	lsts := make([]string, 1)
	lsts[0] = "../../test/test.lst"

	outChan := make(chan []byte, 1)
	errChan := make(chan error, 1)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	w := NewWorker(lsts)
	go func(t *testing.T) {
		err := w.StartWatchingContext(ctx, outChan)
		if err != nil {
			errChan <- err
		}
	}(t)

	time.Sleep(time.Second)

	absLstPath, err := filepath.Abs(lsts[0])
	assert.Equal(t, err, nil)

	f, err := os.OpenFile(absLstPath, os.O_APPEND|os.O_WRONLY, 0600)
	assert.Equal(t, err, nil)

	fileinfo, err := f.Stat()
	assert.Equal(t, err, nil)

	filesize := fileinfo.Size()

	go func() {

		for i := 0; i < countOfWrite; i++ {

			_, err = f.WriteString(string(rune(0)))
			if err != nil {
				errChan <- err
				return
			}

			err = f.Sync()
			if err != nil {
				errChan <- err
				return
			}

		}

		err = f.Close()
		assert.Equal(t, err, nil)

	}()

	for i := 0; i < countOfWrite; i++ {
		select {
		case _, ok := <-outChan:
			{
				assert.Equal(t, ok, true)
			}

		case err, ok := <-errChan:
			{
				assert.Equal(t, ok, true)
				if !assert.Equal(t, err, nil) {
					t.FailNow()
				}
			}
		case <-ctx.Done():
			t.Fatal("interrupted by context")
		}

	}

	err = os.Truncate(absLstPath, filesize)

	assert.Equal(t, err, nil)

}

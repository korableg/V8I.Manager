package worker

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/korableg/V8I.Manager/pkg/v8i/v8iwriter"
	"github.com/korableg/V8I.Manager/pkg/v8i/v8iwriter/v8imockwriter"
	"github.com/stretchr/testify/assert"
)

func TestWorker(t *testing.T) {

	const countOfWrite = 10

	sourceLst := filepath.FromSlash("../../test/test.lst")
	destinationLst := filepath.FromSlash("../../test/test1.lst")

	sourceBytes, err := os.ReadFile(sourceLst)
	assert.Equal(t, err, nil)

	err = os.WriteFile(destinationLst, sourceBytes, 0644)
	assert.Equal(t, err, nil)

	lsts := make([]string, 1)
	lsts[0] = destinationLst

	v8iW := v8imockwriter.New()

	v8is := make([]v8iwriter.V8IWriter, 1)
	v8is[0] = v8iW

	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()

	errChan := make(chan error, 1)

	w := NewWorker(lsts, v8is)
	go func(t *testing.T) {
		err := w.StartWatchingContext(ctx)
		if err != nil {
			errChan <- err
		}
	}(t)

	time.Sleep(time.Second)

	go func() {

		f, err := os.OpenFile(lsts[0], os.O_APPEND|os.O_WRONLY, 0644)
		assert.Equal(t, err, nil)

		for i := 0; i < countOfWrite; i++ {

			time.Sleep(time.Millisecond * 500)

			_, err = f.WriteString(string(rune(0)))
			assert.Equal(t, err, nil)

			err = f.Sync()
			assert.Equal(t, err, nil)

		}

		err = f.Close()
		assert.Equal(t, err, nil)

		err = os.Remove(destinationLst)
		assert.Equal(t, err, nil)

		close(v8iW.W)

	}()

	writes := 0

Loop:
	for {
		select {
		case _, ok := <-v8iW.W:
			if !ok {
				break Loop
			}
			writes++
		case <-ctx.Done():
			break Loop
		case err := <-errChan:
			assert.Equal(t, err, nil)
		}
	}

	assert.Equal(t, countOfWrite, writes)

}

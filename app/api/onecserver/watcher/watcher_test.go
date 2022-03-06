package watcher

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestWatcher(t *testing.T) {
	const (
		pathSource      = "../../../test/test.lst"
		pathDestination = "../../../test/test_destination.lst"
		actualFires     = 2
	)

	dataSource, err := ioutil.ReadFile(pathSource)
	require.Nil(t, err)

	err = ioutil.WriteFile(pathDestination, dataSource, 0644)
	require.Nil(t, err)

	defer func() {
		err = os.Remove(pathDestination)
		assert.Nil(t, err)
	}()

	stopCh := make(chan struct{}, 1)
	w, err := NewWatcher(pathDestination, stopCh)
	require.Nil(t, err)

	go func() {
		ctx, _ := context.WithTimeout(context.Background(), 11*time.Second)
		<-ctx.Done()
		close(stopCh)
	}()

	go func() {
		time.AfterFunc(6*time.Second,
			func() {
				_ = ioutil.WriteFile(pathDestination, []byte("test"), 0644)
			})
	}()

	fires := 0
	for range w.Start() {
		fires++
	}

	assert.Equal(t, fires, actualFires)

}

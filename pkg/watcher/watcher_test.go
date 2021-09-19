package watcher

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const pathToLst = "../../test/test.lst"
const pathToNewlst = "../../test/new_test.lst"

func TestWatch(t *testing.T) {

	f, err := os.Create(pathToNewlst)
	if !os.IsExist(err) {
		assert.Equal(t, err, nil)
	}

	err = f.Sync()
	assert.Equal(t, err, nil)

	err = f.Close()
	assert.Equal(t, err, nil)

	time.Sleep(time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	cfn, errChan := Watch(ctx, pathToNewlst)

	go func(t *testing.T, errChan <-chan error) {
		err := <-errChan
		assert.Equal(t, err, nil)
	}(t, errChan)

	go func(ctx context.Context, cfn chan string, errChan chan error) {
		<-ctx.Done()
		close(cfn)
		close(errChan)
	}(ctx, cfn, errChan)

	readBytes, err := os.ReadFile(pathToLst)
	assert.Equal(t, err, nil)

	err = os.WriteFile(pathToNewlst, readBytes, 0644)
	assert.Equal(t, err, nil)

	filename := <-cfn
	assert.Equal(t, pathToNewlst, filename)

	cancel()

	os.Remove(pathToNewlst)

}

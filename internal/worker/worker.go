package worker

import (
	"bytes"
	"context"
	"errors"
	"os"
	"time"

	"github.com/korableg/V8I.Manager/pkg/clusterdb"
	"github.com/korableg/V8I.Manager/pkg/lstparser"
	"github.com/korableg/V8I.Manager/pkg/v8ibuilder"
	"github.com/korableg/V8I.Manager/pkg/watcher"
)

type Worker struct {
	lsts []string
}

func NewWorker(lsts []string) *Worker {

	w := &Worker{
		lsts: lsts,
	}

	return w

}

func (w *Worker) StartWatchingContext(ctx context.Context, out chan<- []byte) error {

	if len(w.lsts) == 0 {
		return errors.New("list of lst files is empty")
	}

	v8iBytes, err := w.buildv8i()
	if err != nil {
		return err
	}

	out <- v8iBytes

	changedChan, errChan := watcher.Watch(ctx, w.lsts...)

	for {
		select {
		case err := <-errChan:
			return err
		case <-ctx.Done():
			return nil
		case _, ok := <-changedChan:

			if !ok {
				return nil
			}

			time.Sleep(time.Millisecond * 300)

			v8iBytes, err := w.buildv8i()
			if err != nil {
				return err
			}

			out <- v8iBytes
		}

	}

}

func (w *Worker) buildv8i() ([]byte, error) {

	clusterDBs, err := lstToClusterDbs(w.lsts)
	if err != nil {
		return nil, err
	}

	v8iBytes, err := v8ibuilder.Build(clusterDBs...)
	if err != nil {
		return nil, err
	}

	return v8iBytes, nil

}

func LstToV8i(lstFileNames []string, v8iFileNames []string) error {

	clusterDBs, err := lstToClusterDbs(lstFileNames)
	if err != nil {
		return err
	}

	v8iData, err := v8ibuilder.Build(clusterDBs...)
	if err != nil {
		return err
	}

	err = V8IBytesToFiles(v8iData, v8iFileNames)
	if err != nil {
		return err
	}

	return nil

}

func V8IBytesToFiles(v8iBytes []byte, v8iFileNames []string) error {
	for _, v8iFilename := range v8iFileNames {
		err := os.WriteFile(v8iFilename, v8iBytes, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func lstToClusterDbs(lstFilenames []string) ([]*clusterdb.ClusterDB, error) {

	lstData := bytes.NewBuffer(nil)

	for _, lstFilename := range lstFilenames {
		d, err := os.ReadFile(lstFilename)
		if err != nil {
			return nil, err
		}
		_, err = lstData.Write(d)
		if err != nil {
			return nil, err
		}
	}

	clusterDBs, err := lstparser.Parse(lstData.Bytes())
	if err != nil {
		return nil, err
	}

	return clusterDBs, nil

}

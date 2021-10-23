package worker

import (
	"bytes"
	"context"
	"errors"
	"os"
	"time"

	"github.com/korableg/V8I.Manager/internal/watcher"
	"github.com/korableg/V8I.Manager/pkg/clusterdb"
	"github.com/korableg/V8I.Manager/pkg/lstparser"
	"github.com/korableg/V8I.Manager/pkg/v8i/v8ibuilder"
	"github.com/korableg/V8I.Manager/pkg/v8i/v8iwriter"
	log "github.com/sirupsen/logrus"
)

type Worker struct {
	lsts []string
	v8is []v8iwriter.V8IWriter
}

func NewWorker(lsts []string, v8is []v8iwriter.V8IWriter) *Worker {

	w := &Worker{
		lsts: lsts,
		v8is: v8is,
	}

	return w

}

func (w *Worker) StartWatchingContext(ctx context.Context) error {

	if len(w.lsts) == 0 {
		return errors.New("list of lst files is empty")
	}

	v8iBytes, err := w.buildv8i()
	if err != nil {
		return err
	}

	err = w.writev8i(v8iBytes)
	if err != nil {
		return err
	}

	changedChan, errChan := watcher.Watch(ctx, w.lsts...)

	log.Info("start watching lst files")

	for {
		select {
		case err := <-errChan:
			return err
		case <-ctx.Done():
			return nil
		case p, ok := <-changedChan:

			if !ok {
				return nil
			}

			log.Infof("the lst file by path %s was written", p)

			time.Sleep(time.Millisecond * 300)

			v8iBytes, err := w.buildv8i()
			if err != nil {
				return err
			}
			err = w.writev8i(v8iBytes)
			if err != nil {
				return err
			}

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

func (w *Worker) writev8i(v8iBytes []byte) error {

	for _, v := range w.v8is {
		_, err := v.Write(v8iBytes)
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

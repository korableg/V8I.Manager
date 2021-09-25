//go:build windows

package main

import (
	"context"

	"github.com/korableg/V8I.Manager/internal/globals"
	"github.com/korableg/V8I.Manager/internal/worker"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"golang.org/x/sys/windows/svc"
)

type service struct{}

func (s *service) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown

	changes <- svc.Status{State: svc.StartPending}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	v8iChan := make(chan []byte, 1)
	errChan := make(chan error, 1)

	lst := viper.GetStringSlice(_lstFlag)
	v8i := viper.GetStringSlice(_v8iFlag)

	workerInstance := worker.NewWorker(lst)

	go func() {
		log.Infof("starting watch %s files", lst)
		err := workerInstance.StartWatchingContext(ctx, v8iChan)
		if err != nil {
			errChan <- err
		}
	}()

	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

loop:
	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				cancel()
				break loop
			default:
				log.Errorf("unexpected control request #%d", c)
			}
		case <-ctx.Done():
			break loop
		case err, ok := <-errChan:
			{
				cancel()
				if !ok || err == nil {
					log.Error("error channel has been closed or a false positive has occured")
					return
				}
				if err != nil {
					log.Error(err)
					return false, 1
				}
			}
		case v8iBytes, ok := <-v8iChan:
			if !ok {
				cancel()
				log.Error("v8i channel has been close")
				return
			}
			worker.V8IBytesToFiles(v8iBytes, v8i)
			log.Infof("saved v8i files into %s", v8i)
		}
	}
	changes <- svc.Status{State: svc.StopPending}
	return
}

func runService() error {

	log.Info("starting service")

	err := svc.Run(globals.AppName, &service{})
	if err != nil {
		return errors.Wrap(err, "during start service")
	}

	log.Info("service stopped")

	return nil

}

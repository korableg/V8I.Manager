//go:build windows

package windowsservice

import (
	"context"

	"github.com/korableg/V8I.Manager/internal/config"
	"github.com/korableg/V8I.Manager/internal/worker"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/windows/svc"
)

type service struct {
	cfg config.Config
}

func New(cfg config.Config) *service {

	s := &service{
		cfg: cfg,
	}

	return s

}

func (s *service) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown

	changes <- svc.Status{State: svc.StartPending}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errChan := make(chan error, 1)

	workerInstance := worker.NewWorker(s.cfg.Lsts(), s.cfg.V8is())

	go func() {
		log.Infof("starting watch %s files", s.cfg.Lsts())
		err := workerInstance.StartWatchingContext(ctx)
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
				if !ok {
					log.Error("error channel has been closed or a false positive has occured")
				}
				if err != nil {
					log.Error(err)
				}
				return false, 1
			}
		}
	}
	changes <- svc.Status{State: svc.StopPending}
	return
}

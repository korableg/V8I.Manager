//go:build windows

package main

import (
	"fmt"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"golang.org/x/sys/windows/svc"
)

type service struct{}

func (s *service) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown

	changes <- svc.Status{State: svc.StartPending}

	//changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

loop:
	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				break loop
			default:
				log.Error(fmt.Sprintf("unexpected control request #%d", c))
			}
		}
	}
	changes <- svc.Status{State: svc.StopPending}
	return
}

func runService() {

	log.Info("starting service")

	err := svc.Run(svcName, &service{})
	if err != nil {
		log.Fatal(errors.Wrap(err, "during start service"))
	}

	log.Info("service stopped")

}

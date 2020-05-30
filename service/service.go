// +build windows

package main

import (
	"OneCIBasesCreator/OneCIBasesCreator"
	iArgs "OneCIBasesCreator/args"
	"fmt"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"
	"os"
	"sync"
	"time"
)

type service struct{}

var elog debug.Log
var mutex sync.Mutex
var err error
var d = time.Minute
var ticker *time.Ticker

var pathLst, pathiBases []string

var lastModifiedLst map[string]time.Time

func (s *service) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue

	pathLst, err = iArgs.PathLst()
	handleError(err)

	pathiBases, err = iArgs.PathIBases()
	handleError(err)

	changes <- svc.Status{State: svc.StartPending}

	ticker = time.NewTicker(d)
	tick := ticker.C

	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

	mutex.Lock()
	go createiBases()

loop:
	for {
		select {
		case <-tick:
			mutex.Lock()
			go createiBases()
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				break loop
			case svc.Pause:
				changes <- svc.Status{State: svc.Paused, Accepts: cmdsAccepted}
				ticker.Stop()
			case svc.Continue:
				changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
				ticker = time.NewTicker(d)
				tick = ticker.C
			default:
				elog.Error(1, fmt.Sprintf("unexpected control request #%d", c))
			}
		}
	}
	changes <- svc.Status{State: svc.StopPending}
	return
}

func runService(name string, isDebug bool) {
	var err error
	if isDebug {
		elog = debug.New(name)
	} else {
		elog, err = eventlog.Open(name)
		if err != nil {
			return
		}
	}
	defer elog.Close()

	elog.Info(1, fmt.Sprintf("starting %s service", name))
	run := svc.Run
	if isDebug {
		run = debug.Run
	}
	err = run(name, &service{})
	if err != nil {
		elog.Error(1, fmt.Sprintf("%s service failed: %v", name, err))
		return
	}
	elog.Info(1, fmt.Sprintf("%s service stopped", name))
}

func createiBases() {
	defer mutex.Unlock()

	if needCreate() {
		err := OneCIBasesCreator.Create(pathLst, pathiBases)
		if err != nil {
			elog.Error(1, err.Error())
		}
	}
}

func needCreate() bool {

	need := false

	if lastModifiedLst == nil {
		lastModifiedLst = make(map[string]time.Time, 0)
		for _, path := range pathLst {
			lastModifiedLst[path] = time.Unix(0, 0)
		}
	}

	for _, path := range pathLst {
		stat, err := os.Stat(path)
		if err != nil {
			elog.Error(1, err.Error())
			continue
		}
		modTime := stat.ModTime()
		if modTime != lastModifiedLst[path] {
			need = true
			lastModifiedLst[path] = modTime
		}

	}

	return need

}

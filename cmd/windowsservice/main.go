//go:build windows

package main

import (
	"io"
	"os"
	"strings"

	welog "github.com/korableg/V8I.Manager/internal/windowseventlogadapter"
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"

	"golang.org/x/sys/windows/svc"
)

const svcName = "V8I.Manager"
const svcDescription = "Observes changes in the files of the 1S server list and creates a list of databases for the user"

func main() {

	//var Version = "dev"

	isService, err := svc.IsWindowsService()
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to determine that the running application is running as a service"))
	}

	if isService {

		winloginfo, err := welog.New(svcName, welog.Info, 1)
		if err != nil {
			log.Fatal("windows event log INFO level haven't initialized")
		}

		defer winloginfo.Close()

		winlogwarning, err := welog.New(svcName, welog.Warning, 2)
		if err != nil {
			log.Fatal("windows event log WARNING level haven't initialized")
		}

		defer winlogwarning.Close()

		winlogerror, err := welog.New(svcName, welog.Error, 3)
		if err != nil {
			log.Fatal("windows event log ERROR level haven't initialized")
		}

		defer winlogerror.Close()

		log.SetOutput(io.Discard)
		log.AddHook(&writer.Hook{
			Writer: winloginfo,
			LogLevels: []log.Level{
				log.InfoLevel,
				log.DebugLevel,
				log.TraceLevel,
			},
		})
		log.AddHook(&writer.Hook{
			Writer: winlogwarning,
			LogLevels: []log.Level{
				log.WarnLevel,
			},
		})
		log.AddHook(&writer.Hook{
			Writer: winlogerror,
			LogLevels: []log.Level{
				log.ErrorLevel,
				log.FatalLevel,
				log.PanicLevel,
			},
		})

		runService()

		return
	}

	cmd := strings.ToLower(os.Args[1])
	switch cmd {
	case "install":
		err = installService(os.Args[2:])
	case "remove":
		err = removeService()
	case "start":
		err = startService()
	case "stop":
		err = controlService(svc.Stop, svc.Stopped)
	default:
		// Открыть в пакетном режиме
	}

	if err != nil {
		log.Fatal(err)
	}

	return

}

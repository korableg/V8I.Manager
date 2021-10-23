package main

import (
	"fmt"
	"io"

	"github.com/korableg/V8I.Manager/internal/globals"
	welog "github.com/korableg/V8I.Manager/internal/windowseventlogadapter"
	"github.com/korableg/V8I.Manager/internal/windowsservice"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"
	"github.com/spf13/cobra"
	"golang.org/x/sys/windows/svc"
)

func _installCmd(cmd *cobra.Command, args []string) error {

	err := windowsservice.Install(globals.AppName, globals.Description, _cfgFile)
	if err != nil {
		return err
	}

	log.Info("the service was successfull installed")
	return nil

}

func _removeCmd(cmd *cobra.Command, args []string) error {
	err := windowsservice.Remove(globals.AppName)
	if err != nil {
		return err
	}
	log.Info("the service was successfull removed")
	return nil
}

func initConfigOs() error {

	isService, err := svc.IsWindowsService()
	if err != nil {
		return err
	}

	cfg.SetIsWindowsService(isService)

	return nil

}

func addCommands(cmd *cobra.Command) {

	cmd.AddCommand(&cobra.Command{
		Use:   "install",
		Short: fmt.Sprintf("installs %s as a windows service", globals.AppName),
		RunE:  _installCmd,
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "remove",
		Short: fmt.Sprintf("removes %s from list of windows services", globals.AppName),
		RunE:  _removeCmd,
	})

}

func runService() error {

	loggerInfo, err := welog.New(globals.AppName, welog.Info, 1)
	if err != nil {
		return errors.Wrap(err, "windows event log INFO level couldn't initialized")
	}
	defer loggerInfo.Close()

	loggerWarning, err := welog.New(globals.AppName, welog.Warning, 2)
	if err != nil {
		return errors.Wrap(err, "windows event log WARNING level couldn't initialized")
	}
	defer loggerWarning.Close()

	loggerError, err := welog.New(globals.AppName, welog.Error, 3)
	if err != nil {
		return errors.Wrap(err, "windows event log ERROR level couldn't initialized")
	}
	defer loggerError.Close()

	log.SetOutput(io.Discard)
	log.AddHook(&writer.Hook{
		Writer: loggerInfo,
		LogLevels: []log.Level{
			log.InfoLevel,
			log.DebugLevel,
			log.TraceLevel,
		},
	})
	log.AddHook(&writer.Hook{
		Writer: loggerWarning,
		LogLevels: []log.Level{
			log.WarnLevel,
		},
	})
	log.AddHook(&writer.Hook{
		Writer: loggerError,
		LogLevels: []log.Level{
			log.ErrorLevel,
			log.FatalLevel,
			log.PanicLevel,
		},
	})

	log.Info("starting service")

	service := windowsservice.New(cfg)

	err = svc.Run(globals.AppName, service)
	if err != nil {
		return err
	}

	log.Info("service stopped")

	return nil

}

//go:build windows

package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/korableg/V8I.Manager/internal/globals"
	welog "github.com/korableg/V8I.Manager/internal/windowseventlogadapter"
	"github.com/korableg/V8I.Manager/internal/worker"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"golang.org/x/sys/windows/svc"
)

const (
	_cfgFlag = "cfg"
	_lstFlag = "lst"
	_v8iFlag = "v8i"
)

var Version = "dev"

func main() {

	cmd := &cobra.Command{
		Use:     filepath.Base(os.Args[0]),
		Short:   globals.Description,
		Version: Version,
		RunE: func(cmd *cobra.Command, args []string) error {

			isService, err := svc.IsWindowsService()
			if err != nil {
				return errors.Wrap(err, "couldn't determine that the application is a service")
			}

			if isService {
				err = initService()
				if err != nil {
					return err
				}
			}

			err = prepareConfig()
			if err != nil {
				return err
			}

			err = checkSources()
			if err != nil {
				return err
			}

			if isService {
				return runService()
			}

			lst := viper.GetStringSlice(_lstFlag)
			v8i := viper.GetStringSlice(_v8iFlag)

			err = worker.LstToV8i(lst, v8i)
			if err != nil {
				return err
			}

			log.Infof("v8i files by paths %s successfully created", v8i)
			return nil

		},
	}

	cmd.SetHelpCommand(&cobra.Command{Hidden: true})

	cmd.AddCommand(&cobra.Command{
		Use:   "install",
		Short: fmt.Sprintf("installs %s as a windows service", globals.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {

			err := prepareConfig()
			if err != nil {
				return err
			}

			err = installService()
			if err != nil {
				return err
			}

			log.Info("the service was successfull installed")
			return nil

		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "remove",
		Short: fmt.Sprintf("removes %s from list of windows services", globals.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := removeService()
			if err != nil {
				return err
			}
			log.Info("the service was successfull removed")
			return nil
		},
	})

	cmd.CompletionOptions.DisableDefaultCmd = true

	cmd.PersistentFlags().StringP(_cfgFlag, "c", "", "file with the application's settings")
	cmd.PersistentFlags().StringSliceP(_lstFlag, "l", nil, "comma-separated list of lst files")
	cmd.PersistentFlags().StringSliceP(_v8iFlag, "v", nil, "comma-separated list of v8i files")

	viper.BindPFlag(_cfgFlag, cmd.PersistentFlags().Lookup(_cfgFlag))
	viper.BindPFlag(_lstFlag, cmd.PersistentFlags().Lookup(_lstFlag))
	viper.BindPFlag(_v8iFlag, cmd.PersistentFlags().Lookup(_v8iFlag))

	err := cmd.Execute()
	if err != nil {
		log.Fatal(err)
	}

}

func checkSources() error {

	lstFlagValue := viper.GetStringSlice(_lstFlag)
	v8iFlagValue := viper.GetStringSlice(_v8iFlag)

	for _, fileName := range lstFlagValue {
		_, err := os.Stat(fileName)
		if err != nil {
			return err
		}
	}

	for _, fileName := range v8iFlagValue {
		_, err := os.Stat(fileName)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
		if os.IsNotExist(err) {
			f, err := os.Create(fileName)
			if err != nil {
				return err
			}
			err = f.Close()
			if err != nil {
				return err
			}
			err = os.Remove(fileName)
			if err != nil {
				return err
			}
		}
	}

	return nil

}

func initService() error {

	winloginfo, err := welog.New(globals.AppName, welog.Info, 1)
	if err != nil {
		return errors.Wrap(err, "windows event log INFO level haven't initialized")
	}

	winlogwarning, err := welog.New(globals.AppName, welog.Warning, 2)
	if err != nil {
		return errors.Wrap(err, "windows event log WARNING level haven't initialized")
	}

	winlogerror, err := welog.New(globals.AppName, welog.Error, 3)
	if err != nil {
		return errors.Wrap(err, "windows event log ERROR level haven't initialized")
	}

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

	return nil

}

func prepareConfig() error {

	cfgFlagValue := viper.GetString(_cfgFlag)
	lstFlagValue := viper.GetStringSlice(_lstFlag)
	v8iFlagValue := viper.GetStringSlice(_v8iFlag)

	if !(len(lstFlagValue) > 0 && len(v8iFlagValue) > 0) {

		if len(cfgFlagValue) == 0 {
			return errors.New("should be determine correct path in --cfg flag or --lst and --v8i flags")
		}

		viper.SetConfigFile(cfgFlagValue)
		err := viper.ReadInConfig()
		if err != nil {
			return err
		}

		log.Debugf("config file %s has been read", viper.ConfigFileUsed())
	}

	return nil

}

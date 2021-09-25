//go:build windows

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/korableg/V8I.Manager/internal/globals"
	welog "github.com/korableg/V8I.Manager/internal/windowseventlogadapter"
	"github.com/spf13/viper"
	"golang.org/x/sys/windows/svc/mgr"
)

func exePath() (string, error) {
	prog := os.Args[0]
	p, err := filepath.Abs(prog)
	if err != nil {
		return "", err
	}
	fi, err := os.Stat(p)
	if err == nil {
		if !fi.Mode().IsDir() {
			return p, nil
		}
		err = fmt.Errorf("%s is directory", p)
	}
	if filepath.Ext(p) == "" {
		p += ".exe"
		fi, err := os.Stat(p)
		if err == nil {
			if !fi.Mode().IsDir() {
				return p, nil
			}
			err = fmt.Errorf("%s is directory", p)
		}
	}
	return "", err
}

func installService() error {
	exepath, err := exePath()
	if err != nil {
		return err
	}

	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(globals.AppName)
	if err == nil {
		s.Close()
		return fmt.Errorf("service %s already exists", globals.AppName)
	}

	args, err := serviceArgs()
	if err != nil {
		return err
	}

	s, err = m.CreateService(
		globals.AppName,
		exepath,
		mgr.Config{
			DisplayName: globals.AppName,
			Description: globals.Description,
			StartType:   mgr.StartAutomatic,
		},
		args...)

	if err != nil {
		return err
	}
	defer s.Close()
	err = welog.Install(globals.AppName)
	if err != nil {
		s.Delete()
		return fmt.Errorf("SetupEventLogSource() failed: %s", err)
	}
	return nil
}

func removeService() error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(globals.AppName)
	if err != nil {
		return fmt.Errorf("service %s is not installed", globals.AppName)
	}
	defer s.Close()
	err = s.Delete()
	if err != nil {
		return err
	}
	err = welog.Remove(globals.AppName)
	if err != nil {
		return fmt.Errorf("RemoveEventLogSource() failed: %s", err)
	}
	return nil
}

func serviceArgs() ([]string, error) {

	serviceArgs := make([]string, 0)
	cfgFile := viper.ConfigFileUsed()
	if len(cfgFile) > 0 {
		cfgPath, err := filepath.Abs(cfgFile)
		if err != nil {
			return nil, err
		}
		serviceArgs = append(serviceArgs, "--cfg", cfgPath)
		return serviceArgs, nil
	}

	serviceArgs = append(serviceArgs,
		"--lst",
		strings.Join(viper.GetStringSlice(_lstFlag), ","),
		"--v8i",
		strings.Join(viper.GetStringSlice(_v8iFlag), ","),
	)

	return serviceArgs, nil

}

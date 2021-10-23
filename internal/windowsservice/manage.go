//go:build windows

package windowsservice

import (
	"fmt"
	"os"
	"path/filepath"

	welog "github.com/korableg/V8I.Manager/internal/windowseventlogadapter"
	"github.com/pkg/errors"
	"golang.org/x/sys/windows/svc/mgr"
)

func Install(name, description, configPath string) error {
	exepath, err := exePath()
	if err != nil {
		return err
	}

	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(name)
	if err == nil {
		s.Close()
		return fmt.Errorf("service %s already exists", name)
	}

	args, err := serviceArgs(configPath)
	if err != nil {
		return err
	}

	s, err = m.CreateService(
		name,
		exepath,
		mgr.Config{
			DisplayName: name,
			Description: description,
			StartType:   mgr.StartAutomatic,
		},
		args...)

	if err != nil {
		return err
	}
	defer s.Close()
	err = welog.Install(name)
	if err != nil {
		s.Delete()
		return fmt.Errorf("SetupEventLogSource() failed: %s", err)
	}
	return nil
}

func Remove(name string) error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(name)
	if err != nil {
		return fmt.Errorf("service %s is not installed", name)
	}
	defer s.Close()
	err = s.Delete()
	if err != nil {
		return err
	}
	err = welog.Remove(name)
	if err != nil {
		return fmt.Errorf("RemoveEventLogSource() failed: %s", err)
	}
	return nil
}

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

func serviceArgs(configPath string) ([]string, error) {

	serviceArgs := make([]string, 2)
	if configPath == "" {
		return nil, errors.New("cfg file doesn't filled")
	}

	cfgPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, err
	}
	serviceArgs[0] = "--cfg"
	serviceArgs[1] = cfgPath

	return serviceArgs, nil

}

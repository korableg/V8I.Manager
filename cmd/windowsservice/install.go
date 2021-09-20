//go:build windows

package main

import (
	"fmt"
	"os"
	"path/filepath"

	welog "github.com/korableg/V8I.Manager/internal/windowseventlogadapter"
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

func installService(args []string) error {
	exepath, err := exePath()
	if err != nil {
		return err
	}
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(svcName)
	if err == nil {
		s.Close()
		return fmt.Errorf("service %s already exists", svcName)
	}
	s, err = m.CreateService(
		svcName,
		exepath,
		mgr.Config{
			DisplayName: svcName,
			Description: svcDescription,
			StartType:   mgr.StartAutomatic,
		},
		args...)
	if err != nil {
		return err
	}
	defer s.Close()
	err = welog.Install(svcName)
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
	s, err := m.OpenService(svcName)
	if err != nil {
		return fmt.Errorf("service %s is not installed", svcName)
	}
	defer s.Close()
	err = s.Delete()
	if err != nil {
		return err
	}
	err = welog.Remove(svcName)
	if err != nil {
		return fmt.Errorf("RemoveEventLogSource() failed: %s", err)
	}
	return nil
}

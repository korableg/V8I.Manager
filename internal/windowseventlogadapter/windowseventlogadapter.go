package windowseventlogadapter

import (
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"
)

const (
	Info    = eventlog.Info
	Warning = eventlog.Warning
	Error   = eventlog.Error
)

func Install(source string) error {
	return eventlog.InstallAsEventCreate(source, eventlog.Info|eventlog.Warning|eventlog.Error)
}

func Remove(source string) error {
	return eventlog.Remove(source)
}

type WindowsEventLogAdapter struct {
	elog  debug.Log
	level int
	eid   uint32
}

func New(source string, level int, eid uint32) (*WindowsEventLogAdapter, error) {

	elog, err := eventlog.Open(source)
	if err != nil {
		return nil, err
	}

	w := &WindowsEventLogAdapter{
		elog:  elog,
		level: level,
		eid:   eid,
	}

	return w, nil

}

func (w *WindowsEventLogAdapter) Write(p []byte) (int, error) {

	var err error
	var msg = string(p)

	switch w.level {
	case Info:
		err = w.elog.Info(w.eid, msg)
	case Warning:
		err = w.elog.Warning(w.eid, msg)
	case Error:
		err = w.elog.Error(w.eid, msg)
	}

	if err != nil {
		return 0, err
	}

	return len(p), nil

}

func (w *WindowsEventLogAdapter) Close() error {
	return w.elog.Close()
}

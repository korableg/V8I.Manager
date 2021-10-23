package config

import "github.com/korableg/V8I.Manager/pkg/v8i/v8iwriter"

type Config interface {
	Lsts() []string
	SetLsts(lsts []string) error
	V8is() []v8iwriter.V8IWriter
	SetV8is(v8is []v8iwriter.V8IWriter) error
	IsWindowsService() bool
	SetIsWindowsService(isService bool)
}

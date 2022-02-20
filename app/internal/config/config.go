package config

import (
	"fmt"
	"github.com/korableg/V8I.Manager/app/api/user/auth"
	"gopkg.in/yaml.v3"
	"io/ioutil"

	"github.com/korableg/V8I.Manager/app/internal/sqlitedb"
	"github.com/korableg/V8I.Manager/app/internal/transport/httpserver"
)

type (
	Config struct {
		Sqlite sqlitedb.Config   `yaml:"sqlite"`
		Http   httpserver.Config `yaml:"http"`
		Auth   auth.Config       `yaml:"auth"`
	}
)

func NewConfig(path string) (Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("error config file: %w", err)
	}

	c := Config{}
	if err = yaml.Unmarshal(data, &c); err != nil {
		return Config{}, fmt.Errorf("unmarshal config: %w", err)
	}

	return c, nil
}

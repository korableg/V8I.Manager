package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"

	"github.com/go-playground/validator/v10"
	"github.com/korableg/V8I.Manager/app/api/user/auth"
	"github.com/korableg/V8I.Manager/app/internal/sqlitedb"
	"github.com/korableg/V8I.Manager/app/internal/transport/httpserver"
)

type (
	Config struct {
		Sqlite sqlitedb.Config   `yaml:"sqlite" validate:"required,dive"`
		Http   httpserver.Config `yaml:"http" validate:"required,dive"`
		Auth   auth.Config       `yaml:"auth" validate:"required,dive"`
	}
)

func NewConfig(path string, validate *validator.Validate) (Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("error config file: %w", err)
	}

	c := Config{}
	if err = yaml.Unmarshal(data, &c); err != nil {
		return Config{}, fmt.Errorf("unmarshal config: %w", err)
	}

	if err = validate.Struct(c); err != nil {
		return Config{}, fmt.Errorf("validate config: %w", err)
	}

	return c, nil
}

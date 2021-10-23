package config

import (
	"os"

	"github.com/korableg/V8I.Manager/pkg/v8i/v8iwriter"
	"github.com/korableg/V8I.Manager/pkg/v8i/v8iwriter/v8ifilewriter"
	"github.com/spf13/viper"
)

const (
	_isWindowsService = "isWindowsService"
	_v8is             = "v8i"
	_lsts             = "lst"
)

type viperConfig struct {
	v *viper.Viper
}

func New(path string) (*viperConfig, error) {

	v := viper.New()

	if path != "" {
		v.SetConfigFile(path)
	} else {
		v.AddConfigPath(".")
		v.SetConfigType("yaml")
		v.SetConfigName("config")
	}

	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	vc := &viperConfig{
		v: v,
	}

	err = vc.checkLstSources()
	if err != nil {
		return nil, err
	}

	vc.prepareConfig()

	return vc, nil

}

func (v *viperConfig) V8is() []v8iwriter.V8IWriter {
	return v.v.Get(_v8is).([]v8iwriter.V8IWriter)
}

func (v *viperConfig) SetV8is(v8is []v8iwriter.V8IWriter) error {
	v.v.Set(_v8is, v8is)
	return v.v.WriteConfig()
}

func (v *viperConfig) Lsts() []string {
	return v.v.GetStringSlice(_lsts)
}

func (v *viperConfig) SetLsts(lsts []string) error {
	v.v.Set(_lsts, lsts)
	return v.v.WriteConfig()
}

func (v *viperConfig) IsWindowsService() bool {
	return v.v.GetBool(_isWindowsService)
}

func (v *viperConfig) SetIsWindowsService(isService bool) {
	v.v.Set(_isWindowsService, isService)
}

func (v *viperConfig) checkLstSources() error {

	lsts := v.v.GetStringSlice(_lsts)

	for _, fileName := range lsts {
		_, err := os.Stat(fileName)
		if err != nil {
			return err
		}
	}

	return nil

}

func (v *viperConfig) prepareConfig() {

	v8is := v.v.GetStringSlice(_v8is)
	v8iWriters := make([]v8iwriter.V8IWriter, len(v8is))
	for i, v := range v8is {
		v8iWriters[i] = v8ifilewriter.New(v)
	}

	v.v.Set(_v8is, v8iWriters)

}

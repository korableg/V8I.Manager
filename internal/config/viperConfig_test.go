package config

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const configPath = "../../assets/config_example.yaml"

func TestViperConfig(t *testing.T) {

	var cfg Config

	absConfigPath, err := filepath.Abs(configPath)
	require.Equal(t, nil, err)

	cfg, err = New(absConfigPath)
	require.Equal(t, nil, err)

	lsts := cfg.Lsts()
	assert.Equal(t, 1, len(lsts))

	v8is := cfg.V8is()
	assert.Equal(t, 1, len(v8is))

}

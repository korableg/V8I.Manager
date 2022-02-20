package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	const path = "./config_test.yaml"

	c, err := NewConfig(path)
	require.Nil(t, err)

	assert.Equal(t, c.Sqlite.Path, "./db.db")
	assert.Equal(t, c.Http.Address, "localhost")
	assert.Equal(t, c.Http.Port, 8080)
	assert.Equal(t, c.Auth.Secret, "lenrkaeirhbekjvnaeir")
}

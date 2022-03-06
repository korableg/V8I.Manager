package dbbuilder

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDBBuilder(t *testing.T) {
	const path = "../../../test/test.lst"

	b, err := NewBuilder()
	require.Nil(t, err)

	dbs, err := b.Build(path)
	assert.Nil(t, err)
	assert.NotEqual(t, len(dbs), 0)
}

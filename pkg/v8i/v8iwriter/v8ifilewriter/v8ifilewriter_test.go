package v8ifilewriter

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/korableg/V8I.Manager/pkg/clusterdb"
	"github.com/korableg/V8I.Manager/pkg/v8i/v8ibuilder"
	"github.com/korableg/V8I.Manager/pkg/v8i/v8iwriter"
	"github.com/stretchr/testify/assert"
)

func TestV8ifilewriter(t *testing.T) {

	dbs := make([]*clusterdb.ClusterDB, 1)
	dbs[0] = &clusterdb.ClusterDB{
		ID:          "1124234",
		Server:      "test",
		Ref:         "testref",
		Description: "test base",
		Name:        "TEST",
		Folder:      "TEST",
	}

	v8iBytes, err := v8ibuilder.Build(dbs...)
	assert.Equal(t, err, nil)

	var w v8iwriter.V8IWriter

	v8ifilename := filepath.FromSlash("../../../../test/test1.v8i")

	w = New(v8ifilename)
	_, err = w.Write(v8iBytes)
	assert.Equal(t, err, nil)

	err = os.Remove(v8ifilename)
	assert.Equal(t, err, nil)

}

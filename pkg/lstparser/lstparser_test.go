package lstparser

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const pathtolst = "../../test/test.lst"

func TestLstParser(t *testing.T) {

	lstdata, err := os.ReadFile(pathtolst)
	if err != nil {
		t.Fatal(err)
	}

	dbs, err := Parse(lstdata)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, len(dbs), 3)

}

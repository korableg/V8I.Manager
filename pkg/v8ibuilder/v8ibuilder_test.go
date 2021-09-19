package v8ibuilder

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/dimchansky/utfbom"
	"github.com/korableg/V8I.Manager/pkg/lstparser"
	"github.com/stretchr/testify/assert"
)

const pathtolst = "../../test/test.lst"
const pathtoV8I = "../../test/ibases.v8i"

func TestV8IBuilder(t *testing.T) {

	lstdata, err := os.ReadFile(pathtolst)
	assert.Equal(t, err, nil)

	file, err := os.Open(pathtoV8I)
	assert.Equal(t, err, nil)

	v8itestdata, err := ioutil.ReadAll(utfbom.SkipOnly(file))
	assert.Equal(t, err, nil)

	clusterDBs, err := lstparser.Parse(lstdata)
	assert.Equal(t, err, nil)

	v8i, err := Build(clusterDBs...)
	assert.Equal(t, err, nil)

	assert.Equal(t, v8i, v8itestdata)

}

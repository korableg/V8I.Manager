package OneCIBasesCreator

import (
	"testing"
)

var PathLst = []string {"./test/test.lst"}
var PathIBases = []string {"./test/ibases.v8i"}

func TestAgentIBases_Create(t *testing.T) {
	err := Create(PathLst, PathIBases)
	if err != nil {
		t.Fatal(err)
	}
}
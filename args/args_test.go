package args

import (
	"errors"
	"testing"
)

func TestPathIBases(t *testing.T) {

	_, err := PathIBases()
	if err != nil && !errors.Is(err, ErrPathIBasesEmpty) {
		t.Error(err)
	}

}

func TestPathLst(t *testing.T) {

	_, err := PathLst()
	if err != nil && err != ErrPathLstEmpty {
		t.Error(err)
	}

}

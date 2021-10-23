package v8iwriter

import (
	"fmt"
	"io"
)

type V8IWriter interface {
	io.Writer
	fmt.Stringer
}

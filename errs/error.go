package errs

import (
	"fmt"
	"path/filepath"
	"runtime"
)

type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

func New(text string) error {
	_, fp, line, ok := runtime.Caller(1)
	if ok {
		return &errorString{fmt.Sprintf("%s:%d %s", filepath.Base(fp), line, text)}
	} else {
		return &errorString{text}
	}
}

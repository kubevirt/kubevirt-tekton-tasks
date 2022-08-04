package zerrors

import (
	"fmt"
)

type softError struct {
	msg  string
	Soft bool
}

func (e softError) Error() string {
	return e.msg
}

func (e softError) IsSoft() bool {
	return e.Soft
}

func NewSoftError(format string, a ...interface{}) error {
	return &softError{Soft: true, msg: fmt.Sprintf(format, a...)}
}

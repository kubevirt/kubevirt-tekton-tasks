package zerrors

import (
	"fmt"
)

type MissingRequiredError struct {
	msg  string
	Soft bool
}

func (e MissingRequiredError) Error() string {
	return e.msg
}

func (e MissingRequiredError) IsSoft() bool {
	return e.Soft
}

func NewMissingRequiredError(format string, a ...interface{}) error {
	return &MissingRequiredError{Soft: true, msg: fmt.Sprintf(format, a...)}
}

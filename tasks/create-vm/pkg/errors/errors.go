package errors

import (
	"fmt"
)

type MissingRequiredError struct {
	s string
}

func (e MissingRequiredError) Error() string {
	return e.s
}

func NewMissingRequiredArgError(name string) error {
	return &MissingRequiredError{s: fmt.Sprintf("missing required argument -%v", name)}
}

func NewMissingRequiredError(s string) error {
	return &MissingRequiredError{s}
}

type NotFoundError struct {
	s string
}

func (e NotFoundError) Error() string {
	return e.s
}

func NewNotFoundError(s string) error {
	return &MissingRequiredError{s}
}

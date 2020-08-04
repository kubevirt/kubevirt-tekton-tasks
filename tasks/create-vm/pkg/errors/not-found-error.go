package errors

type NotFoundError struct {
	msg string
}

func (e NotFoundError) Error() string {
	return e.msg
}

func NewNotFoundError(msg string) error {
	return &NotFoundError{msg}
}

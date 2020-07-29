package errors

import "strings"

type MultiError struct {
	errors []error
}

func (e MultiError) Error() string {
	var str strings.Builder
	if e.errors != nil {
		for _, err := range e.errors {
			if err != nil {
				errMsg := err.Error()
				if errMsg != "" {
					str.WriteString(errMsg)
					if errMsg[len(errMsg)-1] != '\n' {
						str.WriteRune('\n')
					}
				}

			}
		}
	}

	return str.String()
}

func (e MultiError) Get(idx int) error {
	return e.errors[idx]
}

func GetErrorFromMultiError(err error, idx int) error {
	if err == nil {
		return nil
	}

	var multiErr, ok = err.(*MultiError)
	if !ok {
		return err
	}
	if idx < 0 {
		return nil
	}
	return multiErr.Get(idx)
}

func NewMultiErrorOrNil(errs []error) error {
	if errs == nil {
		return nil
	}

	isEmpty := true
	for _, err := range errs {
		if err != nil {
			isEmpty = false
			break
		}
	}

	if isEmpty {
		return nil
	}

	return &MultiError{errs}
}

func NewMultiFilteredErrorOrNil(errs []error) error {
	if errs == nil {
		return nil
	}

	var filteredErrors []error

	for _, err := range errs {
		if err != nil {
			filteredErrors = append(filteredErrors, err)
		}
	}

	return NewMultiErrorOrNil(filteredErrors)
}

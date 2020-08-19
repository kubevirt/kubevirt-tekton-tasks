package zerrors

import (
	"strings"
)

type MultiError struct {
	keyOrder    []string
	errorMap    map[string]error
	short       bool
	shortErrMsg string
}

func NewMultiError() *MultiError {
	return &MultiError{}
}

func (e *MultiError) init() {
	if e == nil {
		*e = MultiError{}
	}
}

func (e MultiError) Error() string {
	if e.IsEmpty() {
		return ""
	}

	var str strings.Builder
	if e.short {
		str.WriteString(e.shortErrMsg)
		str.WriteRune(' ')
		str.WriteString(strings.Join(e.keyOrder, ", "))
	} else {
		for _, key := range e.keyOrder {
			if err, ok := e.errorMap[key]; ok {
				if errMsg := err.Error(); errMsg != "" {
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

func (e *MultiError) IsEmpty() bool {
	return e == nil || len(e.keyOrder) <= 0 || len(e.errorMap) <= 0
}

func (e *MultiError) Len() int {
	if e == nil {
		return 0
	}
	return len(e.keyOrder)
}

func (e *MultiError) IsSoft() bool {
	if !e.IsEmpty() {
		for _, key := range e.keyOrder {
			if err, ok := e.errorMap[key]; ok {
				if !IsErrorSoft(err) {
					return false
				}
			}
		}
	}

	return true
}

func (e *MultiError) Add(name string, err error) {
	e.init()

	if name != "" && err != nil {
		if e.errorMap == nil {
			e.errorMap = make(map[string]error, 1)
		}

		e.keyOrder = append(e.keyOrder, name)
		e.errorMap[name] = err
	}
}

// Add but for chaining
func (e *MultiError) AddC(name string, err error) *MultiError {
	e.Add(name, err)
	return e
}

func (e *MultiError) Get(name string) error {
	if e.IsEmpty() {
		return nil
	}
	return e.errorMap[name]
}

func (e *MultiError) AsOptional() error {
	if !e.IsEmpty() {
		return e
	}
	return nil
}

func (e *MultiError) print(isShort bool, msg string) *MultiError {
	e.init()
	e.short = isShort
	e.shortErrMsg = msg
	return e
}

func (e *MultiError) ShortPrint(msg string) *MultiError {
	return e.print(true, msg)
}

func (e *MultiError) LongPrint() *MultiError {
	return e.print(false, "")
}

func GetErrorFromMultiError(err error, name string) error {
	if err == nil || name == "" {
		return nil
	}

	var multiErr, ok = err.(*MultiError)
	if !ok {
		return err
	}
	return multiErr.Get(name)
}

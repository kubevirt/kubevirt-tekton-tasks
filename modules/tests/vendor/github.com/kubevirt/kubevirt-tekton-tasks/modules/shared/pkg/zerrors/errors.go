package zerrors

import "k8s.io/apimachinery/pkg/api/errors"

type SoftError interface {
	IsSoft() bool
}

func IsErrorSoft(err error) bool {
	if err == nil {
		return false
	}
	var softErr, ok = err.(SoftError)

	return ok && softErr.IsSoft()
}

func IsStatusError(err error, allowedStatusCodes ...int32) bool {
	if err == nil {
		return false
	}

	if statusErr, ok := err.(*errors.StatusError); ok {
		for _, code := range allowedStatusCodes {
			if code == statusErr.ErrStatus.Code {
				return true
			}
		}
	}

	return false
}

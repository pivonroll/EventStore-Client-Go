package errors

import (
	stdErrors "errors"
	"fmt"
)

type ErrorCode string

type Error interface {
	Error() string
	Code() ErrorCode
}

type errorImpl struct {
	err  error
	code ErrorCode
}

func (r errorImpl) Error() string {
	if r.err != nil {
		return fmt.Sprintf("Error code: %s. Reason %v", r.code, r.err)
	}

	return fmt.Sprintf("Error code: %s", r.code)
}

func (r errorImpl) Code() ErrorCode {
	return r.code
}

func NewErrorCodeMsg(code ErrorCode, errorString string) Error {
	return errorImpl{code: code, err: stdErrors.New(errorString)}
}

func NewErrorCodeMsgf(code ErrorCode, msg string, params ...interface{}) Error {
	return errorImpl{code: code, err: stdErrors.New(fmt.Sprintf(msg, params...))}
}

func NewErrorCode(code ErrorCode) Error {
	return errorImpl{
		code: code,
	}
}

func NewError(code ErrorCode, err error) Error {
	return errorImpl{
		code: code,
		err:  err,
	}
}

package errors

import (
	stderrors "errors"
	"fmt"
)

type Code string

const (
	CodeUnknown         Code = "unknown"
	CodeInvalidArgument Code = "invalid_argument"
	CodeBrowserFailed   Code = "browser_failed"
	CodeCaptureFailed   Code = "capture_failed"
	CodeWriteFailed     Code = "write_failed"
)

type Error struct {
	Code    Code
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}

	if e.Err == nil {
		return e.Message
	}

	return fmt.Sprintf("%s: %v", e.Message, e.Err)
}

func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Err
}

func New(code Code, message string) error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func Wrap(code Code, message string, err error) error {
	if err == nil {
		return New(code, message)
	}

	return &Error{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func CodeOf(err error) Code {
	if err == nil {
		return CodeUnknown
	}

	var appErr *Error
	if stderrors.As(err, &appErr) {
		return appErr.Code
	}

	return CodeUnknown
}

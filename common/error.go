package common

import "fmt"

type Error struct {
	msg   string
	cause error
}

func (e *Error) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("Error: %s.\nCaused by: %s", e.msg, e.cause.Error())
	} else {
		return fmt.Sprintf("Error: %s.", e.msg)
	}
}

func NewError(msg string) *Error {
	return &Error{msg, nil}
}

func NewErrorWithCause(msg string, cause error) *Error {
	return &Error{msg, cause}
}

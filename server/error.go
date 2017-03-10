package server

import "fmt"

const (
	ErrServerFailure int32 = -10000
	ErrRankNotFound  int32 = -10001
)

type Error struct {
	PrevErr error
	Cause   string
}

func NewError(cause string, prevErr error) *Error {
	return &Error{
		PrevErr: prevErr,
		Cause:   cause,
	}
}

func (e Error) Error() string {
	return fmt.Sprintf("%s: %v", e.Cause, e.PrevErr)
}

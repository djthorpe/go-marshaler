package marshaler

import "fmt"

///////////////////////////////////////////////////////////////////////////////
// TYPES

type Error int

///////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	ErrSuccess Error = iota
	ErrBadParameter
)

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (e Error) Error() string {
	switch e {
	case ErrSuccess:
		return "ErrSuccess"
	case ErrBadParameter:
		return "ErrBadParameter"
	default:
		return "[?? Invalid Error value]"
	}
}

func (e Error) With(args ...interface{}) error {
	return fmt.Errorf("%s: %w", fmt.Sprint(args...), e)
}

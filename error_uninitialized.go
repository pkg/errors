package errors

import "fmt"

// ErrorUninitialized is returned when a required dependency for this action has not been set up properly.
type ErrorUninitialized struct {
	BaseError
}

// Uninitialized returns a new instance of ErrorUninitialized.
func Uninitialized(format string, a ...interface{}) error {
	return ErrorUninitialized{
		NewBase("uninitialized :: " + fmt.Sprintf(format, a...)),
	}
}

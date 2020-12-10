package errors

import (
	"fmt"
)

// ErrorAlreadyExists is returned when a new resource can't be created due to a conflict with an existing resource.
type ErrorAlreadyExists struct {
	BaseError
}

// AlreadyExists returns a new instance of ErrorAlreadyExists.
func AlreadyExists(format string, a ...interface{}) error {
	return ErrorAlreadyExists{
		NewBase("already exists :: " + fmt.Sprintf(format, a...)),
	}
}

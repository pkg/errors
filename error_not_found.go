package errors

import (
	"fmt"
)

// ErrorNotFound is return when expected object is not present
type ErrorNotFound struct {
	BaseError
}

// NotFound returns a new instance of ErrorNotFound
func NotFound(message string) error {
	return ErrorNotFound{
		NewBase(fmt.Sprintf("not found :: %s", message)),
	}
}

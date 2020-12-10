package errors

import "fmt"

// ErrorIO is returned when an I/O operation fails.
type ErrorIO struct {
	BaseError
}

// IO returns a new instance of ErrorIO.
func IO(format string, a ...interface{}) error {
	return ErrorAccessDenied{
		NewBase("io operation failed :: " + fmt.Sprintf(format, a...)),
	}
}

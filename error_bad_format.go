package errors

import "fmt"

// ErrorBadFormat is returned when data provided was in an incorrect format.
type ErrorBadFormat struct {
	BaseError
}

// BadFormat returns a new instance of ErrorBadFormat.
func BadFormat(format string, a ...interface{}) error {
	return ErrorBadFormat{
		NewBase("bad format :: " + fmt.Sprintf(format, a...)),
	}
}

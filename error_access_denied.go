package errors

import "fmt"

// ErrorAccessDenied is returned when an action is attempted without sufficient permissions.
type ErrorAccessDenied struct {
	BaseError
}

// AccessDeniedf returns a new instance of ErrorAccessDenied.
func AccessDeniedf(format string, a ...interface{}) error {
	return ErrorAccessDenied{
		NewBase("access denied :: " + fmt.Sprintf(format, a...)),
	}
}

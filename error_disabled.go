package errors

import "fmt"

// ErrorDisabled is returned when the action being attempted is disabled but not necessarily inherently to the user.
type ErrorDisabled struct {
	BaseError
}

// Disabled returns a new instance of ErrorDisabled.
func Disabled(format string, a ...interface{}) error {
	return ErrorDisabled{
		NewBase("disabled :: " + fmt.Sprintf(format, a...)),
	}
}

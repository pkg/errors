package errors

import "fmt"

// ErrorLimited is returned when an action is attempted but has been temporarily limited (rate limited or otherwise).
type ErrorLimited struct {
	BaseError
}

// Limited returns a new instance of ErrorLimited.
func Limited(format string, a ...interface{}) error {
	return ErrorLimited{
		NewBase("limited :: " + fmt.Sprintf(format, a...)),
	}
}

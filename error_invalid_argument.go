package errors

import "fmt"

// ErrorInvalidArgument is returned when provided argument is not valid
type ErrorInvalidArgument struct {
	BaseError
}

// InvalidArgument returns a new instance of ErrorInvalidArgument
func InvalidArgument(argument, reason string) error {
	return ErrorInvalidArgument{
		NewBase(fmt.Sprintf("invalid argument :: %s - %s", argument, reason)),
	}
}

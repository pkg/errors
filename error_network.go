package errors

import (
	"fmt"
)

// ErrorNetwork is returned when a network operation fails.
type ErrorNetwork struct {
	BaseError
}

// Network returns a new instance of ErrorNetwork.
func Network(format string, a ...interface{}) error {
	return ErrorNetwork{
		NewBase("network operation failed :: " + fmt.Sprintf(format, a...)),
	}
}

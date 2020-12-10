package errors

import (
	"fmt"
)

// ErrorConflictingChange is returned when a change to a resource would conflict with an existing resource.
type ErrorConflictingChange struct {
	BaseError
}

// ConflictingChange returns a new instance of ErrorConflictingChange.
func ConflictingChange(format string, a ...interface{}) error {
	return ErrorConflictingChange{
		NewBase("conflicting change :: " + fmt.Sprintf(format, a...)),
	}
}

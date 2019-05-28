package errors

import "fmt"

// ErrorUnsupportedType is returned when trying to process an unknown type
type ErrorUnsupportedType struct {
	BaseError
}

// UnsupportedType returns a new instance of ErrorUnsupportedType
func UnsupportedType(message string) error {
	return ErrorUnsupportedType{
		new(fmt.Sprintf("unsupported type :: %s", message)),
	}
}

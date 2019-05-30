package errors

import (
	"fmt"
)

// ErrorDataLoad is returned when database query failed
type ErrorDataLoad struct {
	BaseError
}

// DataLoad returns a new instance of ErrorDataLoad
func DataLoad(message string) error {
	return ErrorDataLoad{
		NewBase(fmt.Sprintf("data load failed :: %s", message)),
	}
}

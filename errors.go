// Package errors implements functions to manipulate errors.
package errors

import (
	"errors"
	"fmt"
)

// New returns an error that formats as the given text.
func New(text string) error {
	return errors.New(text)
}

// Errorf returns a formatted error.
func Errorf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}

// Cause returns the underlying cause of the error, if possible.
// An error value has a cause if it implements the following
// method:
//
// Cause() error
//
//
// If the error does not implement Cause, the original error will
// be returned. If the error is nil, nil will be returned without further
// investigation.
func Cause(err error) error {
	if err == nil {
		return nil
	}
	type causer interface {
		Cause() error
	}
	if err, ok := err.(causer); ok {
		return err.Cause()
	}
	return err
}

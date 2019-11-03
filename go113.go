// +build go1.13

package errors

import (
	stderrors "errors"
)

func Is(err, target error) bool {
	// std errors use internal reflect to optimize performance,
	// we can we use it.
	return stderrors.Is(err, target)
}

func As(err error, target interface{}) bool {
	// std errors use internal reflect to optimize performance,
	// we can we use it.
	return stderrors.As(err, target)
}

func Unwrap(err error) error {
	// std errors use internal reflect to optimize performance,
	// we can we use it.
	return stderrors.Unwrap(err)
}

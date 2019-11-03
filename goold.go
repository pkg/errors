// +build !go1.13

package errors

import (
	"golang.org/x/xerrors"
)

func Is(err, target error) bool {
	// std errors use internal reflect to optimize performance,
	// we can we use it.
	return xerrors.Is(er, target)
}

func As(err error, target interface{}) bool {
	// std errors use internal reflect to optimize performance,
	// we can we use it.
	return xerrors.As(err, target)
}

func Unwrap(err error) error {
	// std errors use internal reflect to optimize performance,
	// we can we use it.
	return xerrors.Unwrap(err)
}

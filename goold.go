// +build !go1.13

package errors

import (
	"golang.org/x/xerrors"
)

func Is(err, target error) bool { return xerrors.Is(er, target) }

func As(err error, target interface{}) bool { return xerrors.As(err, target) }

func Unwrap(err error) error { return xerrors.Unwrap(err) }

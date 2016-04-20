// Package errors implements functions to manipulate errors.
package errors

import (
	"fmt"
	"runtime"
)

// New returns an error that formats as the given text.
func New(message string) error {
	pc, _, _, _ := runtime.Caller(1) // the caller of New
	return struct {
		error
		pc uintptr
	}{
		fmt.Errorf(message),
		pc,
	}
}

// Errorf returns a formatted error.
func Errorf(format string, args ...interface{}) error {
	pc, _, _, _ := runtime.Caller(1) // the caller of Errorf
	return struct {
		error
		pc uintptr
	}{
		fmt.Errorf(format, args...),
		pc,
	}
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

func underlying(err error) (error, bool) {
	if err == nil {
		return nil, false
	}
	type underlying interface {
		underlying() error
	}
	if err, ok := err.(underlying); ok {
		return err.underlying(), true
	}
	return nil, false
}

type traced struct {
	error // underlying error
	pc    uintptr
}

func (t *traced) underlying() error { return t.error }

// Trace adds caller information to the error.
// If error is nil, nil will be returned.
func Trace(err error) error {
	if err == nil {
		return nil
	}
	pc, _, _, _ := runtime.Caller(1) // the caller of Trace
	return traced{
		error: err,
		pc:    pc,
	}
}

type annotated struct {
	error // underlying error
	pc uintptr
}

func (a *annotated) Cause() error { return a.error }

// Annotate returns a new error annotating the error provided
// with the message, and the location of the caller of Annotate.
// The underlying error can be recovered by calling Cause.
// If err is nil, nil will be returned.
func Annotate(err error, message string) error {
	if err == nil {
		return nil
	}
	pc, _, _, _ := runtime.Caller(1) // the caller of Annotate
	return annotated{
		error: err,
		pc: pc,
	}
}

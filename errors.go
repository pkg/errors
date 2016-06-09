// Package errors provides simple error handling primitives.
//
// The traditional error handling idiom in Go is roughly akin to
//
//      if err != nil {
//              return err
//      }
//
// which applied recursively up the call stack results in error reports
// without context or debugging information. The errors package allows
// programmers to add context to the failure path in their code in a way
// that does not destroy the original value of the error.
//
// Adding context to an error
//
// The errors.Wrap function returns a new error that adds context to the
// original error. For example
//
//      _, err := ioutil.ReadAll(r)
//      if err != nil {
//              return errors.Wrap(err, "read failed")
//      }
//
// Retrieving the cause of an error
//
// Using errors.Wrap constructs a stack of errors, adding context to the
// preceding error. Depending on the nature of the error it may be necessary
// to reverse the operation of errors.Wrap to retrieve the original error
// for inspection. Any error value which implements this interface
//
//     type Causer interface {
//             Cause() error
//     }
//
// can be inspected by errors.Cause. errors.Cause will recursively retrieve
// the topmost error which does not implement causer, which is assumed to be
// the original cause. For example:
//
//     switch err := errors.Cause(err).(type) {
//     case *MyError:
//             // handle specifically
//     default:
//             // unknown error
//     }
//
// Retrieving the stack trace of an error or wrapper
//
// New, Errorf, Wrap, and Wrapf record a stack trace at the point they are
// invoked. This information can be retrieved with the following interface.
//
//     type Stacktrace interface {
//             Stacktrace() []Frame
//     }
package errors

import (
	"fmt"
	"io"
)

// _error is an error implementation returned by New and Errorf
// that implements its own fmt.Formatter.
type _error struct {
	msg string
	*stack
}

func (e _error) Error() string { return e.msg }

func (e _error) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v: ", e.Stacktrace()[0])
		}
		fallthrough
	case 's':
		io.WriteString(s, e.msg)
	}
}

// New returns an error that formats as the given text.
func New(text string) error {
	return _error{
		text,
		callers(),
	}
}

// Errorf formats according to a format specifier and returns the string
// as a value that satisfies error.
func Errorf(format string, args ...interface{}) error {
	return _error{
		fmt.Sprintf(format, args...),
		callers(),
	}
}

type cause struct {
	cause error
	msg   string
}

func (c cause) Error() string { return fmt.Sprintf("%s: %v", c.msg, c.Cause()) }
func (c cause) Cause() error  { return c.cause }

// wrapper is an error implementation returned by Wrap and Wrapf
// that implements its own fmt.Formatter.
type wrapper struct {
	cause
	*stack
}

func (w wrapper) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v\n", w.Cause())
			fmt.Fprintf(s, "%+v: %s", w.Stacktrace()[0], w.msg)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, w.Error())
	}
}

// Wrap returns an error annotating err with message.
// If err is nil, Wrap returns nil.
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return wrapper{
		cause: cause{
			cause: err,
			msg:   message,
		},
		stack: callers(),
	}
}

// Wrapf returns an error annotating err with the format specifier.
// If err is nil, Wrapf returns nil.
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return wrapper{
		cause: cause{
			cause: err,
			msg:   fmt.Sprintf(format, args...),
		},
		stack: callers(),
	}
}

// Cause returns the underlying cause of the error, if possible.
// An error value has a cause if it implements the following
// interface:
//
//     type Causer interface {
//            Cause() error
//     }
//
// If the error does not implement Cause, the original error will
// be returned. If the error is nil, nil will be returned without further
// investigation.
func Cause(err error) error {
	type causer interface {
		Cause() error
	}

	for err != nil {
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return err
}

// Fprint prints the error to the supplied writer.
// If the error implements the Causer interface described in Cause
// Fprint will recurse into the error's cause.
// Fprint will also print the file and line of the error.
// If err is nil, nothing is printed.
//
// Deprecated: Fprint will be removed in version 0.7.
func Fprint(w io.Writer, err error) {
	if err == nil {
		return
	}
	fmt.Fprintf(w, "%+v\n", err)
}

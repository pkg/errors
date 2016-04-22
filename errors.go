// Package errors implements functions to manipulate errors.
package errors

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
)

type loc uintptr

func (l loc) Location() (string, int) {
	pc := uintptr(l)
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "unknown", 0
	}

	_, prefix, _, _ := runtime.Caller(0)
	file, line := fn.FileLine(pc)
	if i := strings.LastIndex(prefix, "github.com/pkg/errors"); i > 0 {
		file = file[i:]
	}
	return file, line
}

// New returns an error that formats as the given text.
func New(text string) error {
	return struct {
		error
		loc
	}{
		fmt.Errorf(text),
		loc(pc()),
	}
}

type e struct {
	cause   error
	message string
	loc
}

func (e *e) Error() string {
	return e.message + ": " + e.cause.Error()
}

func (e *e) Cause() error {
	return e.cause
}

// Wrap returns an error annotating the cause with message.
// If cause is nil, Wrap returns nil.
func Wrap(cause error, message string) error {
	if cause == nil {
		return nil
	}
	return &e{
		cause:   cause,
		message: message,
		loc:     loc(pc()),
	}
}

type causer interface {
	Cause() error
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
	for err != nil {
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return err
}

type locationer interface {
	Location() (string, int)
}

// Print prints the error to Stderr.
// If the error implements the Causer interface described in Cause
// Print will recurse into the error's cause.
// If the error implements the inteface:
//
//     type Location interface {
//            Location() (file string, line int)
//     }
//
// Print will also print the file and line of the error.
func Print(err error) {
	Fprint(os.Stderr, err)
}

// Fprint prints the error to the supplied writer.
// The format of the output is the same as Print.
// If err is nil, nothing is printed.
func Fprint(w io.Writer, err error) {
	for err != nil {
		location, ok := err.(locationer)
		if ok {
			file, line := location.Location()
			fmt.Fprintf(w, "%s:%d: ", file, line)
		}
		switch err := err.(type) {
		case *e:
			fmt.Fprintln(w, err.message)
		default:
			fmt.Fprintln(w, err.Error())
		}

		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
}

func pc() uintptr {
	pc, _, _, _ := runtime.Caller(2)
	return pc
}

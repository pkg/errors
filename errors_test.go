package errors

import (
	"fmt"
	"io"
	"reflect"
	"testing"
)

func TestNewError(t *testing.T) {
	tests := []struct {
		err  string
		want error
	}{
		{"", fmt.Errorf("")},
		{"foo", fmt.Errorf("foo")},
		{"foo", New("foo")},
	}

	for _, tt := range tests {
		got := New(tt.err)
		if got.Error() != tt.want.Error() {
			t.Errorf("New.Error(): got: %q, want %q", got, tt.want)
		}
	}
}

func TestNewEqualNew(t *testing.T) {
	// test that two calls to New return the same error when called from the same location
	var errs []error
	for i := 0; i < 2; i++ {
		errs = append(errs, New("error"))
	}
	a, b := errs[0], errs[1]
	if !reflect.DeepEqual(a, b) {
		t.Errorf("Expected two calls to New from the same location to give the same error: %#v, %#v", a, b)
	}
}

func TestNewNotEqualNew(t *testing.T) {
	// test that two calls to New return different errors when called from different locations
	a, b := New("error"), New("error")
	if reflect.DeepEqual(a, b) {
		t.Errorf("Expected two calls to New from the different locations give the same error: %#v, %#v", a, b)
	}
}

type nilError struct{}

func (nilError) Error() string { return "nil error" }

type causeError struct {
	cause error
}

func (e *causeError) Error() string { return "cause error" }
func (e *causeError) Cause() error  { return e.cause }

func TestCause(t *testing.T) {
	x := New("error")
	tests := []struct {
		err  error
		want error
	}{{
		// nil error is nil
		err:  nil,
		want: nil,
	}, {
		// explicit nil error is nil
		err:  (error)(nil),
		want: nil,
	}, {
		// typed nil is nil
		err:  (*nilError)(nil),
		want: (*nilError)(nil),
	}, {
		// uncaused error is unaffected
		err:  io.EOF,
		want: io.EOF,
	}, {
		// caused error returns cause
		err:  &causeError{cause: io.EOF},
		want: io.EOF,
	}, {
		err:  x, // return from errors.New
		want: x,
	}}

	for i, tt := range tests {
		got := Cause(tt.err)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("test %d: got %#v, want %#v", i+1, got, tt.want)
		}
	}
}

func TestTraceNotEqual(t *testing.T) {
	// test that two calls to trace do not return identical errors
	err := New("error")
	a := err
	var errs []error
	for i := 0; i < 2; i++ {
		err = Trace(err)
		errs = append(errs, err)
	}
	b, c := errs[0], errs[1]
	if reflect.DeepEqual(a, b) {
		t.Errorf("a and b equal: %#v, %#v", a, b)
	}
	if reflect.DeepEqual(b, c) {
		t.Errorf("b and c equal: %#v, %#v", b, c)
	}
	if reflect.DeepEqual(a, c) {
		t.Errorf("a and c equal: %#v, %#v", a, c)
	}
}

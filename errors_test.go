package errors

import (
	"fmt"
	"io"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	got := New("test error")
	want := fmt.Errorf("test error")

	if !reflect.DeepEqual(got, want) {
		t.Errorf("New: got %#v, want %#v", got, want)
	}
}

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

type nilError struct{}

func (nilError) Error() string { return "nil error" }

type causeError struct {
	cause error
}

func (e *causeError) Error() string { return "cause error" }
func (e *causeError) Cause() error  { return e.cause }

func TestCause(t *testing.T) {
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
	}}

	for i, tt := range tests {
		got := Cause(tt.err)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("test %d: got %#v, want %#v", i+1, got, tt.want)
		}
	}
}

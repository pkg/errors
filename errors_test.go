package errors

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		err  string
		want error
	}{
		{"", fmt.Errorf("")},
		{"foo", fmt.Errorf("foo")},
		{"foo", New("foo")},
		{"string with format specifiers: %v", errors.New("string with format specifiers: %v")},
	}

	for _, tt := range tests {
		got := New(tt.err)
		if got.Error() != tt.want.Error() {
			t.Errorf("New.Error(): got: %q, want %q", got, tt.want)
		}
	}
}

func TestWrapNil(t *testing.T) {
	got := Wrap(nil, "no error")
	if got != nil {
		t.Errorf("Wrap(nil, \"no error\"): got %#v, expected nil", got)
	}
}

func TestWrap(t *testing.T) {
	tests := []struct {
		err     error
		message string
		want    string
	}{
		{io.EOF, "read error", "read error: EOF"},
		{Wrap(io.EOF, "read error"), "client error", "client error: read error: EOF"},
	}

	for _, tt := range tests {
		got := Wrap(tt.err, tt.message).Error()
		if got != tt.want {
			t.Errorf("Wrap(%v, %q): got: %v, want %v", tt.err, tt.message, got, tt.want)
		}
	}
}

type nilError struct{}

func (nilError) Error() string { return "nil error" }

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
		err:  Wrap(io.EOF, "ignored"),
		want: io.EOF,
	}, {
		err:  x, // return from errors.New
		want: x,
	}, {
		WithMessage(nil, "whoops"),
		nil,
	}, {
		WithMessage(io.EOF, "whoops"),
		io.EOF,
	}, {
		WithStack(nil),
		nil,
	}, {
		WithStack(io.EOF),
		io.EOF,
	}}

	for i, tt := range tests {
		got := Cause(tt.err)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("test %d: got %#v, want %#v", i+1, got, tt.want)
		}
	}
}

func TestWrapfNil(t *testing.T) {
	got := Wrapf(nil, "no error")
	if got != nil {
		t.Errorf("Wrapf(nil, \"no error\"): got %#v, expected nil", got)
	}
}

func TestWrapf(t *testing.T) {
	tests := []struct {
		err     error
		message string
		want    string
	}{
		{io.EOF, "read error", "read error: EOF"},
		{Wrapf(io.EOF, "read error without format specifiers"), "client error", "client error: read error without format specifiers: EOF"},
		{Wrapf(io.EOF, "read error with %d format specifier", 1), "client error", "client error: read error with 1 format specifier: EOF"},
	}

	for _, tt := range tests {
		got := Wrapf(tt.err, tt.message).Error()
		if got != tt.want {
			t.Errorf("Wrapf(%v, %q): got: %v, want %v", tt.err, tt.message, got, tt.want)
		}
	}
}

func TestErrorf(t *testing.T) {
	tests := []struct {
		err  error
		want string
	}{
		{Errorf("read error without format specifiers"), "read error without format specifiers"},
		{Errorf("read error with %d format specifier", 1), "read error with 1 format specifier"},
		{Errorf("read error: %w", io.EOF), "read error: EOF"},
		{Errorf("%w: client error", io.EOF), "EOF: client error"},
		{Errorf("read %d bytes error: %w", 42, io.EOF), "read 42 bytes error: EOF"},
		{Errorf("%w", io.EOF), "EOF"},
	}

	for _, tt := range tests {
		got := tt.err.Error()
		if got != tt.want {
			t.Errorf("Errorf(%v): got: %q, want %q", tt.err, got, tt.want)
		}
	}
}

func TestErrorfWrapping(t *testing.T) {
	base := errors.New("base error")
	wrapped := Errorf("wrapped: %w", base)

	// Verify Unwrap() chain reaches base via stdlib errors.Is
	u, ok := wrapped.(interface{ Unwrap() error })
	if !ok {
		t.Fatal("expected wrapped error to implement Unwrap() error")
	}
	if !errors.Is(wrapped, base) {
		t.Errorf("errors.Is(wrapped, base) = false, want true")
	}
	if got := errors.Unwrap(u.Unwrap()); got != base {
		t.Errorf("Unwrap(Unwrap(wrapped)) = %v, want %v", got, base)
	}

	// Verify nested Unwrap chains
	outer := Errorf("outer: %w", wrapped)
	if !errors.Is(outer, wrapped) {
		t.Errorf("errors.Is(outer, wrapped) = false, want true")
	}
	if !errors.Is(outer, base) {
		t.Errorf("errors.Is(outer, base) = false, want true")
	}

	// Verify non-%w returns *fundamental without Unwrap()
	plain := Errorf("plain %s", "error")
	if _, ok := plain.(*fundamental); !ok {
		t.Errorf("plain type = %T, want *fundamental", plain)
	}
	if _, ok := plain.(interface{ Unwrap() error }); ok {
		t.Errorf("Errorf without %%w should not implement Unwrap() error")
	}

	// Verify %w with nil error operand returns *fundamental without Unwrap()
	var nilErr error
	wrapNil := Errorf("wrap nil: %w", nilErr)
	if _, ok := wrapNil.(*fundamental); !ok {
		t.Errorf("wrapNil type = %T, want *fundamental", wrapNil)
	}
	if _, ok := wrapNil.(interface{ Unwrap() error }); ok {
		t.Errorf("wrapNil should not implement Unwrap() error")
	}
	if got := Cause(wrapNil); got != wrapNil {
		t.Errorf("Cause(wrapNil) = %v, want %v", got, wrapNil)
	}

	// Verify StackTrace
	st, ok := wrapped.(interface{ StackTrace() StackTrace })
	if !ok {
		t.Fatal("expected wrapped error to implement StackTrace()")
	}
	if len(st.StackTrace()) == 0 {
		t.Error("expected non-empty stack trace")
	}

	// Verify formatting verbs
	if got := fmt.Sprintf("%s", wrapped); got != "wrapped: base error" {
		t.Errorf("%%s format = %q, want %q", got, "wrapped: base error")
	}
	if got := fmt.Sprintf("%v", wrapped); got != "wrapped: base error" {
		t.Errorf("%%v format = %q, want %q", got, "wrapped: base error")
	}
	if got := fmt.Sprintf("%q", wrapped); got != `"wrapped: base error"` {
		t.Errorf("%%q format = %q, want %q", got, `"wrapped: base error"`)
	}
	plusV := fmt.Sprintf("%+v", wrapped)
	if !strings.Contains(plusV, "wrapped: base error") || !strings.Contains(plusV, "TestErrorfWrapping") {
		t.Errorf("%%+v format missing expected message or stack: %q", plusV)
	}
}

func TestErrorfMultiWrap(t *testing.T) {
	err1 := errors.New("err1")
	err2 := errors.New("err2")
	multi := Errorf("first: %w, second: %w", err1, err2)

	if !errors.Is(multi, err1) {
		t.Errorf("errors.Is(multi, err1) = false, want true")
	}
	if !errors.Is(multi, err2) {
		t.Errorf("errors.Is(multi, err2) = false, want true")
	}
	if got := multi.Error(); got != "first: err1, second: err2" {
		t.Errorf("multi.Error() = %q, want %q", got, "first: err1, second: err2")
	}
	if st, ok := multi.(interface{ StackTrace() StackTrace }); !ok || len(st.StackTrace()) == 0 {
		t.Errorf("expected multi-wrapped error to have non-empty StackTrace")
	}
	if got := fmt.Sprintf("%q", multi); got != `"first: err1, second: err2"` {
		t.Errorf("%%q format = %q, want %q", got, `"first: err1, second: err2"`)
	}
	plusV := fmt.Sprintf("%+v", multi)
	if !strings.Contains(plusV, "first: err1, second: err2") || !strings.Contains(plusV, "TestErrorfMultiWrap") {
		t.Errorf("%%+v format missing expected message or stack: %q", plusV)
	}

	// Verify all-nil %w operands return *fundamental without Unwrap()
	var nilErr error
	zeroFromMulti := Errorf("first: %w, second: %w", nilErr, nilErr)
	if _, ok := zeroFromMulti.(*fundamental); !ok {
		t.Errorf("zeroFromMulti type = %T, want *fundamental", zeroFromMulti)
	}
	if _, ok := zeroFromMulti.(interface{ Unwrap() error }); ok {
		t.Errorf("zeroFromMulti should not implement Unwrap() error")
	}
}

func TestWithStackNil(t *testing.T) {
	got := WithStack(nil)
	if got != nil {
		t.Errorf("WithStack(nil): got %#v, expected nil", got)
	}
}

func TestWithStack(t *testing.T) {
	tests := []struct {
		err  error
		want string
	}{
		{io.EOF, "EOF"},
		{WithStack(io.EOF), "EOF"},
	}

	for _, tt := range tests {
		got := WithStack(tt.err).Error()
		if got != tt.want {
			t.Errorf("WithStack(%v): got: %v, want %v", tt.err, got, tt.want)
		}
	}
}

func TestWithMessageNil(t *testing.T) {
	got := WithMessage(nil, "no error")
	if got != nil {
		t.Errorf("WithMessage(nil, \"no error\"): got %#v, expected nil", got)
	}
}

func TestWithMessage(t *testing.T) {
	tests := []struct {
		err     error
		message string
		want    string
	}{
		{io.EOF, "read error", "read error: EOF"},
		{WithMessage(io.EOF, "read error"), "client error", "client error: read error: EOF"},
	}

	for _, tt := range tests {
		got := WithMessage(tt.err, tt.message).Error()
		if got != tt.want {
			t.Errorf("WithMessage(%v, %q): got: %q, want %q", tt.err, tt.message, got, tt.want)
		}
	}
}

func TestWithMessagefNil(t *testing.T) {
	got := WithMessagef(nil, "no error")
	if got != nil {
		t.Errorf("WithMessage(nil, \"no error\"): got %#v, expected nil", got)
	}
}

func TestWithMessagef(t *testing.T) {
	tests := []struct {
		err     error
		message string
		want    string
	}{
		{io.EOF, "read error", "read error: EOF"},
		{WithMessagef(io.EOF, "read error without format specifier"), "client error", "client error: read error without format specifier: EOF"},
		{WithMessagef(io.EOF, "read error with %d format specifier", 1), "client error", "client error: read error with 1 format specifier: EOF"},
	}

	for _, tt := range tests {
		got := WithMessagef(tt.err, tt.message).Error()
		if got != tt.want {
			t.Errorf("WithMessage(%v, %q): got: %q, want %q", tt.err, tt.message, got, tt.want)
		}
	}
}

// errors.New, etc values are not expected to be compared by value
// but the change in errors#27 made them incomparable. Assert that
// various kinds of errors have a functional equality operator, even
// if the result of that equality is always false.
func TestErrorEquality(t *testing.T) {
	vals := []error{
		nil,
		io.EOF,
		errors.New("EOF"),
		New("EOF"),
		Errorf("EOF"),
		Wrap(io.EOF, "EOF"),
		Wrapf(io.EOF, "EOF%d", 2),
		WithMessage(nil, "whoops"),
		WithMessage(io.EOF, "whoops"),
		WithStack(io.EOF),
		WithStack(nil),
	}

	for i := range vals {
		for j := range vals {
			_ = vals[i] == vals[j] // mustn't panic
		}
	}
}

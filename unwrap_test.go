package errors

import (
	stderrors "errors"
	"fmt"
	"io"
	"testing"
)

func TestErrorChainCompat(t *testing.T) {
	err := stderrors.New("error that gets wrapped")
	wrapped := Wrap(err, "wrapped up")
	if !stderrors.Is(wrapped, err) {
		t.Errorf("Wrap does not support Go 1.13 error chains")
	}
}

func TestIs(t *testing.T) {
	err := New("test")

	type args struct {
		err    error
		target error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "with stack",
			args: args{
				err:    WithStack(err),
				target: err,
			},
			want: true,
		},
		{
			name: "with message",
			args: args{
				err:    WithMessage(err, "test"),
				target: err,
			},
			want: true,
		},
		{
			name: "with message format",
			args: args{
				err:    WithMessagef(err, "%s", "test"),
				target: err,
			},
			want: true,
		},
		{
			name: "std errors compatibility",
			args: args{
				err:    fmt.Errorf("wrap it: %w", err),
				target: err,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Is(tt.args.err, tt.args.target); got != tt.want {
				t.Errorf("Is() = %v, want %v", got, tt.want)
			}
		})
	}
}

type customErr struct {
	msg string
}

func (c customErr) Error() string { return c.msg }

func TestAs(t *testing.T) {
	var err = customErr{msg: "test message"}

	type args struct {
		err    error
		target any
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "with stack",
			args: args{
				err:    WithStack(err),
				target: new(customErr),
			},
			want: true,
		},
		{
			name: "with message",
			args: args{
				err:    WithMessage(err, "test"),
				target: new(customErr),
			},
			want: true,
		},
		{
			name: "with message format",
			args: args{
				err:    WithMessagef(err, "%s", "test"),
				target: new(customErr),
			},
			want: true,
		},
		{
			name: "std errors compatibility",
			args: args{
				err:    fmt.Errorf("wrap it: %w", err),
				target: new(customErr),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := As(tt.args.err, tt.args.target); got != tt.want {
				t.Errorf("As() = %v, want %v", got, tt.want)
			}

			ce := tt.args.target.(*customErr)
			if err != *ce {
				t.Errorf("set target error failed, target error is %v", *ce)
			}
		})
	}
}

func TestUnwrap(t *testing.T) {
	err := New("test")

	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "with stack",
			args: args{err: WithStack(err)},
			want: err,
		},
		{
			name: "with message",
			args: args{err: WithMessage(err, "test")},
			want: err,
		},
		{
			name: "with message format",
			args: args{err: WithMessagef(err, "%s", "test")},
			want: err,
		},
		{
			name: "std errors compatibility",
			args: args{err: fmt.Errorf("wrap: %w", err)},
			want: err,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Unwrap(tt.args.err); err != tt.want {
				t.Errorf("Unwrap() error = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestJoin(t *testing.T) {
	err1 := New("err1")
	err2 := New("err2")

	tests := []struct {
		name string
		errs []error
		want string
	}{
		{
			name: "two errors",
			errs: []error{err1, err2},
			want: "err1\nerr2",
		},
		{
			name: "nil filtered",
			errs: []error{err1, nil, err2},
			want: "err1\nerr2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Join(tt.errs...)
			if err == nil {
				t.Fatal("Join() = nil, want non-nil")
			}
			if got := err.Error(); got != tt.want {
				t.Errorf("Join().Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestJoinNil(t *testing.T) {
	if err := Join(); err != nil {
		t.Errorf("Join() = %v, want nil", err)
	}
	if err := Join(nil, nil); err != nil {
		t.Errorf("Join(nil, nil) = %v, want nil", err)
	}
}

func TestWrapAsType(t *testing.T) {
	err := customErr{msg: "test"}
	wrapped := Wrap(err, "wrapped")

	tests := []struct {
		name string
		fn   func(error) (customErr, bool)
	}{
		{name: "AsType", fn: AsType[customErr]},
		{name: "stderrors.AsType", fn: stderrors.AsType[customErr]},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := tt.fn(wrapped)
			if !ok {
				t.Fatalf("%s[customErr]() = false, want true", tt.name)
			}
			if got != err {
				t.Errorf("%s[customErr]() = %v, want %v", tt.name, got, err)
			}
		})
	}
}

func TestAsTypeNotFound(t *testing.T) {
	err := io.EOF
	assertNotFound := func(name string, ok bool) {
		t.Helper()
		if ok {
			t.Errorf("%s[customErr](io.EOF) = true, want false", name)
		}
	}

	tests := []struct {
		name string
		fn   func(error) (customErr, bool)
	}{
		{name: "AsType", fn: AsType[customErr]},
		{name: "stderrors.AsType", fn: stderrors.AsType[customErr]},
	}

	for _, tt := range tests {
		_, ok := tt.fn(err)
		assertNotFound(tt.name, ok)
	}
}

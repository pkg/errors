package errors

import (
	"fmt"
	"io"
	"testing"
)

func TestFormat(t *testing.T) {
	tests := []struct {
		error
		format string
		want   string
	}{{

		New("error"),
		"%s",
		"error",
	}, {
		New("error"),
		"%v",
		"error",
	}, {
		New("error"),
		"%+v",
		"github.com/pkg/errors/format_test.go:24: error",
	}, {
		Errorf("%s", "error"),
		"%s",
		"error",
	}, {
		Errorf("%s", "error"),
		"%v",
		"error",
	}, {
		Errorf("%s", "error"),
		"%+v",
		"github.com/pkg/errors/format_test.go:36: error",
	}, {
		Wrap(New("error"), "error2"),
		"%s",
		"error2: error",
	}, {
		Wrap(New("error"), "error2"),
		"%v",
		"error2: error",
	}, {
		Wrap(New("error"), "error2"),
		"%+v",
		"github.com/pkg/errors/format_test.go:48: error\n" +
			"github.com/pkg/errors/format_test.go:48: error2",
	}, {
		Wrap(io.EOF, "error"),
		"%s",
		"error: EOF",
	}, {
		Wrapf(New("error"), "error%d", 2),
		"%s",
		"error2: error",
	}, {
		Wrap(io.EOF, "error"),
		"%v",
		"error: EOF",
	}, {
		Wrap(io.EOF, "error"),
		"%+v",
		"EOF\n" +
			"github.com/pkg/errors/format_test.go:65: error",
	}, {
		Wrapf(New("error"), "error%d", 2),
		"%v",
		"error2: error",
	}, {
		Wrapf(New("error"), "error%d", 2),
		"%+v",
		"github.com/pkg/errors/format_test.go:74: error\n" +
			"github.com/pkg/errors/format_test.go:74: error2",
	}}

	for _, tt := range tests {
		got := fmt.Sprintf(tt.format, tt.error)
		if got != tt.want {
			t.Errorf("fmt.Sprintf(%q, err): got: %q, want: %q", tt.format, got, tt.want)
		}
	}
}

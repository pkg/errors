package errors

import (
	"fmt"
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
		"github.com/pkg/errors/format_test.go:23: error",
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
		"github.com/pkg/errors/format_test.go:35: error",
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
		"github.com/pkg/errors/format_test.go:47: error2",
	}, {
		Wrapf(New("error"), "error%d", 2),
		"%s",
		"error2: error",
	}, {
		Wrapf(New("error"), "error%d", 2),
		"%v",
		"error2: error",
	}, {
		Wrapf(New("error"), "error%d", 2),
		"%+v",
		"github.com/pkg/errors/format_test.go:59: error2",
	}}

	for _, tt := range tests {
		got := fmt.Sprintf(tt.format, tt.error)
		if got != tt.want {
			t.Errorf("fmt.Sprintf(%q, err): got: %q, want: %q", tt.format, got, tt.want)
		}
	}
}

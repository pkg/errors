package errors

import (
	"fmt"
	"io"
	"strings"
	"testing"
)

func testFormat(t *testing.T, err error, format, want string) {
	got := fmt.Sprintf(format, err)
	lines := strings.SplitN(got, "\n", -1)
	for i, w := range strings.SplitN(want, "\n", -1) {
		if lines[i] != w {
			t.Errorf("fmt.Sprintf(%q, err): got: %q, want: %q", format, got, want)
		}
	}
}

func TestFormatNew(t *testing.T) {
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
		"error\n" +
			"github.com/pkg/errors.TestFormatNew\n" +
			"\t/home/dfc/src/github.com/pkg/errors/format_test.go:34",
	}}

	for _, tt := range tests {
		testFormat(t, tt.error, tt.format, tt.want)
	}
}

func TestFormatErrorf(t *testing.T) {
	tests := []struct {
		error
		format string
		want   string
	}{{
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
		"error\n" +
			"github.com/pkg/errors.TestFormatErrorf\n" +
			"\t/home/dfc/src/github.com/pkg/errors/format_test.go:60",
	}}

	for _, tt := range tests {
		testFormat(t, tt.error, tt.format, tt.want)
	}
}

func TestFormatWrap(t *testing.T) {
	tests := []struct {
		error
		format string
		want   string
	}{{
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
		"error\n" +
			"github.com/pkg/errors.TestFormatWrap\n" +
			"\t/home/dfc/src/github.com/pkg/errors/format_test.go:86",
	}, {
		Wrap(io.EOF, "error"),
		"%s",
		"error: EOF",
	}}

	for _, tt := range tests {
		testFormat(t, tt.error, tt.format, tt.want)
	}
}

func TestFormatWrapf(t *testing.T) {
	tests := []struct {
		error
		format string
		want   string
	}{{
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
			"github.com/pkg/errors.TestFormatWrapf\n" +
			"\t/home/dfc/src/github.com/pkg/errors/format_test.go:116: error",
	}, {
		Wrapf(New("error"), "error%d", 2),
		"%v",
		"error2: error",
	}, {
		Wrapf(New("error"), "error%d", 2),
		"%+v",
		"error\n" +
			"github.com/pkg/errors.TestFormatWrapf\n" +
			"\t/home/dfc/src/github.com/pkg/errors/format_test.go:126",
	}}

	for _, tt := range tests {
		testFormat(t, tt.error, tt.format, tt.want)
	}
}

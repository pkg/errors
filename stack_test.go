package errors

import (
	"fmt"
	"runtime"
	"testing"
)

var line, _, _, _ = runtime.Caller(0)

func TestFrameLine(t *testing.T) {
	var tests = []struct {
		Frame
		want int
	}{{
		Frame(line),
		9,
	}, {
		func() Frame {
			var line, _, _, _ = runtime.Caller(0)
			return Frame(line)
		}(),
		20,
	}, {
		func() Frame {
			var line, _, _, _ = runtime.Caller(1)
			return Frame(line)
		}(),
		28,
	}, {
		Frame(0), // invalid PC
		0,
	}}

	for _, tt := range tests {
		got := tt.Frame.line()
		want := tt.want
		if want != got {
			t.Errorf("Frame(%v): want: %v, got: %v", uintptr(tt.Frame), want, got)
		}
	}
}

func TestStackLocation(t *testing.T) {
	st := func() *stack {
		var pcs [32]uintptr
		n := runtime.Callers(1, pcs[:])
		var st stack = pcs[0:n]
		return &st
	}()
	file, line := st.Location()
	wfile, wline := "github.com/pkg/errors/stack_test.go", 47
	if file != wfile || line != wline {
		t.Errorf("stack.Location(): want %q %d, got %q %d", wfile, wline, file, line)
	}
}

func TestFrameFormat(t *testing.T) {
	var tests = []struct {
		Frame
		format string
		want   string
	}{{
		Frame(line),
		"%s",
		"stack_test.go",
	}, {
		Frame(line),
		"%+s",
		"github.com/pkg/errors/stack_test.go",
	}, {
		Frame(0),
		"%s",
		"unknown",
	}, {
		Frame(line),
		"%d",
		"9",
	}, {
		Frame(0),
		"%d",
		"0",
	}, {
		Frame(line),
		"%n",
		"init",
	}, {
		Frame(0),
		"%n",
		"",
	}, {
		Frame(line),
		"%v",
		"stack_test.go:9",
	}, {
		Frame(0),
		"%v",
		"unknown:0",
	}}

	for _, tt := range tests {
		got := fmt.Sprintf(tt.format, tt.Frame)
		want := tt.want
		if want != got {
			t.Errorf("%v %q: want: %q, got: %q", tt.Frame, tt.format, want, got)
		}
	}
}

func TestTrimGOPATH(t *testing.T) {
	var tests = []struct {
		Frame
		want string
	}{{
		Frame(line),
		"github.com/pkg/errors/stack_test.go",
	}}

	for _, tt := range tests {
		pc := tt.Frame.pc()
		fn := runtime.FuncForPC(pc)
		got := trimGOPATH(fn, pc)
		want := tt.want
		if want != got {
			t.Errorf("%v: want %q, got %q", tt.Frame, want, got)
		}
	}
}

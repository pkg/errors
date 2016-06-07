package errors

import (
	"fmt"
	"runtime"
	"testing"
)

var initpc, _, _, _ = runtime.Caller(0)

func TestFrameLine(t *testing.T) {
	var tests = []struct {
		Frame
		want int
	}{{
		Frame(initpc),
		9,
	}, {
		func() Frame {
			var pc, _, _, _ = runtime.Caller(0)
			return Frame(pc)
		}(),
		20,
	}, {
		func() Frame {
			var pc, _, _, _ = runtime.Caller(1)
			return Frame(pc)
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

type X struct{}

func (x X) val() Frame {
	var pc, _, _, _ = runtime.Caller(0)
	return Frame(pc)
}

func (x *X) ptr() Frame {
	var pc, _, _, _ = runtime.Caller(0)
	return Frame(pc)
}

func TestFrameFormat(t *testing.T) {
	var tests = []struct {
		Frame
		format string
		want   string
	}{{
		Frame(initpc),
		"%s",
		"stack_test.go",
	}, {
		Frame(initpc),
		"%+s",
		"github.com/pkg/errors/stack_test.go",
	}, {
		Frame(0),
		"%s",
		"unknown",
	}, {
		Frame(0),
		"%+s",
		"unknown",
	}, {
		Frame(initpc),
		"%d",
		"9",
	}, {
		Frame(0),
		"%d",
		"0",
	}, {
		Frame(initpc),
		"%n",
		"init",
	}, {
		func() Frame {
			var x X
			return x.ptr()
		}(),
		"%n",
		"(*X).ptr",
	}, {
		func() Frame {
			var x X
			return x.val()
		}(),
		"%n",
		"X.val",
	}, {
		Frame(0),
		"%n",
		"",
	}, {
		Frame(initpc),
		"%v",
		"stack_test.go:9",
	}, {
		Frame(initpc),
		"%+v",
		"github.com/pkg/errors/stack_test.go:9",
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
		Frame(initpc),
		"github.com/pkg/errors/stack_test.go",
	}}

	for _, tt := range tests {
		pc := tt.Frame.pc()
		fn := runtime.FuncForPC(pc)
		file, _ := fn.FileLine(pc)
		got := trimGOPATH(fn.Name(), file)
		want := tt.want
		if want != got {
			t.Errorf("%v: want %q, got %q", tt.Frame, want, got)
		}
	}
}

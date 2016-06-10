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

func TestFuncname(t *testing.T) {
	tests := []struct {
		name, want string
	}{
		{"", ""},
		{"runtime.main", "main"},
		{"github.com/pkg/errors.funcname", "funcname"},
		{"funcname", "funcname"},
		{"io.copyBuffer", "copyBuffer"},
		{"main.(*R).Write", "(*R).Write"},
	}

	for _, tt := range tests {
		got := funcname(tt.name)
		want := tt.want
		if got != want {
			t.Errorf("funcname(%q): want: %q, got %q", tt.name, want, got)
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

func TestStacktrace(t *testing.T) {
	type fileline struct {
		file string
		line int
	}
	tests := []struct {
		err  error
		want []fileline
	}{{
		New("ooh"), []fileline{
			{"github.com/pkg/errors/stack_test.go", 181},
		},
	}, {
		Wrap(New("ooh"), "ahh"), []fileline{
			{"github.com/pkg/errors/stack_test.go", 185}, // this is the stack of Wrap, not New
		},
	}, {
		Cause(Wrap(New("ooh"), "ahh")), []fileline{
			{"github.com/pkg/errors/stack_test.go", 189}, // this is the stack of New
		},
	}, {
		func() error { return New("ooh") }(), []fileline{
			{"github.com/pkg/errors/stack_test.go", 193}, // this is the stack of New
			{"github.com/pkg/errors/stack_test.go", 193}, // this is the stack of New's caller
		},
	}, {
		Cause(func() error {
			return func() error {
				return Errorf("hello %s", fmt.Sprintf("world"))
			}()
		}()), []fileline{
			{"github.com/pkg/errors/stack_test.go", 200}, // this is the stack of Errorf
			{"github.com/pkg/errors/stack_test.go", 201}, // this is the stack of Errorf's caller
			{"github.com/pkg/errors/stack_test.go", 202}, // this is the stack of Errorf's caller's caller
		},
	}}
	for _, tt := range tests {
		x, ok := tt.err.(interface {
			Stacktrace() Stacktrace
		})
		if !ok {
			t.Errorf("expected %#v to implement Stacktrace() Stacktrace", tt.err)
			continue
		}
		st := x.Stacktrace()
		for i, want := range tt.want {
			frame := st[i]
			file, line := fmt.Sprintf("%+s", frame), frame.line()
			if file != want.file || line != want.line {
				t.Errorf("frame %d: expected %s:%d, got %s:%d", i, want.file, want.line, file, line)
			}
		}
	}
}

func stacktrace() Stacktrace {
	const depth = 8
	var pcs [depth]uintptr
	n := runtime.Callers(1, pcs[:])
	var st stack = pcs[0:n]
	return st.Stacktrace()
}

func TestStacktraceFormat(t *testing.T) {
	tests := []struct {
		Stacktrace
		format string
		want   string
	}{{
		nil,
		"%s",
		"[]",
	}, {
		nil,
		"%v",
		"[]",
	}, {
		nil,
		"%+v",
		"[]",
	}, {
		make(Stacktrace, 0),
		"%s",
		"[]",
	}, {
		make(Stacktrace, 0),
		"%v",
		"[]",
	}, {
		make(Stacktrace, 0),
		"%+v",
		"[]",
	}, {
		stacktrace()[:2],
		"%s",
		"[stack_test.go stack_test.go]",
	}, {
		stacktrace()[:2],
		"%v",
		"[stack_test.go:230 stack_test.go:269]",
	}, {
		stacktrace()[:2],
		"%+v",
		"[github.com/pkg/errors/stack_test.go:230 github.com/pkg/errors/stack_test.go:273]",
	}}

	for i, tt := range tests {
		got := fmt.Sprintf(tt.format, tt.Stacktrace)
		if got != tt.want {
			t.Errorf("test %d: got: %q, want: %q", i+1, got, tt.want)
		}
	}
}

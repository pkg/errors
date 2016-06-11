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
		"github.com/pkg/errors.init\n" +
			"\t.+/github.com/pkg/errors/stack_test.go",
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
		`\(\*X\).ptr`,
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
		"github.com/pkg/errors.init\n" +
			"\t.+/github.com/pkg/errors/stack_test.go:9",
	}, {
		Frame(0),
		"%v",
		"unknown:0",
	}}

	for _, tt := range tests {
		testFormatRegexp(t, tt.Frame, tt.format, tt.want)
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
	tests := []struct {
		err  error
		want []string
	}{{
		New("ooh"), []string{
			"github.com/pkg/errors.TestStacktrace\n" +
				"\t.+/github.com/pkg/errors/stack_test.go:175",
		},
	}, {
		Wrap(New("ooh"), "ahh"), []string{
			"github.com/pkg/errors.TestStacktrace\n" +
				"\t.+/github.com/pkg/errors/stack_test.go:180", // this is the stack of Wrap, not New
		},
	}, {
		Cause(Wrap(New("ooh"), "ahh")), []string{
			"github.com/pkg/errors.TestStacktrace\n" +
				"\t.+/github.com/pkg/errors/stack_test.go:185", // this is the stack of New
		},
	}, {
		func() error { return New("ooh") }(), []string{
			`github.com/pkg/errors.(func·005|TestStacktrace.func1)` +
				"\n\t.+/github.com/pkg/errors/stack_test.go:190", // this is the stack of New
			"github.com/pkg/errors.TestStacktrace\n" +
				"\t.+/github.com/pkg/errors/stack_test.go:190", // this is the stack of New's caller
		},
	}, {
		Cause(func() error {
			return func() error {
				return Errorf("hello %s", fmt.Sprintf("world"))
			}()
		}()), []string{
			`github.com/pkg/errors.(func·006|TestStacktrace.func2.1)` +
				"\n\t.+/github.com/pkg/errors/stack_test.go:199", // this is the stack of Errorf
			`github.com/pkg/errors.(func·007|TestStacktrace.func2)` +
				"\n\t.+/github.com/pkg/errors/stack_test.go:200", // this is the stack of Errorf's caller
			"github.com/pkg/errors.TestStacktrace\n" +
				"\t.+/github.com/pkg/errors/stack_test.go:201", // this is the stack of Errorf's caller's caller
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
		for j, want := range tt.want {
			testFormatRegexp(t, st[j], "%+v", want)
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
		`\[\]`,
	}, {
		nil,
		"%v",
		`\[\]`,
	}, {
		nil,
		"%+v",
		"",
	}, {
		nil,
		"%#v",
		`\[\]errors.Frame\(nil\)`,
	}, {
		make(Stacktrace, 0),
		"%s",
		`\[\]`,
	}, {
		make(Stacktrace, 0),
		"%v",
		`\[\]`,
	}, {
		make(Stacktrace, 0),
		"%+v",
		"",
	}, {
		make(Stacktrace, 0),
		"%#v",
		`\[\]errors.Frame{}`,
	}, {
		stacktrace()[:2],
		"%s",
		`\[stack_test.go stack_test.go\]`,
	}, {
		stacktrace()[:2],
		"%v",
		`\[stack_test.go:228 stack_test.go:275\]`,
	}, {
		stacktrace()[:2],
		"%+v",
		"\n" +
			"github.com/pkg/errors.stacktrace\n" +
			"\t.+/github.com/pkg/errors/stack_test.go:228\n" +
			"github.com/pkg/errors.TestStacktraceFormat\n" +
			"\t.+/github.com/pkg/errors/stack_test.go:279",
	}, {
		stacktrace()[:2],
		"%#v",
		`\[\]errors.Frame{stack_test.go:228, stack_test.go:287}`,
	}}

	for _, tt := range tests {
		testFormatRegexp(t, tt.Stacktrace, tt.format, tt.want)
	}
}

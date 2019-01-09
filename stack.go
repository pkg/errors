package errors

import (
	"bytes"
	"fmt"
	"io"
	"path"
	"runtime"
	"strconv"
	"strings"
)

// Frame represents a program counter inside a stack frame.
type Frame uintptr

// pc returns the program counter for this frame;
// multiple frames may have the same PC value.
func (f Frame) pc() uintptr { return uintptr(f) - 1 }

// file returns the full path to the file that contains the
// function for this Frame's pc.
func (f Frame) file() string {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return "unknown"
	}
	file, _ := fn.FileLine(f.pc())
	return file
}

// line returns the line number of source code of the
// function for this Frame's pc.
func (f Frame) line() int {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return 0
	}
	_, line := fn.FileLine(f.pc())
	return line
}

// name returns the name of this function, if known.
func (f Frame) name() string {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return "unknown"
	}
	return fn.Name()
}

// Format formats the frame according to the fmt.Formatter interface.
//
//    %s    source file
//    %d    source line
//    %n    function name
//    %v    equivalent to %s:%d
//
// Format accepts flags that alter the printing of some verbs, as follows:
//
//    %+s   function name and path of source file relative to the compile time
//          GOPATH separated by \n\t (<funcname>\n\t<path>)
//    %+v   equivalent to %+s:%d
func (f Frame) Format(s fmt.State, verb rune) {
	f.format(s, s, verb)
}

// format allows stack trace printing calls to be made with a bytes.Buffer.
func (f Frame) format(w io.Writer, s fmt.State, verb rune) {
	switch verb {
	case 's':
		switch {
		case s.Flag('+'):
			io.WriteString(w, f.name())
			io.WriteString(w, "\n\t")
			io.WriteString(w, f.file())
		default:
			io.WriteString(w, path.Base(f.file()))
		}
	case 'd':
		io.WriteString(w, strconv.Itoa(f.line()))
	case 'n':
		io.WriteString(w, funcname(f.name()))
	case 'v':
		f.format(w, s, 's')
		io.WriteString(w, ":")
		f.format(w, s, 'd')
	}
}

// StackTrace is stack of Frames from innermost (newest) to outermost (oldest).
type StackTrace []Frame

// Format formats the stack of Frames according to the fmt.Formatter interface.
//
//    %s	lists source files for each Frame in the stack
//    %v	lists the source file and line number for each Frame in the stack
//
// Format accepts flags that alter the printing of some verbs, as follows:
//
//    %+v   Prints filename, function, and line number for each Frame in the stack.
func (st StackTrace) Format(s fmt.State, verb rune) {
	var b bytes.Buffer
	switch verb {
	case 'v':
		switch {
		case s.Flag('+'):
			b.Grow(len(st) * stackMinLen)
			for _, f := range st {
				b.WriteByte('\n')
				f.format(&b, s, verb)
			}
		case s.Flag('#'):
			fmt.Fprintf(&b, "%#v", []Frame(st))
		default:
			st.formatSlice(&b, s, verb)
		}
	case 's':
		st.formatSlice(&b, s, verb)
	}
	io.Copy(s, &b)
}

// formatSlice will format this StackTrace into the given buffer as a slice of
// Frame, only valid when called with '%s' or '%v'.
func (st StackTrace) formatSlice(b *bytes.Buffer, s fmt.State, verb rune) {
	b.WriteByte('[')
	if len(st) == 0 {
		b.WriteByte(']')
		return
	}

	b.Grow(len(st) * (stackMinLen / 4))
	st[0].format(b, s, verb)
	for _, fr := range st[1:] {
		b.WriteByte(' ')
		fr.format(b, s, verb)
	}
	b.WriteByte(']')
}

// stackMinLen is a best-guess at the minimum length of a stack trace. It
// doesn't need to be exact, just give a good enough head start for the buffer
// to avoid the expensive early growth.
const stackMinLen = 96

// stack represents a stack of program counters.
type stack []uintptr

func (s *stack) Format(st fmt.State, verb rune) {
	switch verb {
	case 'v':
		switch {
		case st.Flag('+'):
			for _, pc := range *s {
				f := Frame(pc)
				fmt.Fprintf(st, "\n%+v", f)
			}
		}
	}
}

func (s *stack) StackTrace() StackTrace {
	f := make([]Frame, len(*s))
	for i := 0; i < len(f); i++ {
		f[i] = Frame((*s)[i])
	}
	return f
}

func callers() *stack {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	var st stack = pcs[0:n]
	return &st
}

// funcname removes the path prefix component of a function's name reported by func.Name().
func funcname(name string) string {
	i := strings.LastIndex(name, "/")
	name = name[i+1:]
	i = strings.Index(name, ".")
	return name[i+1:]
}

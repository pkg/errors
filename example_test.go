package errors_test

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
)

func ExampleNew() {
	err := errors.New("whoops")
	fmt.Println(err)

	// Output: whoops
}

func ExampleNew_printf() {
	err := errors.New("whoops")
	fmt.Printf("%+v", err)

	// Output: github.com/pkg/errors/example_test.go:18: whoops
}

func ExampleWrap() {
	cause := errors.New("whoops")
	err := errors.Wrap(cause, "oh noes")
	fmt.Println(err)

	// Output: oh noes: whoops
}

func fn() error {
	e1 := errors.New("error")
	e2 := errors.Wrap(e1, "inner")
	e3 := errors.Wrap(e2, "middle")
	return errors.Wrap(e3, "outer")
}

func ExampleCause() {
	err := fn()
	fmt.Println(err)
	fmt.Println(errors.Cause(err))

	// Output: outer: middle: inner: error
	// error
}

func ExampleFprint() {
	err := fn()
	errors.Fprint(os.Stdout, err)

	// Output: github.com/pkg/errors/example_test.go:33: error
	// github.com/pkg/errors/example_test.go:34: inner
	// github.com/pkg/errors/example_test.go:35: middle
	// github.com/pkg/errors/example_test.go:36: outer
}

func ExampleWrapf() {
	cause := errors.New("whoops")
	err := errors.Wrapf(cause, "oh noes #%d", 2)
	fmt.Println(err)

	// Output: oh noes #2: whoops
}

func ExampleErrorf() {
	err := errors.Errorf("whoops: %s", "foo")
	fmt.Printf("%+v", err)

	// Output: github.com/pkg/errors/example_test.go:67: whoops: foo
}

func Example_stacktrace() {
	type Stacktrace interface {
		Stacktrace() []errors.Frame
	}

	err, ok := errors.Cause(fn()).(Stacktrace)
	if !ok {
		panic("oops, err does not implement Stacktrace")
	}

	st := err.Stacktrace()
	fmt.Printf("%+v", st[0:2]) // top two framces

	// Output: [github.com/pkg/errors/example_test.go:33 github.com/pkg/errors/example_test.go:78]
}

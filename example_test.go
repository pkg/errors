package errors_test

import (
	"fmt"

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

	// Output: github.com/pkg/errors/example_test.go:17: whoops
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

func ExampleCause_printf() {
	err := fn()
	fmt.Printf("%+v\n", err)

	// Output: github.com/pkg/errors/example_test.go:32: error
	// github.com/pkg/errors/example_test.go:33: inner
	// github.com/pkg/errors/example_test.go:34: middle
	// github.com/pkg/errors/example_test.go:35: outer
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

	// Output: github.com/pkg/errors/example_test.go:66: whoops: foo
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
	fmt.Printf("%+v", st[0:2]) // top two frames

	// Output: [github.com/pkg/errors/example_test.go:32 github.com/pkg/errors/example_test.go:77]
}

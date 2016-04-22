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

	// Output: github.com/pkg/errors/example_test.go:29: outer
	// github.com/pkg/errors/example_test.go:28: middle
	// github.com/pkg/errors/example_test.go:27: inner
	// error
}

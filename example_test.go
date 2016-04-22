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
	return errors.Wrap(errors.Wrap(errors.Wrap(errors.New("error"), "inner"), "middle"), "outer")
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

	// Output: outer
	// middle
	// inner
	// error
}

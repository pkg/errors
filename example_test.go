package errors_test

import (
	"errors"
	"fmt"
)

func ExampleNew() {
	err := errors.New("whoops")
	fmt.Println(err.Error())

	// Output: whoops
}

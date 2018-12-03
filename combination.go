package errors

import (
	"fmt"
	"io"
	"strings"
)

// Combination combines multiple errors into one.
// The Error method returns the strings from the individual Error methods
// joined by the new line character '\n'.
// Combination always has at least one error and returns the first error
// as result of the Cause method.
type Combination struct {
	errs []error
	*stack
}

// Combine returns a Combination error for 2 or more errors which are not nil,
// or the callstack wrapped error if only one error was passed,
// or nil if zero arguments are passed or all passed errors are nil.
// Any returned non nil error will be wrapped with a callstack,
// independent if it's a single error or a Combination type is returned.
// The Combination type's Error method returns the strings from the individual Error methods
// joined by the new line character '\n'.
func Combine(errs ...error) error {
	switch len(errs) {
	case 0:
		return nil
	case 1:
		return WithStackSkip(1, errs[0])
	}

	var lastNotNil error
	numNotNil := 0
	for _, err := range errs {
		if err != nil {
			lastNotNil = err
			numNotNil++
		}
	}
	switch numNotNil {
	case 0:
		return nil
	case 1:
		return WithStackSkip(1, lastNotNil)
	}

	if numNotNil < len(errs) {
		notNilErrs := make([]error, 0, numNotNil)
		for _, err := range errs {
			if err != nil {
				notNilErrs = append(notNilErrs, err)
			}
		}
		errs = notNilErrs
	}

	return &Combination{
		errs,
		callers(0),
	}
}

func (c *Combination) Error() string {
	if len(c.errs) == 1 {
		return c.errs[0].Error()
	}

	var b strings.Builder
	for i, err := range c.errs {
		if i > 0 {
			b.WriteByte('\n')
		}
		b.WriteString(err.Error())
	}
	return b.String()
}

func (c *Combination) Errs() []error {
	return c.errs
}

func (c *Combination) Cause() error {
	return Cause(c.errs[0])
}

func (c *Combination) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			for _, e := range c.errs {
				fmt.Fprintf(s, "%+v\n", Cause(e).Error())
			}
			c.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, c.Error())
	case 'q':
		fmt.Fprintf(s, "%q", c.Error())
	}
}

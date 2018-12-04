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
// independent if a single error or a Combination type is returned.
// The Combination type's Error method returns the strings from the
// individual Error methods joined by the new line character '\n'.
// Note that Cause(error) can only return a single error,
// so in case of a Combination type, Cause returns the first error.
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

// Uncombine returns multible errors if err is a Combination type,
// or the passed single err if was not a Combination type,
// or nil if the passed error was nil.
func Uncombine(err error) []error {
	if err == nil {
		return nil
	}
	if combo, ok := err.(*Combination); ok {
		return combo.errs
	}
	return []error{err}
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

func (c *Combination) Errors() []error {
	return c.errs
}

func (c *Combination) Append(err error) {
	if err != nil {
		c.errs = append(c.errs, WithStackSkip(1, err))
	}
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

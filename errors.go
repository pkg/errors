// Package errors implements functions to manipulate errors.
package errors

type errorString string

func (e errorString) Error() string {
	return string(e)
}

// New returns an error that formats as the given text.
func New(text string) error {
	return errorString(text)
}

// Cause returns the underlying cause of the error, if possible.
// An error value has a cause if it implements the following
// method:
//
// Cause() error
//
//
// If the error does not implement Cause, the original error will
// be returned. If the error is nil, nil will be returned without further
// investigation.
func Cause(err error) error {
	if err == nil {
		return nil
	}
	if err, ok := err.(interface {
		Cause() error
	}); ok {
		return err.Cause()
	}
	return err
}

// cause implements the interface required by Cause.
type cause struct {
	err error
}

func (c *cause) Cause() error {
	return c.err
}

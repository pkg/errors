// Package errors implements functions to manipulate errors.
package errors

type e struct {
	cause   error
	message string
}

func (e *e) Error() string {
	return e.message + ": " + e.cause.Error()
}

func (e *e) Cause() error {
	return e.cause
}

// Wrap returns an error annotating the cause with message.
// If cause is nil, Wrap returns nil.
func Wrap(cause error, message string) error {
	if cause == nil {
		return nil
	}
	return &e{cause: cause, message: message}
}

// Cause returns the underlying cause of the error, if possible.
// An error value has a cause if it implements the following
// interface:
//
//     type Causer interface {
//            Cause() error
//     }
//
// If the error does not implement Cause, the original error will
// be returned. If the error is nil, nil will be returned without further
// investigation.
func Cause(err error) error {
	if err == nil {
		return nil
	}
	type causer interface {
		Cause() error
	}
	if err, ok := err.(causer); ok {
		return err.Cause()
	}
	return err
}

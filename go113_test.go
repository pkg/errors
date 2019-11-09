// +build go1.13

package errors

import (
	stdlib_errors "errors"
	"testing"
)

func TestErrorChainCompat(t *testing.T) {
	err := stdlib_errors.New("error that gets wrapped")
	wrapped := Wrap(err, "wrapped up")
	if !stdlib_errors.Is(wrapped, err) {
		t.Errorf("Wrap does not support Go 1.13 error chains")
	}
}

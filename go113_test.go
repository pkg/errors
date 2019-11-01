// +build go1.13

package errors_test

import (
	stdlib_errors "errors"
	"testing"

	"github.com/pkg/errors"
)

func TestErrorChainCompat(t *testing.T) {
	err := stdlib_errors.New("error that gets wrapped")
	wrapped := errors.Wrap(err, "wrapped up")
	if !stdlib_errors.Is(wrapped, err) {
		t.Errorf("Wrap does not support Go 1.13 error chains")
	}
}

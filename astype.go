//go:build !go1.26

package errors

// AsType finds the first error in err's chain that matches type E,
// and if so, returns that error value and true.
func AsType[E error](err error) (E, bool) {
	var target E
	return target, As(err, &target)
}

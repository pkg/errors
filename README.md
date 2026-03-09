# errors [![GoDoc](https://godoc.org/github.com/pkg/errors?status.svg)](https://pkg.go.dev/github.com/pkg/errors) [![Report card](https://goreportcard.com/badge/github.com/pkg/errors)](https://goreportcard.com/report/github.com/pkg/errors) [![Sourcegraph](https://sourcegraph.com/github.com/pkg/errors/-/badge.svg)](https://sourcegraph.com/github.com/pkg/errors?badge)

Package errors provides simple error handling primitives.

`go get github.com/pkg/errors`

The traditional error handling idiom in Go is roughly akin to
```go
if err != nil {
        return err
}
```
which applied recursively up the call stack results in error reports without context or debugging information. The errors package allows programmers to add context to the failure path in their code in a way that does not destroy the original value of the error.

## Adding context to an error

The errors.Wrap function returns a new error that adds context to the original error. For example
```go
_, err := ioutil.ReadAll(r)
if err != nil {
        return errors.Wrap(err, "read failed")
}
```
## Retrieving the cause of an error

Using `errors.Wrap` constructs a stack of errors, adding context to the preceding error. Depending on the nature of the error it may be necessary to reverse the operation of errors.Wrap to retrieve the original error for inspection. Any error value which implements this interface can be inspected by `errors.Cause`.
```go
type causer interface {
        Cause() error
}
```
`errors.Cause` will recursively retrieve the topmost error which does not implement `causer`, which is assumed to be the original cause. For example:
```go
switch err := errors.Cause(err).(type) {
case *MyError:
        // handle specifically
default:
        // unknown error
}
```

[Read the package documentation for more information](https://godoc.org/github.com/pkg/errors).

## Compared to the standard library's `errors` package

This package was initially built to manage chains of typed errors. Support for this was later added to the standard library's `errors` package via `Unwrap` methods. Unfortunately, the standard library's `errors` package does not support stack traces, so it cannot be used as a drop-in replacement for the features this package provides. This package will mostly work interoperably with the standard library's `errors` package; typed errors added in `pkg/errors` work fine in `errors.Is` and `errors.As`, and you can wrap them with other typed errors.

## Contributing

This package is in maintenance mode. New features that extend the scope of the package are not being accepted and should be implemented in other modules. Exceptions would be features that improve interoperability with the standard library's `errors` package when it receives new features, as `pkg/errors` should mostly be able to act as a drop-in replacement. Bug fixes and reports are welcome. CI will be maintained to make sure the package is tested against new versions of Go and Go linting tools.

## License

BSD-2-Clause

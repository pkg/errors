# errors [![Travis-CI](https://travis-ci.org/pkg/errors.svg)](https://travis-ci.org/pkg/errors) [![GoDoc](https://godoc.org/github.com/pkg/errors?status.svg)](http://godoc.org/github.com/pkg/errors) [![Report card](https://goreportcard.com/badge/github.com/pkg/errors)](https://goreportcard.com/report/github.com/pkg/errors)

Package errors implements functions for manipulating errors.

The traditional error handling idiom in Go is roughly akin to
```
if err != nil {
        return err
}
```
which applied recursively up the call stack results in error reports without context or debugging information. The errors package allows programmers to add context to the failure path in their code in a way that does not destroy the original value of the error.

## Adding context to an error

The errors.Wrap function returns a new error that adds context to the original error. For example
```
_, err := ioutil.ReadAll(r)
if err != nil {
        return errors.Wrap(err, "read failed")
}
```
In addition, errors.Wrap records the file and line where it was called, allowing the programmer to retrieve the path to the original error.

## Retrieving the cause of an error

Using errors.Wrap constructs a stack of errors, adding context to the preceding error. Depending on the nature of the error it may be necessary to recurse the operation of errors.Wrap to retrieve the original error for inspection. Any error value which implements this interface
```
type causer interface {
     Cause() error
}
```
Can be inspected by errors.Cause which will recursively retrieve the topmost error which does nor implement causer, which is assumed to be the original cause. For example:
```
switch err := errors.Cause(err).(type) {
case *MyError:
        // handle specifically
default:
        // unknown error
}
```

Would you like to know more? Read the [original presentation](https://t.co/GGPr7HJZYR), blog post coming soon.

Licence: MIT

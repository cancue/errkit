# errkit

> Forked from [github.com/pkg/errors](https://github.com/pkg/errors), with modifications.

While preserving the original design philosophy, `errkit` introduces minor changes and may behave differently in specific cases—particularly in how `Cause` behaves.

## Disclaimer
- Tests are currently **insufficient** and may not cover all original test cases.
- The behavior of `Cause` is **not identical** to the original package; compatibility is not guaranteed.

## Motivation

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
_, err := os.ReadAll(r)
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

[Read the package documentation for more information](https://godoc.org/github.com/cancue/errkit).


## License

BSD-2-Clause

smquartz/errors
================

[![Build Status](https://travis-ci.org/smquartz/errors.svg?branch=master)](https://travis-ci.org/smquartz/errors)

Package errors adds stacktrace and arbitrary metadata support to errors in go.

This is particularly useful when you want to understand the state of execution
when an error was returned unexpectedly.

It provides the type \*Error which implements the standard golang error
interface, so you can use this library interchangably with code that is
expecting a normal error return.

Usage
-----

Full documentation is available on
[godoc](https://godoc.org/github.com/smquartz/errors), but here's a simple
example:

```go
package crashy

import "github.com/smquartz/errors"

var Crashed = errors.Errorf("oh dear")

func Crash() error {
    return errors.New(Crashed)
}
```

This can be called as follows:

```go
package main

import (
    "crashy"
    "fmt"
    "github.com/smquartz/errors"
)

func main() {
    err := crashy.Crash()
    if err != nil {
        if errors.Is(err, crashy.Crashed) {
            fmt.Println(err.(*errors.Error).ErrorStack())
        } else {
            panic(err)
        }
    }
}
```

Meta-fu
-------

This package is a fork of [github.com/go-errors/errors](https://github.com/go-errors/errors) that modifies
its behaviour slightly and adds a few features, including the ability
to include arbitrary metadata in your errors.

This package is licensed under the MIT license, see LICENSE.MIT for details.

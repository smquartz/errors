smquartz/errors
================

[![Build Status](https://travis-ci.org/smquartz/errors.svg?branch=master)](https://travis-ci.org/smquartz/errors)
[![Godoc Reference](https://godoc.org/github.com/smquartz/errors?status.svg)](https://godoc.org/github.com/smquartz/errors)
[![Coverage Status](https://coveralls.io/repos/github/smquartz/errors/badge.svg?branch=master)](https://coveralls.io/github/smquartz/errors?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/smquartz/errors)](https://goreportcard.com/report/github.com/smquartz/errors)
[![codebeat badge](https://codebeat.co/badges/b4e83d02-a632-4bb8-aa1d-4f3f079a319e)](https://codebeat.co/projects/github-com-smquartz-errors-master)

Package errors adds stacktrace to errors in go.

This is particularly useful when you want to understand the state of execution
when an error was returned unexpectedly.

It provides the type \*Err which implements the standard golang error
interface, so you can use this library interchangeably with code that is
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
its behaviour and adds compliance with [github.com/pkg/errors](https://github.com/pkg/errors).

This package is licensed under the MIT license, see LICENSE.MIT for details.

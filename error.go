// Package errors provides errors that have stack-traces and arbitrary metadata.
//
// This is particularly useful when you want to understand the
// state of execution when an error was returned unexpectedly.
//
// It provides the type *Error which implements the standard
// golang error interface, so you can use this library interchangably
// with code that is expecting a normal error return.
//
//
// This package is a fork of github.com/go-errors/errors that modifies
// its behaviour slightly and adds a few features, including the ability
// to include arbitrary metadata in your errors.
package errors

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

// MaxStackDepth is the maximum number of stackframes permitted on any single
// error.  This does not apply to the sum of all nested errors.
var MaxStackDepth = 50

// Error is an error with an attached stacktrace. It can be used
// wherever the builtin error interface is expected.
type Error struct {
	// underlying or "cause" error
	Err    error
	stack  []uintptr
	frames []StackFrame
	prefix string
	// arbitrary metadata that may be included in the error
	Metadata Metadata
}

// Error returns the underlying error's message.
func (err *Error) Error() string {
	var msg string

	if err.Err != nil {
		msg = err.Err.Error()
	} else {
		msg = fmt.Errorf("%v", err.Err).Error()
	}

	if err.prefix != "" {
		msg = fmt.Sprintf("%s: %s", err.prefix, msg)
	}

	return msg
}

// Stack returns the callstack formatted the same way that go does
// in runtime/debug.Stack()
func (err *Error) Stack() []byte {
	buf := bytes.Buffer{}

	for _, frame := range err.StackFrames() {
		buf.WriteString(frame.String())
	}

	return buf.Bytes()
}

// Callers satisfies the bugsnag ErrorWithCallerS() interface
// so that the stack can be read out.
func (err *Error) Callers() []uintptr {
	return err.stack
}

// ErrorStack returns a string that contains both the
// error message and the callstack.
func (err *Error) ErrorStack() string {
	return err.TypeName() + " " + err.Error() + "\n" + string(err.Stack())
}

// StackFrames returns an array of frames containing information about the
// stack.
func (err *Error) StackFrames() []StackFrame {
	if err.frames == nil {
		err.frames = make([]StackFrame, len(err.stack))

		for i, pc := range err.stack {
			err.frames[i] = NewStackFrame(pc)
		}
	}

	return err.frames
}

// StackTrace implements the pkg/errors.stacktracer interface.  It returns
// an array of frames containing information about the stack.
func (err *Error) StackTrace() errors.StackTrace {
	st := make(errors.StackTrace, len(err.stack))

	for i, pc := range err.stack {
		st[i] = errors.Frame(pc)
	}

	return st
}

// TypeName returns the type this error. e.g. *errors.stringError.
func (err *Error) TypeName() string {
	if _, ok := err.Err.(uncaughtPanic); ok {
		return "panic"
	}
	return reflect.TypeOf(err.Err).String()
}

// Cause returns the underlying cause of an error, if possible.  It returns the
// first error in the stack of nested errors, that is not of type *Error
func (err *Error) Cause() (root error) {
	root = err
	for {
		e, ok := Assert(root)
		if !ok {
			break
		}
		root = e.Err
	}
	return root
}

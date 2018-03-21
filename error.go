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
	Err   error
	stack []uintptr
	// cache of parsed stack
	frames []StackFrame
	// a prefix to prepend to the error message of the underlying error
	prefix string
	// whether to return the deepest nested stacktrace (false) or the shallowest
	// (this instance's) stacktrace
	ignoreNestedStack bool
	// arbitrary metadata that may be included in the error
	Metadata Metadata
}

// SetIgnoreNestedStack sets the ignoreNestedSTack field on the *Error this
// is called on, which determines whether functions that return information
// about the stack, return it of the stack of this *Error (true), or of the
// deepest nested *Error (false; default).
func (err *Error) SetIgnoreNestedStack(val bool) *Error {
	err.ignoreNestedStack = val
	return err
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
// in runtime/debug.Stack().  Note that this function will return
// a formatted callstack of the deepest nested *Error instance, unless
// ignoreNestedStack is set on the *Error.
func (err *Error) Stack() []byte {
	if !err.ignoreNestedStack {
		u, _ := AssertDeepestUnderlying(err)
		// we ignore the error of the above function for brevity, because an error
		// should never be returned from it with this usage, and the appropriate
		// action is to panic, hich will happen anyway if u is nil as in the case
		// of an error
		return u.ParentStack()
	}
	return err.ParentStack()
}

// ParentStack returns the callstack of the parent error, formatted the same
// way that go does in runtime/debug.Stack().  Note that this function will
// return a formatted callstack of the actual *Error this is called upon, not
// that of the deepest nested *Error.
func (err *Error) ParentStack() []byte {
	buf := bytes.Buffer{}

	for _, frame := range err.ParentStackFrames() {
		buf.WriteString(frame.String())
	}

	return buf.Bytes()
}

// ParentCallers satisfies the bugsnag ErrorWithCallerS() interface
// so that the stack can be read out.  It returns the stack of the *Error
// that this function is called on, rather than that of the deepes nested
// *Error.
func (err *Error) ParentCallers() []uintptr {
	return err.stack
}

// Callers satisfies the bugsnag ErrorWithCallers() interface so that the stack
// can be read out.  It returns the stack of the deepest nested *Error, unless
// ignoreNestedStack is set on the *Error.
func (err *Error) Callers() []uintptr {
	if !err.ignoreNestedStack {
		u, _ := AssertDeepestUnderlying(err)
		// we ignore the error of the above function for brevity, because an error
		// should never be returned from it with this usage, and the appropriate
		// action is to panic, hich will happen anyway if u is nil as in the case
		// of an error
		return u.ParentCallers()
	}
	return err.ParentCallers()
}

// ErrorStack returns a string that contains both the
// error message and the callstack.  The callstack is that of the deepest
// nested *Error, rather than that of the *Error this is called on, unless
// ignoreNestedStack is set on the *Error.
func (err *Error) ErrorStack() string {
	return err.TypeName() + " " + err.Error() + "\n" + string(err.Stack())
}

// ParentErrorStack returns a string that contains both the error message and
// the callstack.  The callstack is that of the *Error this is called on, rather
// than the deepest nested *Error.
func (err *Error) ParentErrorStack() string {
	return err.TypeName() + " " + err.Error() + "\n" + string(err.ParentStack())
}

// ParentStackFrames returns an array of frames containing information about
// the stack of the *Error this is called on.
func (err *Error) ParentStackFrames() []StackFrame {
	if err.frames == nil {
		err.frames = make([]StackFrame, len(err.stack))

		for i, pc := range err.stack {
			err.frames[i] = NewStackFrame(pc)
		}
	}

	return err.frames
}

// StackFrames returns an array of frames containing information about the
// stack of the deepest nested *Error, unless ignoreNestedStack is set on the
// *Error, in which case it is about the stack of the *Error this is called on.
func (err *Error) StackFrames() []StackFrame {
	if !err.ignoreNestedStack {
		u, _ := AssertDeepestUnderlying(err)
		// we ignore the error of the above function for brevity, because an error
		// should never be returned from it with this usage, and the appropriate
		// action is to panic, which will happen anyway if u is nil as in the case
		// of an error
		return u.ParentStackFrames()
	}
	return err.ParentStackFrames()
}

// ParentStackTrace implements a function similar that required for the
// pkg/errors.stacktracer interface.  It returns
// an array of frames containing information about the stack of the *Error this
// is called on, rather than the deepest nested *Error.
func (err *Error) ParentStackTrace() errors.StackTrace {
	st := make(errors.StackTrace, len(err.stack))

	for i, pc := range err.stack {
		st[i] = errors.Frame(pc)
	}

	return st
}

// StackTrace implements the pkg/errors.stacktracer interface.  It returns an
// array of frames containing information about the stack of the deepest nested
// *Error, unless ignoreNestedStack is set on the *Error, in which case it is
// about the stack of the *Error this is called on.
func (err *Error) StackTrace() errors.StackTrace {
	if !err.ignoreNestedStack {
		u, _ := AssertDeepestUnderlying(err)
		// we ignore the error of the above function for brevity, because an error
		// should never be returned from it with this usage, and the appropriate
		// action is to panic, which will happen anyway if u is nil as in the case
		// of an error
		return u.ParentStackTrace()
	}
	return err.ParentStackTrace()
}

// TypeName returns the type this error. e.g. *errors.stringError.
func (err *Error) TypeName() string {
	if _, ok := err.Err.(uncaughtPanic); ok {
		return "panic"
	}
	return reflect.TypeOf(err.Err).String()
}

// Cause returns the underlying cause of an error.  It returns the immediate
// cause of an error, not the "root" cause, which may be nested further.
func (err *Error) Cause() error {
	return err.Err
}

// RootCause returns the root underlying cause of an error.  It returns the
// first error in the stack of nested errors, that is not of type *Error
func (err *Error) RootCause() (root error) {
	root = err
	for {
		e, ok := Assert(root)
		if !ok {
			break
		}
		root = e.Cause()
	}
	return root
}

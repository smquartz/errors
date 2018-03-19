// Package errors provides errors that have stack-traces.
//
// This is particularly useful when you want to understand the
// state of execution when an error was returned unexpectedly.
//
// It provides the type *Error which implements the standard
// golang error interface, so you can use this library interchangably
// with code that is expecting a normal error return.
//
// For example:
//
//  package crashy
//
//  import "github.com/go-errors/errors"
//
//  var Crashed = errors.Errorf("oh dear")
//
//  func Crash() error {
//      return errors.New(Crashed)
//  }
//
// This can be called as follows:
//
//  package main
//
//  import (
//      "crashy"
//      "fmt"
//      "github.com/go-errors/errors"
//  )
//
//  func main() {
//      err := crashy.Crash()
//      if err != nil {
//          if errors.Is(err, crashy.Crashed) {
//              fmt.Println(err.(*errors.Error).ErrorStack())
//          } else {
//              panic(err)
//          }
//      }
//  }
//
// This package was original written to allow reporting to Bugsnag,
// but after I found similar packages by Facebook and Dropbox, it
// was moved to one canonical location so everyone can benefit.
package errors

import (
	"bytes"
	"reflect"
)

// The maximum number of stackframes on any error.
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

// Is detects whether the error is equal to a given error. Errors
// are considered equal by this function if they are the same object,
// or if they both contain the same error inside an errors.Error.
func Is(e error, original error) bool {

	if e == original {
		return true
	}

	if e, ok := e.(*Error); ok {
		return Is(e.Err, original)
	}

	if original, ok := original.(*Error); ok {
		return Is(e, original.Err)
	}

	return false
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

// TypeName returns the type this error. e.g. *errors.stringError.
func (err *Error) TypeName() string {
	if _, ok := err.Err.(uncaughtPanic); ok {
		return "panic"
	}
	return reflect.TypeOf(err.Err).String()
}

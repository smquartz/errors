package errors

import (
	"fmt"
	"runtime"
)

// New makes an Error from the given value. If that value is already an
// error then it will be used directly, if not, it will be passed to
// fmt.Errorf("%v"). The stacktrace will point to the line of code that
// called New.
func New(e interface{}) *Error {
	return Wrap(e, 1)
}

// Wrap makes an Error from the given value. If that value is already an
// error then it will be used directly, if not, it will be passed to
// fmt.Errorf("%v"). The skip parameter indicates how far up the stack
// to start the stacktrace. 0 is from the current call, 1 from its caller, etc.
func Wrap(e interface{}, skip int) *Error {
	var err error

	switch e := e.(type) {
	case error:
		err = e
	case nil:
		err = nil
	default:
		err = fmt.Errorf("%v", e)
	}

	stack := make([]uintptr, MaxStackDepth)
	length := runtime.Callers(2+skip, stack[:])
	return &Error{
		Err:   err,
		stack: stack[:length],
	}
}

// Wrapf makes an Error from the given value.  If that value is already
// an error then it will be used directly as the underlying error, if not,
// it will be passed to fmt.Errorf("%v").  The prefixf parameter is used to
// add a formatted prefix to the error message when calling Error().  The skip
// pameter indicates how far up the stack to start the stacktrace; 0 is from
// the current call, 1 from its caller, etc.
func Wrapf(e interface{}, prefixf string, skip int, a ...interface{}) *Error {
	err := Wrap(e, skip+1)
	err.prefix = fmt.Sprintf(prefixf, a...)
	return err
}

// Errorf creates a new error with the given message. You can use it
// as a drop-in replacement for fmt.Errorf() to provide descriptive
// errors in return values.
func Errorf(format string, a ...interface{}) *Error {
	return Wrap(fmt.Errorf(format, a...), 1)
}

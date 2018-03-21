package errors

import (
	"bytes"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

// error strings used by this file
const (
	errCauseIncorrect             = "returned cause error is not correct"
	errNotContainStack            = "does not contain appropriate stack"
	errNotContainType             = "does not contain appropriate error type"
	errIgnoreNestedStackIncorrect = "the value of the ignoreNestedStack field is incorrect"
)

// error format strings used by this file
const (
	parentErrorStackFailed     = ".ParentErrorStack() failed; %v"
	setIgnoreNestedStackFailed = ".SetIgnoreNestedStack() failed; %v"
)

func TestSetIgnoreNestedStack(t *testing.T) {
	e := New(testMsgFoo)
	e2 := New(testMsgFoo)
	e2.ignoreNestedStack = true

	if e.SetIgnoreNestedStack(true).ignoreNestedStack != true {
		t.Errorf(setIgnoreNestedStackFailed, errIgnoreNestedStackIncorrect)
	}
	if e.SetIgnoreNestedStack(true).ignoreNestedStack != true {
		t.Errorf(setIgnoreNestedStackFailed, errIgnoreNestedStackIncorrect)
	}
	if e2.SetIgnoreNestedStack(false).ignoreNestedStack != false {
		t.Errorf(setIgnoreNestedStackFailed, errIgnoreNestedStackIncorrect)
	}
	if e2.SetIgnoreNestedStack(false).ignoreNestedStack != false {
		t.Errorf(setIgnoreNestedStackFailed, errIgnoreNestedStackIncorrect)
	}
}

func TestCallers(t *testing.T) {
	// error to test with
	err := New(testMsgFoo)

	// let's check that .Callers returns the correct thing
	if !reflect.DeepEqual(err.stack, err.Callers()) {
		t.Errorf(constructorStringFailed, errStacksNotMatch)
	}

	// error to test with
	// is nested and has ignoreNestedStack set
	err = New(New(testMsgFoo)).SetIgnoreNestedStack(true)

	// let's check that .Callers returns the correct thing
	if !reflect.DeepEqual(err.stack, err.Callers()) {
		t.Errorf(constructorStringFailed, errStacksNotMatch)
	}

	// error to test with
	// is nested and does not have ignoreNestedStack
	err = New(New(testMsgFoo))

	// let's check that .Callers returns the correct thing
	if !reflect.DeepEqual(err.Err.(*Error).stack, err.Callers()) {
		t.Errorf(constructorStringFailed, errStacksNotMatch)
	}
}

func TestCause(t *testing.T) {
	// test case: *Error with underlying nil error
	if New(nil).Cause() != nil {
		t.Errorf(constructorNilFailed, errCauseIncorrect)
	}

	// test case: *Error with underlying plain error
	underlying := fmt.Errorf(testMsgFoo)
	if New(underlying).Cause() != underlying {
		t.Errorf(constructorPlainErrorFailed, errCauseIncorrect)
	}

	// test case: *Error with underlying *Error with underlying plain error
	underlying2 := New(underlying)
	if New(underlying2).Cause() != underlying2 {
		t.Errorf(constructorErrorFailed, errCauseIncorrect)
	}
}

func TestRootCause(t *testing.T) {
	// test case: *Error with underlying nil error
	if New(nil).RootCause() != nil {
		t.Errorf(constructorNilFailed, errCauseIncorrect)
	}

	// test case: *Error with underlying plain error
	underlying := fmt.Errorf(testMsgFoo)
	if New(underlying).RootCause() != underlying {
		t.Errorf(constructorPlainErrorFailed, errCauseIncorrect)
	}

	// test case: *Error with underlying *Error with underlying plain error
	underlying2 := New(underlying)
	if New(underlying2).RootCause() != underlying {
		t.Errorf(constructorErrorFailed, errCauseIncorrect)
	}

	// test case: *Error with underlying *Error with underlying *Error with
	// underlying plain error
	underlying3 := New(underlying2)
	if New(underlying3).RootCause() != underlying {
		t.Errorf(constructorErrorFailed, errCauseIncorrect)
	}

	// test case *Error with underlying *Error with underlying *Error with
	// underlying *Error with underlying plain error
	underlying4 := New(underlying3)
	if New(underlying4).RootCause() != underlying {
		t.Errorf(constructorErrorFailed, errCauseIncorrect)
	}

	// test case *Error with underlying *Error with underlying *Error with
	// underlying *Error with underlying nil error
	underlying = nil
	if New(New(New(New(underlying)))).RootCause() != underlying {
		t.Errorf(constructorErrorFailed, errCauseIncorrect)
	}
}

func TestParentErrorStack(t *testing.T) {
	// test error to use; Error constructor
	e := New(New(New(testMsgFoo)))
	// ErrorStack of e itself
	s := e.ParentErrorStack()
	if !strings.HasPrefix(s, e.TypeName()) {
		t.Errorf(parentErrorStackFailed, errNotContainType)
	}
	if !strings.Contains(s, e.Error()) {
		t.Errorf(parentErrorStackFailed, errWrongErrorMessage)
	}
	if !strings.Contains(s, string(e.ParentStack())) {
		t.Errorf(parentErrorStackFailed, errNotContainStack)
	}
}

func TestStackFormat(t *testing.T) {

	defer func() {
		err := recover()
		if err != 'a' {
			t.Fatal(err)
		}

		e, expected := Errorf("hi"), callers()

		bs := [][]uintptr{e.stack, expected}

		if err := compareStacks(bs[0], bs[1]); err != nil {
			t.Errorf("Stack didn't match")
			t.Errorf(err.Error())
		}

		stack := string(e.Stack())

		if !strings.Contains(stack, "a: b(5)") {
			t.Errorf("Stack trace does not contain source line: 'a: b(5)'")
			t.Errorf(stack)
		}
		if !strings.Contains(stack, "error_test.go:") {
			t.Errorf("Stack trace does not contain file name: 'error_test.go:'")
			t.Errorf(stack)
		}
	}()

	a()
}

func a() error {
	b(5)
	return nil
}

func b(i int) {
	c()
}

func c() {
	panic('a')
}

// compareStacks will compare a stack created using the errors package (actual)
// with a reference stack created with the callers function (expected). The
// first entry is compared inexact since the actual and expected stacks cannot
// be created at the exact same program counter position so the first entry
// will always differ somewhat. Returns nil if the stacks are equal enough and
// an error containing a detailed error message otherwise.
func compareStacks(actual, expected []uintptr) error {
	if len(actual) != len(expected) {
		return stackCompareError("Stacks does not have equal length", actual, expected)
	}
	for i, pc := range actual {
		if i == 0 {
			firstEntryDiff := (int)(expected[i]) - (int)(pc)
			if firstEntryDiff < -27 || firstEntryDiff > 27 {
				return stackCompareError(fmt.Sprintf("First entry PC diff to large (%d)", firstEntryDiff), actual, expected)
			}
		} else if pc != expected[i] {
			return stackCompareError(fmt.Sprintf("Stacks does not match entry %d (and maybe others)", i), actual, expected)
		}
	}
	return nil
}

func stackCompareError(msg string, actual, expected []uintptr) error {
	return fmt.Errorf("%s\nActual stack trace:\n%s\nExpected stack trace:\n%s", msg, readableStackTrace(actual), readableStackTrace(expected))
}

func callers() []uintptr {
	return callersSkip(1)
}

func callersSkip(skip int) []uintptr {
	callers := make([]uintptr, MaxStackDepth)
	length := runtime.Callers(skip+2, callers[:])
	return callers[:length]
}

func readableStackTrace(callers []uintptr) string {
	var result bytes.Buffer
	frames := callersToFrames(callers)
	for _, frame := range frames {
		result.WriteString(fmt.Sprintf("%s:%d (%#x)\n\t%s\n", frame.File, frame.Line, frame.PC, frame.Function))
	}
	return result.String()
}

func callersToFrames(callers []uintptr) []runtime.Frame {
	frames := make([]runtime.Frame, 0, len(callers))
	framesPtr := runtime.CallersFrames(callers)
	for {
		frame, more := framesPtr.Next()
		frames = append(frames, frame)
		if !more {
			return frames
		}
	}
}

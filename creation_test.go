package errors

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

// error strings used in tests
// declared globally to allow reuse for consistency
const (
	errNil                  = "<nil>"
	errStacksNotMatch       = "stacks did not match"
	errWrongErrorMessage    = "produces error with the wrong error message"
	errErrStackWrongFormat  = ".ErrorStack() is in the wrong format"
	errWrongUnderlyingError = "produces error with the wrong underlying error"
	errNotMatchWrap         = "the error returned does not match the equivelant returned by Wrap()"
	errSkipFailed           = "failed to successfully skip parts of the stack"
)

// constructor failed format strings used in tests
// declared globally to allow reuse for consistency
var (
	constructorStringFailed     = "constructor with a string failed; %v"
	constructorPlainErrorFailed = "constructor with a plain error failed; %v"
	constructorErrorFailed      = "constructor with a *Error failed; %v"
	constructorNilFailed        = "constructor with nil failed; %v"
)

// test error message contents used in tests
const (
	testMsgFoo             = "foo barred a baz"
	testPrefixFoobar       = "foobar"
	testFormatPrefixFoobar = "foo %v the bar"
	testFormatArgumentBaz  = "baz"
)

func TestSkipWorks(t *testing.T) {

	defer func() {
		err := recover()
		if err != 'a' {
			t.Fatal(err)
		}

		bs := [][]uintptr{Wrap(testMsgFoo, 2).stack, callersSkip(2)}

		if err := compareStacks(bs[0], bs[1]); err != nil {
			t.Errorf(errStacksNotMatch)
			t.Errorf(err.Error())
		}
	}()

	a()
}

func TestNew(t *testing.T) {
	// *Error type returned as error, used in this test
	e := func() error {
		return New(testMsgFoo)
	}()

	// test using a string constructor

	// test that the error message is as it was set
	if e.Error() != testMsgFoo {
		t.Errorf(constructorStringFailed, errWrongErrorMessage)
	}
	// test the underlying error is a plain one as returned by fmt.Errorf()
	if !reflect.DeepEqual(e.(*Error).Err, fmt.Errorf(testMsgFoo)) {
		t.Errorf(constructorStringFailed, errWrongUnderlyingError)
	}

	// test using a plain error constructor

	// test that the error message is as it was set
	if New(fmt.Errorf(testMsgFoo)).Error() != testMsgFoo {
		t.Errorf(constructorPlainErrorFailed, errWrongErrorMessage)
	}
	// test that the underlying error is as it was set
	if !reflect.DeepEqual(New(fmt.Errorf(testMsgFoo)).Err, fmt.Errorf(testMsgFoo)) {
		t.Errorf(constructorPlainErrorFailed, errWrongUnderlyingError)
	}

	// test using a *Error constructor
	err := New(e)
	// test that the underlying error is as it was set
	if !reflect.DeepEqual(err.Err, e) {
		t.Errorf(constructorErrorFailed, errWrongUnderlyingError)
	}
	// test that the error message is correct
	if err.Error() != e.Error() {
		t.Errorf(constructorErrorFailed, errWrongErrorMessage)
	}

	// test using a nil constructor

	// test that the error message is correct
	if New(nil).Error() != errNil {
		t.Errorf(constructorNilFailed, errWrongErrorMessage)
	}
	// test that the underlying error is as it was set
	if New(nil).Err != nil {
		t.Errorf(constructorNilFailed, errWrongUnderlyingError)
	}

	// create a slice to compare New()'s stack, with callers()'s stack
	bs := [][]uintptr{New(testMsgFoo).stack, callers()}
	// actually compare New()'s stack with caller()'s stack
	if err := compareStacks(bs[0], bs[1]); err != nil {
		t.Errorf(constructorStringFailed, errStacksNotMatch)
		t.Errorf(err.Error())
	}
	// check .ErrorStack() retruns the correct format
	if err.ErrorStack() != err.TypeName()+" "+err.Error()+"\n"+string(err.Stack()) {
		t.Errorf(constructorErrorFailed, errErrStackWrongFormat)
	}
	// check .ErrorStack() returns the correct format if ignoreNestedStack is true
	if err.SetIgnoreNestedStack(true).ErrorStack() != err.TypeName()+" "+err.Error()+"\n"+string(err.ParentStack()) {
		t.Errorf(constructorErrorFailed, errErrStackWrongFormat)
	}
}

func TestWrapError(t *testing.T) {
	// *Error type returned as error, used in this test
	e := func() error {
		return Wrap(testMsgFoo, 1)
	}()

	// test using a string constructor

	// test that the error message is as it was set
	if e.Error() != testMsgFoo {
		t.Errorf(constructorStringFailed, errWrongErrorMessage)
	}
	// test the underlying error is a plain one as returned by fmt.Errorf()
	if !reflect.DeepEqual(e.(*Error).Err, fmt.Errorf(testMsgFoo)) {
		t.Errorf(constructorStringFailed, errWrongUnderlyingError)
	}

	// test using a plain error constructor

	// test that the error message is as it was set
	if Wrap(fmt.Errorf(testMsgFoo), 0).Error() != testMsgFoo {
		t.Errorf(constructorPlainErrorFailed, errWrongErrorMessage)
	}
	// test that the underlying error is as it was set
	if !reflect.DeepEqual(Wrap(fmt.Errorf(testMsgFoo), 0).Err, fmt.Errorf(testMsgFoo)) {
		t.Errorf(constructorPlainErrorFailed, errWrongUnderlyingError)
	}

	// test using a *Error constructor
	err := Wrap(e, 0)
	// test that the underlying error is as it was set
	if !reflect.DeepEqual(err.Err, e) {
		t.Errorf(constructorErrorFailed, errWrongUnderlyingError)
	}
	// test that the error message is correct
	if err.Error() != e.Error() {
		t.Errorf(constructorErrorFailed, errWrongErrorMessage)
	}

	// test using a nil constructor

	// test that the error message is correct
	if Wrap(nil, 0).Error() != errNil {
		t.Errorf(constructorNilFailed, errWrongErrorMessage)
	}
	// test that the underlying error is as it was set
	if Wrap(nil, 0).Err != nil {
		t.Errorf(constructorNilFailed, errWrongUnderlyingError)
	}

	// create a slice to compare Wrap()'s stack, with callers()'s stack
	bs := [][]uintptr{Wrap(testMsgFoo, 0).stack, callers()}
	// actually compare Wrap()'s stack with caller()'s stack
	if err := compareStacks(bs[0], bs[1]); err != nil {
		t.Errorf(constructorStringFailed, errStacksNotMatch)
		t.Errorf(err.Error())
	}
	// check .ErrorStack() retruns the correct format
	if err.ErrorStack() != err.TypeName()+" "+err.Error()+"\n"+string(err.Stack()) {
		t.Errorf(constructorErrorFailed, errErrStackWrongFormat)
	}
	// check .ErrorStack() returns the correct format if ignoreNestedStack is true
	if err.SetIgnoreNestedStack(true).ErrorStack() != err.TypeName()+" "+err.Error()+"\n"+string(err.ParentStack()) {
		t.Errorf(constructorErrorFailed, errErrStackWrongFormat)
	}
}

func TestWrapfError(t *testing.T) {
	// *Error type returned as error, used in this test
	e := func() error {
		return Wrapf(testMsgFoo, testFormatPrefixFoobar, 1, testFormatArgumentBaz)
	}()

	// test using a string constructor

	// test that the error message is as it was set
	if e.Error() != fmt.Sprintf("%s: %s", fmt.Sprintf(testFormatPrefixFoobar, testFormatArgumentBaz), testMsgFoo) {
		t.Errorf(constructorStringFailed, errWrongErrorMessage)
	}
	// test the underlying error is a plain one as returned by fmt.Errorf()
	if !reflect.DeepEqual(e.(*Error).Err, fmt.Errorf(testMsgFoo)) {
		t.Errorf(constructorStringFailed, errWrongUnderlyingError)
	}

	// test using a plain error constructor
	err := Wrapf(fmt.Errorf(testMsgFoo), testFormatPrefixFoobar, 0, testFormatArgumentBaz)
	// test that the error message is as it was set
	if err.Error() != fmt.Sprintf("%s: %s", fmt.Sprintf(testFormatPrefixFoobar, testFormatArgumentBaz), testMsgFoo) {
		t.Errorf(constructorPlainErrorFailed, errWrongErrorMessage)
	}
	// test that the underlying error is as it was set
	if !reflect.DeepEqual(err.Err, fmt.Errorf(testMsgFoo)) {
		t.Errorf(constructorPlainErrorFailed, errWrongUnderlyingError)
	}

	// test using a nil constructor
	err = Wrapf(nil, testFormatPrefixFoobar, 0, testFormatArgumentBaz)
	// test that the error message is correct
	if err.Error() != fmt.Sprintf("%s: %s", fmt.Sprintf(testFormatPrefixFoobar, testFormatArgumentBaz), fmt.Errorf("%v", nil).Error()) {
		t.Errorf(constructorNilFailed, errWrongErrorMessage)
	}
	// test that the underlying error is as it was set
	if err.Err != nil {
		t.Errorf(constructorNilFailed, errWrongUnderlyingError)
	}

	// test using a *Error constructor
	err = Wrapf(e, testFormatPrefixFoobar, 0, testFormatArgumentBaz)
	// test that the underlying error is as it was set
	if !reflect.DeepEqual(err.Err, e) {
		t.Errorf(constructorErrorFailed, errWrongUnderlyingError)
	}
	// test that the error message is correct
	if err.Error() != fmt.Sprintf("%s: %s", fmt.Sprintf(testFormatPrefixFoobar, testFormatArgumentBaz), e.Error()) {
		t.Errorf(constructorErrorFailed, errWrongErrorMessage)
	}

	// create a slice to compare Wrapf()'s stack, with callers()'s stack
	bs := [][]uintptr{Wrapf(testMsgFoo, testFormatPrefixFoobar, 0).stack, callers()}
	// actually compare Wrapf()'s stack with caller()'s stack
	if err := compareStacks(bs[0], bs[1]); err != nil {
		t.Errorf(constructorStringFailed, errStacksNotMatch)
		t.Errorf(err.Error())
	}
	// check .ErrorStack() retruns the correct format
	if err.ErrorStack() != err.TypeName()+" "+err.Error()+"\n"+string(err.Stack()) {
		t.Errorf(constructorErrorFailed, errErrStackWrongFormat)
	}
	// check .ErrorStack() returns the correct format if ignoreNestedStack is true
	if err.SetIgnoreNestedStack(true).ErrorStack() != err.TypeName()+" "+err.Error()+"\n"+string(err.ParentStack()) {
		t.Errorf(constructorErrorFailed, errErrStackWrongFormat)
	}

	original := e.(*Error)
	if !strings.HasSuffix(original.StackFrames()[0].File, "creation_test.go") || strings.HasSuffix(original.StackFrames()[1].File, "creation_test.go") {
		t.Errorf(constructorStringFailed, errSkipFailed)
	}
	original.SetIgnoreNestedStack(true)
	if !strings.HasSuffix(original.StackFrames()[0].File, "creation_test.go") || strings.HasSuffix(original.StackFrames()[1].File, "creation_test.go") {
		t.Errorf(constructorStringFailed, errSkipFailed)
	}
	if !strings.HasSuffix(NewStackFrame(uintptr(original.StackTrace()[0])).File, "creation_test.go") || strings.HasSuffix(NewStackFrame(uintptr(original.StackTrace()[1])).File, "creation_test.go") {
		t.Errorf(constructorStringFailed, errSkipFailed)
	}
	original.SetIgnoreNestedStack(false)
	if !strings.HasSuffix(NewStackFrame(uintptr(original.StackTrace()[0])).File, "creation_test.go") || strings.HasSuffix(NewStackFrame(uintptr(original.StackTrace()[1])).File, "creation_test.go") {
		t.Errorf(constructorStringFailed, errSkipFailed)
	}
}

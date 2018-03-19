package errors

import (
	"fmt"
	"reflect"
	"testing"
)

// error strings used is this file
const (
	errNilAsserted                                    = "nil was successfully asserted to *Error"
	errPlainErrorAsserted                             = "a plain error was successfully asserted to *Error"
	errErrorNotAsserted                               = "a *Error failed to be asserted to a *Error"
	errAssertedErrorNotMatch                          = "a *Error was asserted to a *Error that did not match the actual *Error"
	errNilAssertedWithUnderlying                      = errNilAsserted + " with an underlying *Error"
	errErrorNotAppropriate                            = "error returned was not appropriate given the conditions"
	errPlainErrorAssertedWithUnderlying               = errPlainErrorAsserted + " with an underlying *Error"
	errErrorWithPlainUnderlyingAssertedWithUnderlying = "a *Error with a plain underlying error was asserted as a *Error with an underlying *Error"
	errUnderlyingErrorNotAsserted                     = "a *Error's underlying *Error failed to be asserted to a *Error"
)

// error format strings used in this file
const (
	assertFailed            = "Assert() failed; %v"
	assertUnderlyingFailed  = "AssertUnderlying() failed; %v"
	getRootUnderlyingFailed = "GetRootUnderlying() failed; %v"
	gRUFNilFailed           = "GetRootUnderlying() failed with nil error as argument; %v"
	gRUFPlainFailed         = "GetRootUnderlying() failed with plain error as argument; %v"
	gRUFErrorFailed         = "GetRootUnderlying() failed with *Error with nested plain error as argument; %v"
	gRUFNestedErrorFailed   = "GetRootUnderlying() failed with *Error with nested *Error(s) as argument; %v"
)

func TestAssert(t *testing.T) {
	var err error

	// test case: nil error
	err = nil
	_, ok := Assert(err)
	if ok {
		t.Errorf(assertFailed, errNilAsserted)
	}

	// test case: plain error
	err = fmt.Errorf(testMsgFoo)
	_, ok = Assert(err)
	if ok {
		t.Errorf(assertFailed, errPlainErrorAsserted)
	}

	// test case: *Error
	actual := New(testMsgFoo)
	err = func() error { return actual }()
	e, ok := Assert(err)
	if !ok {
		t.Errorf(assertFailed, errErrorNotAsserted)
	}
	if e != actual {
		t.Errorf(assertFailed, errAssertedErrorNotMatch)
	}
}

func TestAssertUnderlying(t *testing.T) {
	var parentError error

	// test case: nil error
	parentError = nil
	_, err := AssertUnderlying(parentError)
	if err == nil {
		t.Errorf(assertUnderlyingFailed, errNilAssertedWithUnderlying)
	} else if !Is(err, ErrNotError) {
		t.Errorf(assertUnderlyingFailed, errErrorNotAppropriate)
	}

	// test case: plain error
	parentError = fmt.Errorf(testMsgFoo)
	_, err = AssertUnderlying(parentError)
	if err == nil {
		t.Errorf(assertUnderlyingFailed, errPlainErrorAssertedWithUnderlying)
	} else if !Is(err, ErrNotError) {
		t.Errorf(assertUnderlyingFailed, errErrorNotAppropriate)
	}

	// test case: *Error with plain underlying
	parentError = func() error { return New(fmt.Errorf(testMsgFoo)) }()
	_, err = AssertUnderlying(parentError)
	if err == nil {
		t.Errorf(assertUnderlyingFailed, errErrorWithPlainUnderlyingAssertedWithUnderlying)
	} else if !Is(err, ErrUnderlyingNotError) {
		t.Errorf(assertUnderlyingFailed, errErrorNotAppropriate)
	}

	// test case: *Error with *Error underlying
	underlying := parentError
	parentError = func() error { return WrapPrefix(underlying, testPrefixFoobar, 1) }()
	u, err := AssertUnderlying(parentError)
	if err != nil {
		if Is(err, ErrNotError) {
			t.Errorf(assertUnderlyingFailed, errErrorNotAsserted)
		} else if Is(err, ErrUnderlyingNotError) {
			t.Errorf(assertUnderlyingFailed, errUnderlyingErrorNotAsserted)
		} else {
			t.Errorf(assertUnderlyingFailed, errErrorNotAppropriate)
		}
	}
	if !reflect.DeepEqual(u, underlying) {
		t.Errorf(assertUnderlyingFailed, errWrongUnderlyingError)
	}

}

func TestGetRootUnderlying(t *testing.T) {
	var err error
	// test case: nil error
	err = nil
	if root := GetRootUnderlying(err); root != err {
		t.Errorf(gRUFNilFailed, errWrongUnderlyingError)
	}

	// test case: plain error
	err = fmt.Errorf(testMsgFoo)
	if root := GetRootUnderlying(err); root != err {
		t.Errorf(gRUFPlainFailed, errWrongUnderlyingError)
	}

	// test case: *Error with underlying plain error
	underlying := fmt.Errorf(testMsgFoo)
	err = func() error { return New(underlying) }()
	if root := GetRootUnderlying(err); root != underlying {
		t.Errorf(gRUFErrorFailed, errWrongUnderlyingError)
	}

	// test case: *Error with underlying *Error with underlying plain error
	underlying2 := New(underlying)
	err = New(underlying2)
	if root := GetRootUnderlying(err); root != underlying {
		t.Errorf(gRUFNestedErrorFailed, errWrongUnderlyingError)
	}

	// test case: *Error with underlying *Error with underlying *Error with
	// underlying plain error
	underlying3 := New(underlying2)
	err = New(underlying3)
	if root := GetRootUnderlying(err); root != underlying {
		t.Errorf(gRUFNestedErrorFailed, errWrongUnderlyingError)
	}
}

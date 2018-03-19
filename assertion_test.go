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
	assertFailed           = "Assert() failed; %v"
	assertUnderlyingFailed = "AssertUnderlying() failed; %v"
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

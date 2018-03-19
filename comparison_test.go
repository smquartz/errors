package errors

import (
	"fmt"
	"io"
	"testing"
)

// error format strings used in this file for consistency
const (
	evaluatedIncorrectlySame      = "evaluated %v as the same as %v when it is not"
	evaluatedIncorrectlyDifferent = "evaluated %v as different to %v when it is not"
)

func TestIs(t *testing.T) {

	if Is(nil, io.EOF) {
		t.Errorf(evaluatedIncorrectlySame, "nil", "io.EOF")
	}

	if !Is(io.EOF, io.EOF) {
		t.Errorf(evaluatedIncorrectlyDifferent, "io.EOF", "io.EOF")
	}

	if !Is(io.EOF, New(io.EOF)) {
		t.Errorf(evaluatedIncorrectlyDifferent, "io.EOF", "New(io.EOF)")
	}

	if !Is(New(io.EOF), New(io.EOF)) {
		t.Errorf(evaluatedIncorrectlyDifferent, "New(io.EOF)", "New(io.EOF)")
	}

	if Is(io.EOF, fmt.Errorf("io.EOF")) {
		t.Errorf(evaluatedIncorrectlySame, "io.EOF", "fmt.Errorf(\"io.EOF\")")
	}

	if !Is(New(New(io.EOF)), io.EOF) {
		t.Errorf(evaluatedIncorrectlyDifferent, "New(New(io.EOF))", "io.EOF")
	}

	if !Is(New(New(io.EOF)), Wrapf(Wrapf(io.EOF, testMsgFoo, 1), testPrefixFoobar, 1)) {
		t.Errorf(evaluatedIncorrectlyDifferent, "New(New(io.EOF))", "WrapPrefix(WrapPrefix(io.EOF, testMsgFoo, 1), testPrefixFoobar, 1)")
	}

}

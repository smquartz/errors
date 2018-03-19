package errors

import (
	"fmt"
	"testing"
)

func TestNewStackFrameNilProgramCounter(t *testing.T) {
	NewStackFrame(0)
}

func TestStackFrameStringWithNonsenseFile(t *testing.T) {
	frame := NewStackFrame(0)
	str := frame.String()
	if str != fmt.Sprintf("%s:%d (0x%x)\n", frame.File, frame.LineNumber, frame.ProgramCounter) {
		t.Errorf("frame.String() was somehow able to read a nonexistent file")
	}
}

func TestStackFromeSourceLineWithNonsenseLineNumber(t *testing.T) {
	err := New(nil)
	frame := err.StackFrames()[0]
	frame.LineNumber = -42
	str, _ := frame.SourceLine()
	if str != "???" {
		t.Errorf("frame.SourceLine() did not recognise a nonsense line number")
	}
}

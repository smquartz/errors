package errors

import "fmt"

func ExampleError_Error(err error) {
	fmt.Println(err.Error())
}

func ExampleError_ErrorStack(err error) {
	fmt.Println(err.(*Error).ErrorStack())
}

func ExampleError_Stack(err *Error) {
	fmt.Println(err.Stack())
}

func ExampleError_TypeName(err *Error) {
	fmt.Println(err.TypeName(), err.Error())
}

func ExampleError_StackFrames(err *Error) {
	for _, frame := range err.StackFrames() {
		fmt.Println(frame.File, frame.LineNumber, frame.Package, frame.Name)
	}
}

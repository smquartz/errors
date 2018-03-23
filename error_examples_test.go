package errors

import "fmt"

func ExampleError_Error(err error) {
	fmt.Println(err.Error())
}

func ExampleError_ErrorStack(err error) {
	fmt.Println(err.(*Err).ErrorStack())
}

func ExampleError_Stack(err *Err) {
	fmt.Println(err.Stack())
}

func ExampleError_TypeName(err *Err) {
	fmt.Println(err.TypeName(), err.Error())
}

func ExampleError_StackFrames(err *Err) {
	for _, frame := range err.StackFrames() {
		fmt.Println(frame.File, frame.LineNumber, frame.Package, frame.Name)
	}
}

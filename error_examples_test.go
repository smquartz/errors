package errors

import "fmt"

func ExampleError_Error() {
	// example error returning function
	err := func() error { return New("hai") }()
	fmt.Println(err.Error())
}

func ExampleError_ErrorStack() {
	// example error returning function
	err := func() error { return New("hai") }()
	// print the error stack
	fmt.Println(err.(*Err).ErrorStack())
}

func ExampleError_Stack() {
	// example *Err returning function
	err := func() *Err { return New("hai") }()
	// print the stack
	fmt.Println(err.Stack())
}

func ExampleError_TypeName() {
	// example *Err returning function
	err := func() *Err { return New("hai") }()
	// get the type name
	fmt.Println(err.TypeName(), err.Error())
}

func ExampleError_StackFrames() {
	// example *Err returning function
	err := func() *Err { return New("hai") }()
	// print all the frames
	for _, frame := range err.StackFrames() {
		fmt.Println(frame.File, frame.LineNumber, frame.Package, frame.Name)
	}
}

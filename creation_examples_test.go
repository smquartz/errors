package errors

import (
	"fmt"
	"io"
)

func ExampleNew() {
	// call some error returning function
	err := func() error { return io.ErrUnexpectedEOF }()
	// calling New attaches the current stacktrace to the existing UnexpectedEOF error
	err = New(err)
	// do something with the error
	fmt.Println(err)
}

func ExampleWrap() {

	// if recovered from panic
	if err := recover(); err != nil {
		// wrap the error, adding stacktrace, skipping one frame
		err = Wrap(err, 1)
		// do something with the error
		fmt.Println(err)
	}

	// else, do something else
	fmt.Println(a())
}

func ExampleErrorf() {
	// example function
	halve := func(x int) (int, error) {
		// if number cannot be halved without remainedr
		if x%2 != 0 {
			return 0, Errorf("cannot halve %v without remainder", x)
		}
		// else, return halved number
		return x / 2, nil
	}

	// call the function
	val, err := halve(3)
	// do something with the error
	if err != nil {
		fmt.Println("halve(3) failed", err)
	} else {
		fmt.Println("halve(3) worked:", val)
	}
}

func ExampleWrapError() {
	// example function that returns error
	example := func() (int, error) {
		// Wrap io.EOF with the current stack-trace and return it
		return 0, Wrap(io.EOF, 0)
	}

	// call the function
	_, err := example()
	// do something with the error
	fmt.Println(err)
}

func ExampleWrapError_skip() {
	defer func() {
		if err := recover(); err != nil {
			// skip 1 frame (the deferred function) and then return the wrapped err
			err = Wrap(err, 1)
			// do something with the error
			fmt.Println(err)
		}
	}()
}

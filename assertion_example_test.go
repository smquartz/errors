package errors

import "fmt"

func ExampleAssert() {
	// mock function that returns a error that is actually a *Error
	err := func() error { return New("some error") }()
	// try assert it
	e, ok := Assert(err)
	if !ok {
		// oh well
		return
	}
	// print the callers
	fmt.Println(e.Callers())
}

func ExampleAssertUnderlying() {
	// mock function that returns an error that is actually a *Error
	crashy := func() error { return New("some error") }
	// mock function that returns an error that is actually a *Error
	// this one relies upon and wraps errors of crashy
	biggerCrashy := func() error {
		err := crashy()
		if err != nil {
			return Wrap(err, 1)
		}
		return nil
	}

	// let's actually call biggerCRashy
	err := biggerCrashy()
	// now we want to access the callers of its underlying error
	e, err := AssertUnderlying(err)
	if err != nil {
		// oh well
		return
	}
	// print the callers
	fmt.Println(e.Callers())
}

func ExampleGetRootUnderlying() {
	// mock function that returns a heavily nested error, with a root error
	// that is a plain error
	err := func() error {
		return WrapPrefix(WrapPrefix(WrapPrefix(WrapPrefix(fmt.Errorf("some error"), "1", 1), "2", 1), "3", 1), "4", 1)
	}()

	// let's print the error
	fmt.Printf("The error is %v\n", err)
	// oh no it has all these prefixes we don't want
	// let's get the root underlying error
	err = GetRootUnderlying(err)
	// let's print the root underlying error
	fmt.Printf("The root error is %v\n", err)
}

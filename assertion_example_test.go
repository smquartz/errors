package errors

import "fmt"

func ExampleAssert() {
	// mock function that returns a error that is actually a *Err
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
	// mock function that returns an error that is actually a *Err
	crashy := func() error { return New("some error") }
	// mock function that returns an error that is actually a *Err
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

package errors

import "io"

func ExampleNew(UnexpectedEOF error) error {
	// calling New attaches the current stacktrace to the existing UnexpectedEOF error
	return New(UnexpectedEOF)
}

func ExampleWrap() error {

	if err := recover(); err != nil {
		return Wrap(err, 1)
	}

	return a()
}

func ExampleErrorf(x int) (int, error) {
	if x%2 == 1 {
		return 0, Errorf("can only halve even numbers, got %d", x)
	}
	return x / 2, nil
}

func ExampleWrapError() (error, error) {
	// Wrap io.EOF with the current stack-trace and return it
	return nil, Wrap(io.EOF, 0)
}

func ExampleWrapError_skip() {
	defer func() {
		if err := recover(); err != nil {
			// skip 1 frame (the deferred function) and then return the wrapped err
			err = Wrap(err, 1)
		}
	}()
}

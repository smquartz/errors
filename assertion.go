package errors

// errors used in this file
var (
	ErrNotError           = Errorf("the provided error is not of type *Error")
	ErrUnderlyingNotError = Errorf("the provided error's underlying error is not of type *Error")
)

// Assert is a convenience function that attempts to assert a error to a *Error.
func Assert(err error) (*Error, bool) {
	e, ok := err.(*Error)
	if !ok {
		return nil, ok
	}
	return e, ok
}

// AssertUnderlying is a convenience function that attempts to assert
// an error to a *Error, and then attempts to assert its underlying
// error to a *Error.
func AssertUnderlying(err error) (*Error, error) {
	e, ok := Assert(err)
	if !ok {
		return nil, New(ErrNotError)
	}
	u, ok := Assert(e.Err)
	if !ok {
		return nil, New(ErrUnderlyingNotError)
	}
	return u, nil
}

// GetRootUnderlying is a convenience function that will return the shallowest
// nested error that is not of type *Error.  The aim of this is to return the
// root cause error.
func GetRootUnderlying(err error) (root error) {
	root = err
	for {
		e, ok := Assert(root)
		if !ok {
			break
		}
		root = e.Err
	}
	return root
}

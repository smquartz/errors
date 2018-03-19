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

// AssertNthUnderlying is a convenience function that attempts to assert
// an error to a *Error, and then attempts to recursively assert its
// underlying errors to a *Error, up to the nth (specified) underlying error.
func AssertNthUnderlying(err error, nth int) (u *Error, ierr error) {
	u, ok := Assert(err)
	if !ok {
		return nil, New(ErrNotError)
	}
	for i := 0; i < nth; i++ {
		u, ierr = AssertUnderlying(u)
		if ierr != nil {
			return nil, New(ErrUnderlyingNotError)
		}
	}
	return u, nil
}

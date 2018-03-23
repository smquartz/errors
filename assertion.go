package errors

// ErrNotErr is used when a parent error is passed to an Assert* function that
// is not of type *Err
type ErrNotErr struct {
	*Err
}

// newErrNotErr returns a new initialised instance of ErrNotErr
func newErrNotErr() *ErrNotErr {
	return &ErrNotErr{Err: Errorf("the provided error is not of type *Err")}
}

// ErrUnderlyingNotErr is used when a parent error is passed to an Assert*
// function that contains an underlying error that is not of type *Err
type ErrUnderlyingNotErr struct {
	*Err
}

// newErrUnderlyingNotErr returns a new initialised instance of ErrUnderlyingNotErr
func newUnderlyingNotErr() *ErrUnderlyingNotErr {
	return &ErrUnderlyingNotErr{Err: Errorf("the provided error is not of type *Err")}
}

// Assert is a convenience function that attempts to assert a error to a *Err.
func Assert(err error) (*Err, bool) {
	e, ok := err.(*Err)
	if !ok {
		return nil, ok
	}
	return e, ok
}

// AssertUnderlying is a convenience function that attempts to assert
// an error to a *Err, and then attempts to assert its underlying
// error to a *Err.
func AssertUnderlying(err error) (*Err, error) {
	e, ok := Assert(err)
	if !ok {
		return nil, newErrNotErr()
	}
	u, ok := Assert(e.Underlying)
	if !ok {
		return nil, newUnderlyingNotErr()
	}
	return u, nil
}

// AssertNthUnderlying is a convenience function that attempts to assert
// an error to a *Err, and then attempts to recursively assert its
// underlying errors to a *Err, up to the nth (specified) underlying error.
func AssertNthUnderlying(err error, nth int) (u *Err, ierr error) {
	u, ok := Assert(err)
	if !ok {
		return nil, newErrNotErr()
	}
	for i := 0; i < nth; i++ {
		u, ierr = AssertUnderlying(u)
		if ierr != nil {
			return nil, newUnderlyingNotErr()
		}
	}
	return u, nil
}

// AssertDeepestUnderlying is a convenience function that attempts to return
// the deepest underlying *Err in a stack of errors.
func AssertDeepestUnderlying(err error) (u *Err, ierr error) {
	u, ok := Assert(err)
	if !ok {
		return nil, newErrNotErr()
	}
	var u2 *Err
	for ierr == nil {
		u2, ierr = AssertUnderlying(u)
		if u2 != nil {
			u = u2
		}
	}
	return u, nil
}

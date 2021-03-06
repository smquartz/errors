package errors

// Is detects whether the error is equal to a given error. Errors
// are considered equal by this function if they are the same object,
// or if they both contain the same error inside an errors.Error.
func Is(e error, original error) bool {

	if e == original {
		return true
	}

	if e, ok := e.(*Err); ok {
		return Is(e.Underlying, original)
	}

	if original, ok := original.(*Err); ok {
		return Is(e, original.Underlying)
	}

	return false
}

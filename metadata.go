package errors

// Metadata is an interface intended to be implemented by applications that use
// this package, which is used to store arbitrary useful information in
// errors
type Metadata interface {
	IsMetadata()
}

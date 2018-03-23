package errors

import (
	"github.com/bugsnag/bugsnag-go"
)

// Metadata is an interface intended to be implemented by applications that use
// this package, which is used to store arbitrary useful information in
// errors
type Metadata interface {
	IsMetadata()
}

// SetMeta sets the metadata on the called-upon error, as well as returning a
// pointer to that same error, to allow it to easily be chained inline.
func (e *Err) SetMeta(meta Metadata) *Err {
	e.Metadata = meta
	return e
}

// DefaultMetadata is a default implementation of the Metadata interface.
// It is compatible with bugsnag.MetaData, allowing you to pass it directly
// to bugsnag functions as rawData
type DefaultMetadata bugsnag.MetaData

// IsMetadata is a dummy function that is used to satisfy the Metadata
// interface.
func (DefaultMetadata) IsMetadata() {}

// NewDefaultMetadata returns an initialized instance of the default Metadata
// implementation.  It returns an unexported metadata type, as Metadata.
func NewDefaultMetadata() DefaultMetadata {
	return make(DefaultMetadata)
}

package errors

import (
	"reflect"
	"testing"
)

// error strings used in this file
const (
	errMetadataNotMatch               = "metadata does not match what it was set to"
	errMetadataNotInitializedProperly = "metadata was not initialized properly"
)

// error format strings used in this file
const (
	newDefaultMetadataFailed = "NewDefaultMetadata() failed; %v"
	isMetadataFailed         = "IsMetadata() failed; %v"
)

func TestSetMeta(t *testing.T) {
	// test error to use
	e := New(testMsgFoo)
	// test metadata to use
	m := DefaultMetadata{"foo": {"bar": "baz"}}
	// second set of test metadata to use
	m2 := DefaultMetadata{"baz": {"bar": "foo"}}

	// test that setting its metadata works
	// set the metadata
	e.SetMeta(m)
	if !reflect.DeepEqual(m, e.Metadata) {
		t.Errorf(constructorStringFailed, errMetadataNotMatch)
	}

	// noting that it already has metadata, let's try set it again
	e.SetMeta(m2)
	if !reflect.DeepEqual(m2, e.Metadata) {
		t.Errorf(constructorStringFailed, errMetadataNotMatch)
	}
}

func TestNewDefaultMetadata(t *testing.T) {
	// test metadata to compare to
	m := make(DefaultMetadata)

	// compare the two to output of NewDefaultMetadata()
	if !reflect.DeepEqual(NewDefaultMetadata(), m) {
		t.Errorf(newDefaultMetadataFailed, errMetadataNotInitializedProperly)
	}
}

func TestIsMetadata(t *testing.T) {
	// this test is mostly a dummy test as IsMetadata() should do absolutely
	// nothing in the default implementation

	// get some metadata to test with
	m := NewDefaultMetadata()

	// call IsMetadata()
	m.IsMetadata()
	if false {
		t.Errorf(isMetadataFailed, "this cannot happen")
	}
}

package errors

import (
	"bytes"
	"io"
)

func ExampleIs() {
	// setup a dummy reader
	reader := bytes.NewReader(nil)
	// setup buffer
	var buf []byte

	// read from the reader
	_, err := reader.Read(buf)
	// if underlying error is io.EOF
	if Is(err, io.EOF) {
		return
	}
}

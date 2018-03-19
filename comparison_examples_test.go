package errors

import "io"

func ExampleIs(reader io.Reader, buff []byte) {
	_, err := reader.Read(buff)
	if Is(err, io.EOF) {
		return
	}
}

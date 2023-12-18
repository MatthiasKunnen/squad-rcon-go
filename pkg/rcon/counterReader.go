package rcon

import (
	"io"
)

// countingReader wraps io.Reader and keeps count of the amount of bytes read.
type countingReader struct {
	io.Reader
	Bytes          []byte
	TotalBytesRead int64
}

// Read reads exactly len(output) bytes into output.
// It returns the number of bytes copied and an error if fewer bytes were read.
func (c *countingReader) Read(output []byte) (n int, err error) {
	n, err = io.ReadFull(c.Reader, output)
	c.Bytes = append(c.Bytes, output...)
	c.TotalBytesRead += int64(n)
	return n, err
}

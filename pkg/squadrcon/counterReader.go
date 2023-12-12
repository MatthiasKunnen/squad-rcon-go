package squadrcon

import (
	"io"
)

// countingReader wraps io.Reader and keeps count of the amount of bytes read.
type countingReader struct {
	io.Reader
	Bytes          []byte
	TotalBytesRead int64
}

func (c *countingReader) Read(output []byte) (n int, err error) {
	n, err = c.Reader.Read(output)
	c.Bytes = append(c.Bytes, output...)
	c.TotalBytesRead += int64(n)
	return n, err
}

package test

import "fmt"

// ErrorReadCloser is a mock io.ReadCloser that always returns an error when reading.
type ErrorReadCloser struct{}

func (c *ErrorReadCloser) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("mock read error")
}

func (c *ErrorReadCloser) Close() error {
	return nil
}

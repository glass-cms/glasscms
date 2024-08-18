package test

import (
	"errors"
)

// ErrorReadCloser is a mock io.ReadCloser that always returns an error when reading.
type ErrorReadCloser struct{}

func (c *ErrorReadCloser) Read(_ []byte) (int, error) {
	return 0, errors.New("error reading")
}

func (c *ErrorReadCloser) Close() error {
	return nil
}

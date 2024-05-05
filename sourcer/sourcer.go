package sourcer

import "errors"

// ErrDone is returned when there are no items left in the data source.
var ErrDone = errors.New("no items left in the data source")

// DataSourcer is an iterator that provides data to be parsed.
type DataSourcer interface {
	// Next returns the next piece of data to be parsed.
	Next() (string, error)

	// Remaining returns the number of pieces of data remaining to be parsed.
	Remaining() int

	// Size returns the total number of pieces of data to be parsed.
	Size() int
}

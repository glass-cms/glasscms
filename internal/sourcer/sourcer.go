package sourcer

// DataSourcer is an iterator that provides data to be parsed.
type DataSourcer interface {
	// Next returns the next piece of data to be parsed.
	Next() (Source, error)

	// Remaining returns the number of pieces of data remaining to be parsed.
	Remaining() int

	// Size returns the total number of pieces of data to be parsed.
	Size() int
}

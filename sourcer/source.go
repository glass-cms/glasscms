package sourcer

import "io"

// Source is a data source that can be read from.
type Source interface {
	io.ReadCloser
	Name() string
}

// NilSource is a no-op sentinel Source.
var NilSource = nilSource{
	ReadCloser: io.NopCloser(nil),
}

type nilSource struct {
	io.ReadCloser
}

func (nilSource) Name() string {
	return ""
}

package fs

import (
	"os"
	"time"

	"github.com/djherbis/times"
)

type FileSource struct {
	*os.File

	birthtime time.Time
	modtime   time.Time

	sourceRootPath string
}

func NewFileSource(file *os.File, sourceRootPath string) (*FileSource, error) {
	stats, err := times.StatFile(file)
	if err != nil {
		return nil, err
	}

	return &FileSource{
		File:           file,
		birthtime:      stats.BirthTime(),
		modtime:        stats.ModTime(),
		sourceRootPath: sourceRootPath,
	}, nil
}

func (f FileSource) Name() string {
	return f.File.Name()
}

func (f FileSource) CreateTime() time.Time {
	return f.birthtime
}

func (f FileSource) UpdateTime() time.Time {
	return f.modtime
}

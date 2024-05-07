package sourcer

import (
	"os"
	"path/filepath"
	"strings"
)

var _ DataSourcer = &FileSystemSourcer{}

// FileSystemSourcer is a DataSourcer that reads files from the file system.
type FileSystemSourcer struct {
	files  []string
	cursor int // cursor is the index of the next file to be read
}

// NewFileSystemSourcer creates a new FileSystemSourcer.
func NewFileSystemSourcer(rootPath string) (*FileSystemSourcer, error) {
	var files []string

	err := filepath.WalkDir(rootPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.HasSuffix(d.Name(), ".md") {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &FileSystemSourcer{
		files: files,
	}, nil
}

func (s *FileSystemSourcer) Next() (Source, error) {
	if s.cursor >= len(s.files) {
		return NilSource, ErrDone
	}

	// Get a ReadCloser for the file
	file, err := os.Open(s.files[s.cursor])
	s.cursor++

	if err != nil {
		return NilSource, err
	}

	return NewFileSource(file)
}

func (s *FileSystemSourcer) Remaining() int {
	return s.Size() - s.cursor
}

func (s *FileSystemSourcer) Size() int {
	return len(s.files)
}

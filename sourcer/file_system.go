package sourcer

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var ErrInvalidFileSystemSource = errors.New("invalid file system source")

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

// IsValidFileSystemSource checks if the given path exist and is a directory.
func IsValidFileSystemSource(fp string) error {
	fileInfo, err := os.Stat(fp)
	if err != nil {
		// Wrap the error to provide more context.
		return fmt.Errorf("%w: %s", ErrInvalidFileSystemSource, err)
	}

	if !fileInfo.IsDir() {
		return fmt.Errorf("%w: %s is not a directory", ErrInvalidFileSystemSource, fp)
	}

	return nil
}

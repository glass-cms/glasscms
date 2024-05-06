package sourcer

var _ DataSourcer = &FileSystemSourcer{}

type FileSystemSourcer struct{}

func NewFileSystemSourcer() *FileSystemSourcer {
	return &FileSystemSourcer{}
}

func (s *FileSystemSourcer) Next() (string, error) {
	return "", nil
}

func (s *FileSystemSourcer) Remaining() int {
	return 0
}

func (s *FileSystemSourcer) Size() int {
	return 0
}

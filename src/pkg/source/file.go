package source

import "os"

var _ Reader = (*File)(nil)

func init() {
	Register("file", &File{})
}

type File struct {
	file *os.File
}

func (f *File) Init(source string) error {
	fh, err := os.Open(source)
	if err != nil {
		return err
	}
	f.file = fh
	return nil
}

func (f *File) Read(p []byte) (n int, err error) {
	return f.file.Read(p)
}

func (f *File) Seek(offset int64, whence int) (int64, error) {
	return f.file.Seek(offset, whence)
}

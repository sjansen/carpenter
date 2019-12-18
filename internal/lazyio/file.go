package lazyio

import (
	"io"
	"os"
	"path/filepath"
)

var _ Opener = &File{}

type File struct {
	Dir string
}

func (f *File) Open(path string) (io.WriteCloser, error) {
	path = filepath.Join(f.Dir, filepath.FromSlash(path))
	err := os.MkdirAll(filepath.Dir(path), 0777)
	if err != nil {
		return nil, err
	}
	return os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0666)
}

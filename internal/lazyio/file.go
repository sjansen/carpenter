package lazyio

import (
	"io"
	"os"
	"path/filepath"
)

var _ Opener = &FileOpener{}

type FileOpener struct {
	Dir string
}

func (o *FileOpener) Open(path string) (io.WriteCloser, error) {
	path = filepath.Join(o.Dir, filepath.FromSlash(path))
	err := os.MkdirAll(filepath.Dir(path), 0777)
	if err != nil {
		return nil, err
	}
	return os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0666)
}

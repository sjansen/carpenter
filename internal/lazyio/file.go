package lazyio

import (
	"io"
	"os"
	"path/filepath"
)

var _ InputOpener = &FileReader{}
var _ InputWalker = &FileReader{}
var _ OutputOpener = &FileWriter{}

type FileReader struct {
	Dir string
}

func (fr *FileReader) Open(path string) (io.ReadCloser, error) {
	path = filepath.Join(fr.Dir, filepath.FromSlash(path))
	return os.Open(path)
}

func (fr *FileReader) Walk(fn func(string) error) error {
	return filepath.Walk(fr.Dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		} else if !info.Mode().IsRegular() {
			return nil
		}
		suffix := fr.stripPrefix(path)
		return fn(suffix)
	})
}

func (fr *FileReader) stripPrefix(path string) string {
	n := len(fr.Dir)
	if fr.Dir[n-1] == '/' {
		path = path[n:]
	} else {
		path = path[n+1:]
	}
	return path
}

type FileWriter struct {
	Dir string
}

func (fw *FileWriter) Open(path string) (io.WriteCloser, error) {
	path = filepath.Join(fw.Dir, filepath.FromSlash(path))
	err := os.MkdirAll(filepath.Dir(path), 0777)
	if err != nil {
		return nil, err
	}
	return os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0666)
}

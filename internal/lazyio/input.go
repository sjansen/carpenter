package lazyio

import (
	"compress/gzip"
	"io"
	pathlib "path"
	"strings"
)

type Input struct {
	Path   string
	Opener InputOpener

	reader     io.ReadCloser
	compressed io.ReadCloser
}

func (i *Input) Close() error {
	var err error
	if i.reader != nil {
		err = i.reader.Close()
	}
	if i.compressed != nil {
		err := i.compressed.Close()
		if err != nil {
			return err
		}
	}
	return err
}

func (i *Input) Open() (io.Reader, error) {
	if i.reader == nil {
		r, err := i.Opener.Open(i.Path)
		if err != nil {
			return nil, err
		}
		if strings.HasSuffix(i.Path, ".gz") {
			i.compressed = r
			r, err = gzip.NewReader(r)
			if err != nil {
				return nil, err
			}
		}
		i.reader = r
	}
	return i.reader, nil
}

func (i *Input) StripExt() string {
	s := i.Path
	if strings.HasSuffix(s, ".gz") {
		s = s[:len(s)-3]
	}
	if n := len(pathlib.Ext(s)); n > 0 {
		s = s[:len(s)-n]
	}
	return s
}

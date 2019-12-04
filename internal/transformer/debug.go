package transformer

import (
	"encoding/csv"
	"io"
	"os"
	"path/filepath"
)

type debugFiles struct {
	normalize    *debugFile
	parse        *debugFile
	tokenize     *debugFile
	unrecognized *debugFile
}

func newDebugFiles(prefix, suffix string) *debugFiles {
	if prefix == "" || suffix == "" {
		return nil
	}
	return &debugFiles{
		normalize: &debugFile{
			Open: func() (io.WriteCloser, error) {
				os.MkdirAll(filepath.Join(prefix, "normalize"), 0777)
				filename := filepath.Join(prefix, "normalize", suffix+".csv")
				return os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666)
			},
		},
		parse: &debugFile{
			Open: func() (io.WriteCloser, error) {
				os.MkdirAll(filepath.Join(prefix, "parse"), 0777)
				filename := filepath.Join(prefix, "parse", suffix+".csv")
				return os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666)
			},
		},
		tokenize: &debugFile{
			Open: func() (io.WriteCloser, error) {
				os.MkdirAll(filepath.Join(prefix, "tokenize"), 0777)
				filename := filepath.Join(prefix, "tokenize", suffix+".csv")
				return os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666)
			},
		},
		unrecognized: &debugFile{
			Open: func() (io.WriteCloser, error) {
				os.MkdirAll(filepath.Join(prefix, "unrecognized"), 0777)
				filename := filepath.Join(prefix, "unrecognized", suffix+".csv")
				return os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666)
			},
		},
	}
}

func (d *debugFiles) Close() {
	if d != nil {
		d.normalize.Close()
		d.normalize = nil

		d.parse.Close()
		d.parse = nil

		d.tokenize.Close()
		d.tokenize = nil

		d.unrecognized.Close()
		d.unrecognized = nil
	}
}

type debugFile struct {
	Open func() (io.WriteCloser, error)

	w   io.WriteCloser
	csv *csv.Writer
	err error
}

func (f *debugFile) Close() {
	if f != nil && f.csv != nil {
		f.csv.Flush()
		f.w.Close()
	}
}

func (f *debugFile) Write(row ...string) error {
	switch {
	case f.err != nil:
		return f.err
	case f == nil:
		return nil
	case f.csv == nil:
		if w, err := f.Open(); err == nil {
			f.w = w
			f.csv = csv.NewWriter(w)
		} else {
			f.err = err
			return err
		}
	}

	return f.csv.Write(row)
}

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
			CSV: true,
			Open: func() (io.WriteCloser, error) {
				filename := filepath.Join(prefix, "normalize", suffix+".csv")
				os.MkdirAll(filepath.Dir(filename), 0777)
				return os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666)
			},
		},
		parse: &debugFile{
			CSV: true,
			Open: func() (io.WriteCloser, error) {
				filename := filepath.Join(prefix, "parse", suffix+".csv")
				os.MkdirAll(filepath.Dir(filename), 0777)
				return os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666)
			},
		},
		tokenize: &debugFile{
			CSV: false,
			Open: func() (io.WriteCloser, error) {
				filename := filepath.Join(prefix, "tokenize", suffix+".log")
				os.MkdirAll(filepath.Dir(filename), 0777)
				return os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666)
			},
		},
		unrecognized: &debugFile{
			CSV: true,
			Open: func() (io.WriteCloser, error) {
				filename := filepath.Join(prefix, "unrecognized", suffix+".csv")
				os.MkdirAll(filepath.Dir(filename), 0777)
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
	CSV  bool
	Open func() (io.WriteCloser, error)

	w   io.WriteCloser
	csv *csv.Writer
	err error
}

func (f *debugFile) Close() {
	if f == nil || f.w == nil {
		return
	}
	if f.csv != nil {
		f.csv.Flush()
	}
	f.w.Close()
}

func (f *debugFile) Write(row ...string) error {
	switch {
	case f == nil:
		return nil
	case f.err != nil:
		return f.err
	case f.w == nil:
		w, err := f.Open()
		if err != nil {
			f.err = err
			return err
		}
		f.w = w
		if f.CSV {
			f.csv = csv.NewWriter(w)
		}
	}

	if f.CSV {
		return f.csv.Write(row)
	}

	for _, s := range row {
		_, err := f.w.Write([]byte(s))
		if err != nil {
			return err
		}
	}

	return nil
}

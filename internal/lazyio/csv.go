package lazyio

import (
	"encoding/csv"
	"io"
)

type CSV struct {
	Path   string
	Opener OutputOpener

	w   io.WriteCloser
	csv *csv.Writer
	err error
}

func (f *CSV) Close() error {
	if f == nil || f.w == nil {
		return nil
	}
	if f.csv != nil {
		f.csv.Flush()
	}
	return f.w.Close()
}

func (f *CSV) Error() error {
	if f == nil || f.w == nil {
		return nil
	}
	if f.csv != nil {
		return f.csv.Error()
	}
	return nil
}

func (f *CSV) Flush() {
	if f != nil && f.csv != nil {
		f.csv.Flush()
	}
}

func (f *CSV) Write(row ...string) error {
	switch {
	case f == nil || f.Opener == nil:
		return nil
	case f.err != nil:
		return f.err
	case f.w == nil:
		w, err := f.Opener.Open(f.Path)
		if err != nil {
			f.err = err
			return err
		}
		f.w = w
		f.csv = csv.NewWriter(w)
	}

	return f.csv.Write(row)
}

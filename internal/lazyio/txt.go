package lazyio

import (
	"io"
)

type TXT struct {
	Path   string
	Opener OutputOpener

	w   io.WriteCloser
	err error
}

func (f *TXT) Close() error {
	if f == nil || f.w == nil {
		return nil
	}
	return f.w.Close()
}

func (f *TXT) Write(line string) error {
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
	}

	_, err := f.w.Write([]byte(line))
	if err != nil {
		return err
	}

	_, err = f.w.Write([]byte{'\n'})
	return err
}

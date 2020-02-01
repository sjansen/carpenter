package lazyio

import "io"

type InputOpener interface {
	Open(path string) (io.ReadCloser, error)
	StripPrefix(path string) string
}

type OutputOpener interface {
	Open(path string) (io.WriteCloser, error)
}

package lazyio

import "io"

type InputOpener interface {
	Open(path string) (io.ReadCloser, error)
}

type InputWalker interface {
	Walk(fn func(path string) error) error
}

type OutputOpener interface {
	Open(path string) (io.WriteCloser, error)
}

package lazyio

import "io"

type Opener interface {
	Open(path string) (io.WriteCloser, error)
}

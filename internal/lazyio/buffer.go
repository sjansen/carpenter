package lazyio

import (
	"bytes"
	"io"
)

var _ io.WriteCloser = &buffer{}

type buffer struct {
	bytes.Buffer
}

func (b *buffer) Close() error {
	return nil
}

var _ OutputOpener = &BufferWriter{}

type BufferWriter struct {
	buffers map[string]*buffer
}

func (b *BufferWriter) Open(path string) (io.WriteCloser, error) {
	if b.buffers == nil {
		b.buffers = make(map[string]*buffer, 1)
	}
	buf := &buffer{}
	b.buffers[path] = buf
	return buf, nil
}

func (b *BufferWriter) Buffer(path string) *bytes.Buffer {
	buf, ok := b.buffers[path]
	if ok {
		return &buf.Buffer
	}
	return nil
}

func (b *BufferWriter) Buffers() []string {
	result := make([]string, 0, len(b.buffers))
	for k := range b.buffers {
		result = append(result, k)
	}
	return result
}

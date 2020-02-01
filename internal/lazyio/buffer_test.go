package lazyio_test

import (
	"testing"

	"github.com/sjansen/carpenter/internal/lazyio"
	"github.com/stretchr/testify/require"
)

func TestBufferOpener(t *testing.T) {
	require := require.New(t)

	o := lazyio.BufferWriter{}

	w, err := o.Open("foo")
	require.NoError(err, "open")

	require.Equal([]string{"foo"}, o.Buffers())

	_, err = w.Write([]byte("bar baz"))
	require.NoError(err, "write")

	err = w.Close()
	require.NoError(err, "close")

	buf := o.Buffer("foo")
	require.NotNil(buf)
	require.Equal("bar baz", buf.String())
}

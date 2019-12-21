package lazyio_test

import (
	"testing"

	"github.com/sjansen/carpenter/internal/lazyio"
	"github.com/stretchr/testify/require"
)

func TestTXT(t *testing.T) {
	require := require.New(t)

	var txt *lazyio.TXT
	err := txt.Write("I want my $2!")
	require.NoError(err)
	err = txt.Close()
	require.NoError(err)

	o := &lazyio.BufferOpener{}
	expectedPath := "better/off/dead.txt"
	txt = &lazyio.TXT{
		Path:   expectedPath,
		Opener: o,
	}

	err = txt.Write("I want my $2!")
	require.NoError(err)
	require.Equal([]string{expectedPath}, o.Buffers())

	err = txt.Close()
	require.NoError(err)
	buf := o.Buffer(expectedPath)
	require.Equal("I want my $2!\n", buf.String())
}

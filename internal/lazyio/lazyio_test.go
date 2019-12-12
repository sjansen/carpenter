package lazyio_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/sjansen/carpenter/internal/lazyio"
	"github.com/stretchr/testify/require"
)

type buffer struct {
	bytes.Buffer

	path string
}

func (b *buffer) Close() error {
	b.path = ""
	return nil
}

func (b *buffer) Open(path string) (io.WriteCloser, error) {
	b.path = path
	b.Buffer.Reset()
	return b, nil
}

func TestCSV(t *testing.T) {
	require := require.New(t)

	var csv *lazyio.CSV
	err := csv.Write("a", "b", "c")
	require.NoError(err)
	err = csv.Close()
	require.NoError(err)

	buf := &buffer{}
	expectedPath := "space/magic.csv"
	csv = &lazyio.CSV{
		Path:   expectedPath,
		Opener: buf,
	}
	require.Equal("", buf.path)
	require.Equal("", buf.String())

	err = csv.Write("a", "b", "c")
	require.NoError(err)
	require.Equal(expectedPath, buf.path)

	err = csv.Close()
	require.NoError(err)
	require.Equal("", buf.path)
	require.Equal("a,b,c\n", buf.String())
}

func TestTXT(t *testing.T) {
	require := require.New(t)

	var txt *lazyio.TXT
	err := txt.Write("I want my $2!")
	require.NoError(err)
	err = txt.Close()
	require.NoError(err)

	buf := &buffer{}
	expectedPath := "better/off/dead.txt"
	txt = &lazyio.TXT{
		Path:   expectedPath,
		Opener: buf,
	}
	require.Equal("", buf.path)
	require.Equal("", buf.String())

	err = txt.Write("I want my $2!")
	require.NoError(err)
	require.Equal(expectedPath, buf.path)

	err = txt.Close()
	require.NoError(err)
	require.Equal("", buf.path)
	require.Equal("I want my $2!\n", buf.String())
}

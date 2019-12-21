package lazyio_test

import (
	"testing"

	"github.com/sjansen/carpenter/internal/lazyio"
	"github.com/stretchr/testify/require"
)

func TestCSV(t *testing.T) {
	require := require.New(t)

	var csv *lazyio.CSV
	err := csv.Write("a", "b", "c")
	require.NoError(err)
	err = csv.Close()
	require.NoError(err)

	o := &lazyio.BufferOpener{}
	expectedPath := "space/magic.csv"
	csv = &lazyio.CSV{
		Path:   expectedPath,
		Opener: o,
	}

	err = csv.Write("a", "b", "c")
	require.NoError(err)
	require.Equal([]string{expectedPath}, o.Buffers())

	err = csv.Close()
	require.NoError(err)
	buf := o.Buffer(expectedPath)
	require.Equal("a,b,c\n", buf.String())
}

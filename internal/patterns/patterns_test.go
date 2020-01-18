package patterns

import (
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMatch(t *testing.T) {
	require := require.New(t)

	for _, filename := range []string{
		"testdata/basic.star",
		"testdata/regex.star",
		"testdata/query.star",
	} {
		r, err := os.Open(filename)
		require.NoError(err)

		patterns, err := Load(filename, r)
		require.NoError(err)

		for raw, expected := range patterns.tests {
			url, err := url.Parse(raw)
			require.NoError(err, raw)

			id, normalized, err := patterns.Match(url)
			require.NoError(err, raw)

			if expected.url == "" {
				require.Equal("", id, raw)
			} else {
				require.Equal(expected.id, id, raw)
			}
			require.Equal(expected.url, normalized, raw)
		}
	}
}

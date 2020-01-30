package s3util_test

import (
	"testing"

	"github.com/sjansen/carpenter/internal/s3util"
	"github.com/stretchr/testify/require"
)

func TestParseURI(t *testing.T) {
	require := require.New(t)

	parsed, err := s3util.ParseURI("https://bucket.s3-us-west-2.amazonaws.com/key")
	require.Error(err)
	require.Nil(parsed)

	for _, tc := range []struct {
		uri      string
		expected *s3util.URI
	}{{
		uri: "s3://bucket/key",
		expected: &s3util.URI{
			Bucket: "bucket",
			Key:    "key",
		},
	}, {
		uri: "s3://a/b/c",
		expected: &s3util.URI{
			Bucket: "a",
			Key:    "b/c",
		},
	}, {
		uri: "s3://example.com",
		expected: &s3util.URI{
			Bucket: "example.com",
			Key:    "",
		},
	}} {
		actual, err := s3util.ParseURI(tc.uri)
		require.NoError(err)
		require.Equal(tc.expected, actual)
	}
}

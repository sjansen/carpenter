package s3util_test

import (
	"testing"

	"github.com/sjansen/carpenter/internal/s3util"
	"github.com/stretchr/testify/require"
)

func TestParseURI(t *testing.T) {
	require := require.New(t)

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
		uri: "s3://a/b/c?profile=d&endpoint=f",
		expected: &s3util.URI{
			Bucket:   "a",
			Key:      "b/c",
			Profile:  "d",
			Endpoint: "f",
		},
	}, {
		uri: "s3://example.com?region",
		expected: &s3util.URI{
			Bucket: "example.com",
			Key:    "",
		},
	}} {
		actual, err := s3util.ParseURI(tc.uri)
		require.NoError(err)
		require.Equal(tc.expected, actual)
	}

	parsed, err := s3util.ParseURI("https://bucket.s3-us-west-2.amazonaws.com/key")
	require.Error(err)
	require.Nil(parsed)

	parsed, err = s3util.ParseURI("s3://bucket/key?param=invalid")
	require.Error(err)
	require.Nil(parsed)
}

func TestURI_ToUploaderConfig(t *testing.T) {
	require := require.New(t)

	uri := &s3util.URI{
		Profile:  "profile",
		Region:   "region",
		Bucket:   "bucket",
		Key:      "prefix",
		Endpoint: "endpoint",
	}
	expected := &s3util.Config{
		Profile:  "profile",
		Region:   "region",
		Bucket:   "bucket",
		Prefix:   "prefix",
		Endpoint: "endpoint",
	}

	actual := uri.ToConfig()
	require.Equal(expected, actual)
}

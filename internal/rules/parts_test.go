package rules

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPlainPart(t *testing.T) {
	require := require.New(t)

	p := &PlainPart{
		id:    "id1",
		slash: "always",
		value: "foo",
	}

	id, match := p.Match("foo")
	require.Equal(true, match)
	require.Equal("id1", id)

	id, match = p.Match("bar")
	require.Equal(false, match)
	require.Equal("", id)

	require.Equal("foo", p.Rewrite("foo"))
}

func TestRegexpPart(t *testing.T) {
	require := require.New(t)

	for _, tc := range []struct {
		part        *RegexPart
		id          string
		match       string
		replacement string
	}{{
		&RegexPart{
			id:    "id1",
			slash: "never",
			regex: regexp.MustCompile("foo"),
			replacement: &PlainReplacement{
				value: "bar",
			},
		},
		"id1",
		"foo",
		"bar",
	}, {
		&RegexPart{
			id:          "id2",
			slash:       "never",
			regex:       regexp.MustCompile("qux"),
			replacement: &FunctionReplacement{},
		},
		"id2",
		"qux",
		"qux", // TODO
	}} {
		p := tc.part

		id, match := p.Match(tc.match)
		require.Equal(true, match)
		require.Equal(tc.id, id)

		id, match = p.Match("Spoon!")
		require.Equal(false, match)
		require.Equal("", id)

		require.Equal(tc.replacement, p.Rewrite(tc.match))
	}
}

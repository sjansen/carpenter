package rules

import (
	"errors"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
	"go.starlark.net/starlark"
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

	actual, err := p.Rewrite("foo")
	require.NoError(err)
	require.Equal("foo", actual)
}

func TestRegexpPart(t *testing.T) {
	require := require.New(t)

	for _, tc := range []struct {
		id          string
		match       string
		replacement string
		error       string
		part        *RegexPart
	}{{
		"id1", "foo", "bar", "",
		&RegexPart{
			id:    "id1",
			slash: "always",
			regex: regexp.MustCompile("^foo$"),
			replacement: &PlainReplacement{
				value: "bar",
			},
		},
	}, {
		"id2", "qux", "quux", "",
		&RegexPart{
			id:    "id2",
			slash: "never",
			regex: regexp.MustCompile("^qux$"),
			replacement: &FunctionReplacement{
				thread: &starlark.Thread{},
				fn: starlark.NewBuiltin("rewrite", func(
					t *starlark.Thread,
					fn *starlark.Builtin,
					args starlark.Tuple,
					kwargs []starlark.Tuple,
				) (starlark.Value, error) {
					return starlark.String("quux"), nil
				}),
			},
		},
	}, {
		"id3", "corge", "", "expected String, got int",
		&RegexPart{
			id:    "id3",
			slash: "strip",
			regex: regexp.MustCompile("^c.*$"),
			replacement: &FunctionReplacement{
				thread: &starlark.Thread{},
				fn: starlark.NewBuiltin("rewrite", func(
					t *starlark.Thread,
					fn *starlark.Builtin,
					args starlark.Tuple,
					kwargs []starlark.Tuple,
				) (starlark.Value, error) {
					return starlark.MakeInt(42), nil
				}),
			},
		},
	}, {
		"id4", "grault", "", "418 I'm a teapot",
		&RegexPart{
			id:    "id4",
			slash: "strip",
			regex: regexp.MustCompile("^g.....$"),
			replacement: &FunctionReplacement{
				thread: &starlark.Thread{},
				fn: starlark.NewBuiltin("rewrite", func(
					t *starlark.Thread,
					fn *starlark.Builtin,
					args starlark.Tuple,
					kwargs []starlark.Tuple,
				) (starlark.Value, error) {
					return nil, errors.New("418 I'm a teapot")
				}),
			},
		},
	}} {
		p := tc.part

		id, match := p.Match("should not match")
		require.Equal("", id, tc.id)
		require.Equal(false, match, tc.id)

		id, match = p.Match(tc.match)
		require.Equal(tc.id, id, tc.id)
		require.Equal(true, match, tc.id)

		actual, err := p.Rewrite(tc.match)
		if tc.error != "" {
			require.Equal("", actual, tc.id)
			require.Contains(err.Error(), tc.error, tc.id)
		} else {
			require.Equal(tc.replacement, actual, tc.id)
			require.NoError(err, tc.id)
		}
	}
}

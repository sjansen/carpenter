package patterns

import (
	"errors"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
	"go.starlark.net/starlark"
)

func TestPlainPart(t *testing.T) {
	require := require.New(t)

	p := &plainPart{
		value: "foo",
	}

	require.False(p.greedy())

	match := p.match("foo")
	require.Equal(true, match)

	match = p.match("bar")
	require.Equal(false, match)

	actual, err := p.normalize(&starlark.Thread{}, "foo")
	require.NoError(err)
	require.Equal("foo", actual)
}

func TestRegexpPart(t *testing.T) {
	require := require.New(t)

	for _, tc := range []struct {
		id       string
		match    string
		expected string
		error    string
		part     *regexPart
	}{{
		"id1", "foo", "bar", "",
		&regexPart{
			regex: regexp.MustCompile("^foo$"),
			rewriter: &staticStringRewriter{
				value: "bar",
			},
		},
	}, {
		"id2", "qux", "quux", "",
		&regexPart{
			regex: regexp.MustCompile("^qux$"),
			rewriter: &callableStringRewriter{
				Callable: starlark.NewBuiltin("rewrite", func(
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
		&regexPart{
			regex: regexp.MustCompile("^c.*$"),
			rewriter: &callableStringRewriter{
				Callable: starlark.NewBuiltin("rewrite", func(
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
		&regexPart{
			regex: regexp.MustCompile("^g.....$"),
			rewriter: &callableStringRewriter{
				Callable: starlark.NewBuiltin("rewrite", func(
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

		p.suffix = true
		require.True(p.greedy())
		p.suffix = false
		require.False(p.greedy())

		match := p.match("should not match")
		require.Equal(false, match, tc.id)

		match = p.match(tc.match)
		require.Equal(true, match, tc.id)

		actual, err := p.normalize(&starlark.Thread{}, tc.match)
		if tc.error != "" {
			require.Contains(err.Error(), tc.error, tc.id)
		} else {
			require.NoError(err, tc.id)
		}
		require.Equal(tc.expected, actual, tc.id)
	}
}

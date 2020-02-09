package patterns

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.starlark.net/starlark"
)

var echo = starlark.NewBuiltin("echo", func(
	t *starlark.Thread,
	fn *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {
	return args.Index(1), nil
})

var fail = starlark.NewBuiltin("rewrite", func(
	t *starlark.Thread,
	fn *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {
	return nil, errors.New("i am error")
})

func TestParam(t *testing.T) {
	require := require.New(t)

	for _, tc := range []struct {
		id       string
		dedup    dedup
		error    string
		expected []string
		values   []string
		param    *param
	}{{
		"p1", keepAll, "",
		[]string{},
		[]string{},
		&param{
			remove: false,
			rewriter: &staticQueryRewriter{
				value: "X",
			},
		},
	}, {
		"p2", keepAll, "",
		nil,
		[]string{"foo", "bar", "baz"},
		&param{
			remove: true,
			rewriter: &staticQueryRewriter{
				value: "X",
			},
		},
	}, {
		"p3", keepAll, "",
		[]string{"X", "X", "X"},
		[]string{"foo", "bar", "baz"},
		&param{
			remove: false,
			rewriter: &staticQueryRewriter{
				value: "X",
			},
		},
	}, {
		"p4", keepFirst, "",
		[]string{"foo"},
		[]string{"foo", "bar", "baz"},
		&param{
			remove: false,
			rewriter: &callableQueryRewriter{
				Callable: echo,
			},
		},
	}, {
		"p5", keepLast, "",
		[]string{"baz"},
		[]string{"foo", "bar", "baz"},
		&param{
			remove: false,
			rewriter: &callableQueryRewriter{
				Callable: echo,
			},
		},
	}, {
		"p6", keepFirst, "i am error",
		nil,
		[]string{"foo", "bar", "baz"},
		&param{
			remove: false,
			rewriter: &callableQueryRewriter{
				Callable: fail,
			},
		},
	}, {
		"p7", keepLast, "i am error",
		nil,
		[]string{"foo", "bar", "baz"},
		&param{
			remove: false,
			rewriter: &callableQueryRewriter{
				Callable: fail,
			},
		},
	}, {
		"p8", keepAll, "i am error",
		nil,
		[]string{"foo", "bar", "baz"},
		&param{
			remove: false,
			rewriter: &callableQueryRewriter{
				Callable: fail,
			},
		},
	}} {
		p := tc.param

		actual, err := p.normalize(&starlark.Thread{}, tc.dedup, tc.id, tc.values)
		if tc.error != "" {
			require.Contains(err.Error(), tc.error, tc.id)
		} else {
			require.NoError(err, tc.id)
		}
		require.Equal(tc.expected, actual, tc.id)
	}
}

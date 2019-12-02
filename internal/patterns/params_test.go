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
	return args.Index(0), nil
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
		dedup    string
		error    string
		expected []string
		values   []string
		param    *param
	}{{
		"p1", "never", "",
		[]string{},
		[]string{},
		&param{
			remove: false,
			rewriter: &rewriteStatic{
				value: "X",
			},
		},
	}, {
		"p2", "never", "",
		nil,
		[]string{"foo", "bar", "baz"},
		&param{
			remove: true,
			rewriter: &rewriteStatic{
				value: "X",
			},
		},
	}, {
		"p3", "never", "",
		[]string{"X", "X", "X"},
		[]string{"foo", "bar", "baz"},
		&param{
			remove: false,
			rewriter: &rewriteStatic{
				value: "X",
			},
		},
	}, {
		"p4", "first", "",
		[]string{"foo"},
		[]string{"foo", "bar", "baz"},
		&param{
			remove: false,
			rewriter: &rewriteFunction{
				Callable: echo,
			},
		},
	}, {
		"p5", "last", "",
		[]string{"baz"},
		[]string{"foo", "bar", "baz"},
		&param{
			remove: false,
			rewriter: &rewriteFunction{
				Callable: echo,
			},
		},
	}, {
		"p6", "first", "i am error",
		nil,
		[]string{"foo", "bar", "baz"},
		&param{
			remove: false,
			rewriter: &rewriteFunction{
				Callable: fail,
			},
		},
	}, {
		"p7", "last", "i am error",
		nil,
		[]string{"foo", "bar", "baz"},
		&param{
			remove: false,
			rewriter: &rewriteFunction{
				Callable: fail,
			},
		},
	}, {
		"p8", "never", "i am error",
		nil,
		[]string{"foo", "bar", "baz"},
		&param{
			remove: false,
			rewriter: &rewriteFunction{
				Callable: fail,
			},
		},
	}} {
		p := tc.param

		actual, err := p.normalize(&starlark.Thread{}, tc.dedup, tc.values)
		if tc.error != "" {
			require.Contains(err.Error(), tc.error, tc.id)
		} else {
			require.NoError(err, tc.id)
		}
		require.Equal(tc.expected, actual, tc.id)
	}
}

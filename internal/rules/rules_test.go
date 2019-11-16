package rules

import (
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var script = `
register_urls({
    'id': 'views.foo',
    'parts': ['foo'],
    'slash': 'always',
}, {
    'id': 'views.bar',
    'parts': ['bar'],
    'slash': 'never',
}, {
    'id': 'views.baz',
    'parts': ['baz'],
    'slash': 'strip',
}, {
    'id': 'views.qux',
    'parts': [('qux', 'quux')],
    'slash': 'always',
}, {
    'id': 'views.multi',
    'parts': ['corge', ('grault', 'garply')],
    'slash': 'always',
})
`
var expected = Rules{
	{
		id:    "views.foo",
		slash: "always",
		parts: []matcher{
			&plainPart{
				value: "foo",
			},
		},
	}, {
		id:    "views.bar",
		slash: "never",
		parts: []matcher{
			&plainPart{
				value: "bar",
			},
		},
	}, {
		id:    "views.baz",
		slash: "strip",
		parts: []matcher{
			&plainPart{
				value: "baz",
			},
		},
	}, {
		id:    "views.qux",
		slash: "always",
		parts: []matcher{
			&regexPart{
				regex: regexp.MustCompile("qux"),
				replacement: &plainReplacement{
					value: "quux",
				},
			},
		},
	}, {
		id:    "views.multi",
		slash: "always",
		parts: []matcher{
			&plainPart{
				value: "corge",
			},
			&regexPart{
				regex: regexp.MustCompile("grault"),
				replacement: &plainReplacement{
					value: "garply",
				},
			},
		},
	},
}

func TestLoad(t *testing.T) {
	require := require.New(t)

	r := strings.NewReader(script)
	actual, err := Load("<test script>", r)
	require.NoError(err)
	require.Equal(expected, actual)
}

func TestLoadErrors(t *testing.T) {
	require := require.New(t)

	for _, tc := range []struct {
		script  string
		message string
	}{{
		script: "Spoon!",
	}, {
		script:  `register_urls(foo=42)`,
		message: "unexpected keyword argument",
	}, {
		script:  `register_urls(42)`,
		message: "expected Dict, got int",
	}, {
		script:  `register_urls({})`,
		message: `missing required key: "id"`,
	}, {
		script:  `register_urls({'id': []})`,
		message: "expected String, got list",
	}, {
		script:  `register_urls({'id': 'foo'})`,
		message: `missing required key: "parts"`,
	}, {
		script:  `register_urls({'id': 'foo', 'parts': None})`,
		message: "expected Iterable, got NoneType",
	}, {
		script:  `register_urls({'id': 'foo', 'parts': ['bar']})`,
		message: `missing required key: "slash"`,
	}, {
		script:  `register_urls({'id': 'foo', 'parts': ['bar'], 'slash': None})`,
		message: "expected String, got NoneType",
	}, {
		script:  `register_urls({'id': 'foo', 'parts': [42], 'slash': 'strip'})`,
		message: "expected String or Tuple, got int",
	}} {
		r := strings.NewReader(tc.script)
		_, err := Load("<test script>", r)
		require.Error(err)
		require.Contains(err.Error(), tc.message)
	}
}

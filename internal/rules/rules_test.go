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
})
`
var expected = &Rules{Matchers: []Matcher{
	&PlainPart{
		id:    "views.foo",
		slash: "always",
		value: "foo",
	},
	&PlainPart{
		id:    "views.bar",
		slash: "never",
		value: "bar",
	},
	&PlainPart{
		id:    "views.baz",
		slash: "strip",
		value: "baz",
	},
	&RegexPart{
		id:    "views.qux",
		slash: "always",
		regex: regexp.MustCompile("qux"),
		replacement: &PlainReplacement{
			value: "quux",
		},
	},
}}

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

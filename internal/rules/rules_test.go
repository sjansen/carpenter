package rules

import (
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var script = `
names = ('waldo|fred', 'plugh')

def fn(x):
    if x == 'xyzzy':
	return 'thud'
    return 'X'

register_urls({
    'id': 'views.always',
    'parts': ['foo'],
    'slash': 'always',
}, {
    'id': 'views.never',
    'parts': ['bar'],
    'slash': 'never',
}, {
    'id': 'views.strip',
    'parts': ['baz'],
    'slash': 'strip',
}, {
    'id': 'views.regex',
    'parts': [('qux', 'quux')],
    'slash': 'always',
}, {
    'id': 'views.multi',
    'parts': [
	'corge',
	('grault', 'garply'),
	names,
	('.+', fn),
    ],
    'slash': 'always',
})
`
var expected = Rules{
	{
		id:    "views.always",
		slash: "always",
		parts: []part{
			&plainPart{
				value: "foo",
			},
		},
	}, {
		id:    "views.never",
		slash: "never",
		parts: []part{
			&plainPart{
				value: "bar",
			},
		},
	}, {
		id:    "views.strip",
		slash: "strip",
		parts: []part{
			&plainPart{
				value: "baz",
			},
		},
	}, {
		id:    "views.regex",
		slash: "always",
		parts: []part{
			&regexPart{
				regex: regexp.MustCompile("qux"),
				rewriter: &rewriteStatic{
					value: "quux",
				},
			},
		},
	}, {
		id:    "views.multi",
		slash: "always",
		parts: []part{
			&plainPart{
				value: "corge",
			},
			&regexPart{
				regex: regexp.MustCompile("grault"),
				rewriter: &rewriteStatic{
					value: "garply",
				},
			},
			&regexPart{
				regex: regexp.MustCompile("waldo|fred"),
				rewriter: &rewriteStatic{
					value: "plugh",
				},
			},
			&regexPart{
				regex:    regexp.MustCompile(".+"),
				rewriter: nil,
			},
		},
	},
}

func TestLoad(t *testing.T) {
	require := require.New(t)

	r := strings.NewReader(script)
	actual, err := Load("<test script>", r)
	require.NoError(err)

	// It isn't possible to create an expected rewriteFunction
	// that reflect.DeepEqual considers equal to actual.
	actualPart := actual[4].parts[3].(*regexPart)
	require.IsType(&rewriteFunction{}, actualPart.rewriter)
	expectedPart := expected[4].parts[3].(*regexPart)
	expectedPart.rewriter = actualPart.rewriter

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

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
    'tests': {
        '/foo':   None,
        '/foo/': '/foo/',
    },
}, {
    'id': 'views.never',
    'parts': ['bar'],
    'slash': 'never',
    'tests': {
        '/bar': '/bar',
        '/bar/': None,
    },
}, {
    'id': 'views.strip',
    'parts': ['baz'],
    'slash': 'strip',
    'tests': {
        '/baz':  '/baz',
        '/baz/': '/baz',
    },
}, {
    'id': 'views.regex',
    'parts': [('qux', 'quux')],
    'slash': 'always',
    'tests': {
        '/qux/': '/quux/',
    },
}, {
    'id': 'views.multi',
    'parts': [
	'corge',
	('grault', 'garply'),
	names,
	('.+', fn),
    ],
    'slash': 'always',
    'tests': {
        '/corge/grault/waldo/xyzzy/': '/corge/garply/plugh/thud/',
        '/corge/grault/fred/42/':     '/corge/garply/plugh/X/',
    },
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
		tests: map[string]string{
			"/foo":  "",
			"/foo/": "/foo/",
		},
	}, {
		id:    "views.never",
		slash: "never",
		parts: []part{
			&plainPart{
				value: "bar",
			},
		},
		tests: map[string]string{
			"/bar":  "/bar",
			"/bar/": "",
		},
	}, {
		id:    "views.strip",
		slash: "strip",
		parts: []part{
			&plainPart{
				value: "baz",
			},
		},
		tests: map[string]string{
			"/baz":  "/baz",
			"/baz/": "/baz",
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
		tests: map[string]string{
			"/qux/": "/quux/",
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
		tests: map[string]string{
			"/corge/grault/waldo/xyzzy/": "/corge/garply/plugh/thud/",
			"/corge/grault/fred/42/":     "/corge/garply/plugh/X/",
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
		script:  `register_urls({'id': 'parts-missing'})`,
		message: `"parts-missing" missing required key: "parts"`,
	}, {
		script:  `register_urls({'id': 'slash-missing', 'parts': []})`,
		message: `"slash-missing" missing required key: "slash"`,
	}, {
		script:  `register_urls({'id': 'tests-missing', 'parts': [], 'slash': 'always'})`,
		message: `"tests-missing" missing required key: "tests"`,
	}, {
		script: `register_urls({
		    'id': 'parts-invalid-1', 'parts': None, 'slash': 'always', 'tests': {},
		})`,
		message: `"parts-invalid-1" expected Iterable, got NoneType`,
	}, {
		script: `register_urls({
		    'id': 'parts-invalid-2', 'parts': [42], 'slash': 'always', 'tests': {},
		})`,
		message: `"parts-invalid-2" expected String or Tuple, got int`,
	}, {
		script: `register_urls({
		    'id': 'parts-invalid-3', 'parts': [(1, 2, 3)], 'slash': 'always', 'tests': {},
		})`,
		message: `"parts-invalid-3" expected 2 item Tuple, got 3`,
	}, {
		script: `register_urls({
		    'id': 'parts-invalid-4', 'parts': [(1, 2)], 'slash': 'always', 'tests': {},
		})`,
		message: `"parts-invalid-4" expected String, got int`,
	}, {
		script: `register_urls({
		    'id': 'parts-invalid-5', 'parts': [('answer', 42)], 'slash': 'always', 'tests': {},
		})`,
		message: `"parts-invalid-5" expected Callable or String, got int`,
	}, {
		script: `register_urls({
		    'id': 'parts-invalid-6', 'parts': [('[a-z', 'X')], 'slash': 'always', 'tests': {},
		})`,
		message: `error parsing regexp`,
	}, {
		script: `register_urls({
		    'id': 'slash-invalid-1', 'parts': [], 'slash': None, 'tests': {},
		})`,
		message: `"slash-invalid-1" expected String, got NoneType`,
	}, {
		script: `register_urls({
		    'id': 'tests-invalid-1', 'parts': [], 'slash': 'always', 'tests': None,
		})`,
		message: `"tests-invalid-1" expected Dict, got NoneType`,
	}, {
		script: `register_urls({
		    'id': 'tests-invalid-2', 'parts': [], 'slash': 'always', 'tests': {6: 9},
		})`,
		message: `"tests-invalid-2" expected String key, got int`,
	}, {
		script: `register_urls({
		    'id': 'tests-invalid-3', 'parts': [], 'slash': 'always', 'tests': {'answer': 42},
		})`,
		message: `"tests-invalid-3" expected None or String value, got int`,
	}} {
		r := strings.NewReader(tc.script)
		_, err := Load("<test script>", r)
		require.Error(err)
		require.Contains(err.Error(), tc.message)
	}
}

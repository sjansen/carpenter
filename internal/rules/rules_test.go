package rules

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var expected = Rules{
	{
		id:     "views.root",
		dedup:  "never",
		slash:  "always",
		parts:  []part{},
		params: map[string]*param{},
		tests: map[string]string{
			"/":       "/",
			"/Spoon!": "",
		},
	}, {
		id:    "views.always",
		dedup: "never",
		slash: "always",
		parts: []part{
			&plainPart{
				value: "foo",
			},
		},
		params: map[string]*param{},
		tests: map[string]string{
			"/foo":  "",
			"/foo/": "/foo/",
		},
	}, {
		id:    "views.never",
		dedup: "never",
		slash: "never",
		parts: []part{
			&plainPart{
				value: "bar",
			},
		},
		params: map[string]*param{},
		tests: map[string]string{
			"/bar":  "/bar",
			"/bar/": "",
		},
	}, {
		id:    "views.strip",
		dedup: "never",
		slash: "strip",
		parts: []part{
			&plainPart{
				value: "baz",
			},
		},
		params: map[string]*param{},
		tests: map[string]string{
			"/baz":  "/baz",
			"/baz/": "/baz",
		},
	}, {
		id:    "views.regex",
		dedup: "never",
		slash: "always",
		parts: []part{
			&regexPart{
				regex: regexp.MustCompile("qux"),
				rewriter: &rewriteStatic{
					value: "quux",
				},
			},
		},
		params: map[string]*param{},
		tests: map[string]string{
			"/qux/": "/quux/",
		},
	}, {
		id:    "query.never",
		dedup: "never",
		slash: "never",
		parts: []part{
			&plainPart{
				value: "search",
			},
		},
		params: map[string]*param{
			"q": {
				remove: false,
				rewriter: &rewriteStatic{
					value: "X",
				},
			},
			"utf8": {
				remove: true,
			},
		},
		tests: map[string]string{
			"/search?q=cats":        "/search?q=X",
			"/search?utf8=✔":        "/search",
			"/search?q=dogs&utf8=✔": "/search?q=X",
		},
	}, {
		id:    "views.multi",
		dedup: "never",
		slash: "strip",
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
		params: map[string]*param{
			"n": {
				remove:   false,
				rewriter: nil,
			},
		},
		tests: map[string]string{
			"/corge/grault/waldo/xyzzy":         "/corge/garply/plugh/thud",
			"/corge/grault/fred/42/":            "/corge/garply/plugh/X",
			"/corge/grault/fred/random/?n=left": "/corge/garply/plugh/X?n=even",
			"/corge/grault/fred/random?n=right": "/corge/garply/plugh/X?n=odd",
		},
	},
}

func TestLoad(t *testing.T) {
	require := require.New(t)

	r, err := os.Open("testdata/valid.star")
	require.NoError(err)

	actual, err := Load("<test script>", r)
	require.NoError(err)

	// It isn't possible to create an expected rewriteFunction
	// that reflect.DeepEqual considers equal to actual.
	last := 6
	require.Len(actual, last+1)
	require.Equal("views.multi", actual[last].id)

	actualPart := actual[last].parts[3].(*regexPart)
	require.IsType(&rewriteFunction{}, actualPart.rewriter)
	expectedPart := expected[last].parts[3].(*regexPart)
	expectedPart.rewriter = actualPart.rewriter

	actualParam := actual[last].params["n"]
	require.IsType(&rewriteFunction{}, actualParam.rewriter)
	expectedParam := expected[last].params["n"]
	expectedParam.rewriter = actualParam.rewriter

	require.Equal(expected, actual)
}

func TestLoadErrors(t *testing.T) {
	ws := regexp.MustCompile(`\s+`)

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
		script:  `register_urls({'id': 'path-missing'})`,
		message: `"path-missing" missing required key: "path"`,
	}, {
		script:  `register_urls({'id': 'path-invalid', 'path': None})`,
		message: `"path-invalid" expected Dict, got NoneType`,
	}, {
		script:  `register_urls({'id': 'parts-missing', 'path': {}})`,
		message: `"parts-missing" missing required key: "parts"`,
	}, {
		script: `register_urls({
		    'id': 'slash-missing',
		    'path': {'parts': []},
		})`,
		message: `"slash-missing" missing required key: "slash"`,
	}, {
		script: `register_urls({
		    'id': 'query-missing',
		    'path': {'parts': [], 'slash': 'always'},
		})`,
		message: `"query-missing" missing required key: "query"`,
	}, {
		script: `register_urls({
		    'id': 'tests-missing',
		    'path': {'parts': [], 'slash': 'always'},
		    'query': {},
		})`,
		message: `"tests-missing" missing required key: "tests"`,
	}, {
		script: `register_urls({
		    'id': 'parts-invalid-1',
		    'path': {'parts': None, 'slash': 'always'},
		    'query': {},
		    'tests': {},
		})`,
		message: `"parts-invalid-1" expected Iterable, got NoneType`,
	}, {
		script: `register_urls({
		    'id': 'parts-invalid-2',
		    'path': {'parts': [42], 'slash': 'always'},
		    'query': {},
		    'tests': {},
		})`,
		message: `"parts-invalid-2" expected String or Tuple, got int`,
	}, {
		script: `register_urls({
		    'id': 'parts-invalid-3',
		    'path': {'parts': [(1, 2, 3)], 'slash': 'always'},
		    'query': {},
		    'tests': {},
		})`,
		message: `"parts-invalid-3" expected 2 item Tuple, got 3`,
	}, {
		script: `register_urls({
		    'id': 'parts-invalid-4',
		    'path': {'parts': [(1, 2)], 'slash': 'always'},
		    'query': {},
		    'tests': {},
		})`,
		message: `"parts-invalid-4" expected String, got int`,
	}, {
		script: `register_urls({
		    'id': 'parts-invalid-5',
		    'path': {'parts': [('answer', 42)], 'slash': 'always'},
		    'query': {},
		    'tests': {},
		})`,
		message: `"parts-invalid-5" expected Callable or String, got int`,
	}, {
		script: `register_urls({
		    'id': 'parts-invalid-6',
		    'path': {'parts': [('[a-z', 'X')], 'slash': 'always'},
		    'query': {},
		    'tests': {},
		})`,
		message: `error parsing regexp`,
	}, {
		script: `register_urls({
		    'id': 'slash-invalid-1',
		    'path': {'parts': [], 'slash': None},
		    'query': {},
		    'tests': {},
		})`,
		message: `"slash-invalid-1" expected String, got NoneType`,
	}, {
		script: `register_urls({
		    'id': 'query-invalid-1',
		    'path': {'parts': [], 'slash': 'always'},
		    'query': None,
		    'tests': {},
		})`,
		message: `"query-invalid-1" expected Dict, got NoneType`,
	}, {
		script: `register_urls({
		    'id': 'query-invalid-2',
		    'path': {'parts': [], 'slash': 'always'},
		    'query': {42: 'answer'},
		    'tests': {},
		})`,
		message: `"query-invalid-2" expected String key, got int`,
	}, {
		script: `register_urls({
		    'id': 'query-invalid-3',
		    'path': {'parts': [], 'slash': 'always'},
		    'query': {"dedup": 0},
		    'tests': {},
		})`,
		message: `"query-invalid-3" expected String, got int`,
	}, {
		script: `register_urls({
		    'id': 'query-invalid-4',
		    'path': {'parts': [], 'slash': 'always'},
		    'query': {"keys": None},
		    'tests': {},
		})`,
		message: `"query-invalid-4" expected "dedup" or "params", got "keys"`,
	}, {
		script: `register_urls({
		    'id': 'tests-invalid-1',
		    'path': {'parts': [], 'slash': 'always'},
		    'query': {},
		    'tests': None,
		})`,
		message: `"tests-invalid-1" expected Dict, got NoneType`,
	}, {
		script: `register_urls({
		    'id': 'tests-invalid-2',
		    'path': {'parts': [], 'slash': 'always'},
		    'query': {},
		    'tests': {6: 9},
		})`,
		message: `"tests-invalid-2" expected String key, got int`,
	}, {
		script: `register_urls({
		    'id': 'tests-invalid-3',
		    'path': {'parts': [], 'slash': 'always'},
		    'query': {},
		    'tests': {'answer': 42},
		})`,
		message: `"tests-invalid-3" expected None or String value, got int`,
	}} {
		tc := tc
		t.Run(ws.ReplaceAllString(tc.script, " "), func(t *testing.T) {
			require := require.New(t)

			r := strings.NewReader(tc.script)
			_, err := Load("<test script>", r)
			require.Error(err)
			require.Contains(err.Error(), tc.message)
		})
	}
}

func TestSelfTest(t *testing.T) {
	require := require.New(t)

	r, err := os.Open("testdata/valid.star")
	require.NoError(err)

	rules, err := Load("<test script>", r)
	require.NoError(err)

	err = rules.SelfTest()
	require.NoError(err)
}

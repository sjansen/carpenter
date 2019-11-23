package rules

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var expected = Rules{
	{
		id:     "root",
		dedup:  "never",
		slash:  "strip",
		parts:  []part{},
		params: map[string]*param{},
		tests: map[string]string{
			"/":       "/",
			"/Spoon!": "",
		},
	}, {
		id:    "always",
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
		id:    "never",
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
		id:    "strip",
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
		id:    "regex",
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
		id:    "multi",
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

	r, err := os.Open("testdata/raw.star")
	require.NoError(err)

	actual, err := Load("<test script>", r)
	require.NoError(err)

	// It isn't possible to create an expected rewriteFunction
	// that reflect.DeepEqual considers equal to actual.
	last := 6
	require.Len(actual, last+1)
	require.Equal("multi", actual[last].id)

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
	files, _ := filepath.Glob("testdata/load/err-*.star")
	for _, tc := range files {
		tc := tc
		prefix := tc[:len(tc)-5]
		t.Run(filepath.Base(prefix), func(t *testing.T) {
			require := require.New(t)

			msg, err := ioutil.ReadFile(prefix + ".txt")
			require.NoError(err)
			expected := strings.TrimSpace(string(msg))

			r, err := os.Open(tc)
			require.NoError(err)

			_, err = Load("<test script>", r)
			require.Error(err)
			require.Equal(expected, err.Error())
		})
	}
}

func TestSelfTest(t *testing.T) {
	require := require.New(t)

	r, err := os.Open("testdata/raw.star")
	require.NoError(err)

	rules, err := Load("<test script>", r)
	require.NoError(err)

	err = rules.SelfTest()
	require.NoError(err)
}

func TestSelfTestErrors(t *testing.T) {
	files, _ := filepath.Glob("testdata/selftest/err-*.star")
	for _, tc := range files {
		tc := tc
		prefix := tc[:len(tc)-5]
		t.Run(filepath.Base(prefix), func(t *testing.T) {
			require := require.New(t)

			msg, err := ioutil.ReadFile(prefix + ".txt")
			require.NoError(err)
			expected := strings.TrimSpace(string(msg))

			r, err := os.Open(tc)
			require.NoError(err)

			rules, err := Load("<test script>", r)
			require.NoError(err)

			err = rules.SelfTest()
			require.Error(err)
			require.Equal(expected, err.Error())
			require.NotEmpty(expected)
		})
	}
}

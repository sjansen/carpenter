package patterns

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	require := require.New(t)

	for filename, tc := range map[string]struct {
		expected *Patterns
		// It's not possible to declare rewriters that reflect.DeepEqual
		// recognizes, so fixers copy rewriters from actual to expected.
		fixer func(t *testing.T, expected, actual *tree)
	}{
		"testdata/basic.star": {
			expected: basicTree,
		},
		"testdata/regex.star": {
			expected: regexTree,
			fixer:    regexTreeFixer,
		},
	} {
		r, err := os.Open(filename)
		require.NoError(err)

		actual, err := Load(filename, r)
		require.NoError(err)

		if tc.fixer != nil {
			tc.fixer(t, &tc.expected.tree, &actual.tree)
		}
		require.Equal(tc.expected, actual)
	}
}

func TestLoadErrors(t *testing.T) {
	files, _ := filepath.Glob("testdata/load-errors/*.star")
	for _, tc := range files {
		tc := tc
		prefix := tc[:len(tc)-5]
		t.Run(filepath.Base(prefix), func(t *testing.T) {
			require := require.New(t)

			r, err := os.Open(tc)
			require.NoError(err)

			_, err = Load(filepath.Base(tc), r)
			require.Error(err)
			actual := err.Error()

			msg, err := ioutil.ReadFile(prefix + ".err")
			if os.IsNotExist(err) {
				msg = []byte(actual)
				ioutil.WriteFile(prefix+".err", msg, 0664)
			} else {
				require.NoError(err)
				expected := strings.TrimSpace(string(msg))

				if expected != actual {
					msg = []byte(actual)
					ioutil.WriteFile(prefix+".actual", msg, 0664)
				}

				require.Equal(expected, actual)
			}
		})
	}
}

func TestLoadPatterns(t *testing.T) {
	require := require.New(t)

	for filename, tc := range map[string]struct {
		expected []*pattern
		// It's not possible to declare rewriters that reflect.DeepEqual
		// recognizes, so fixers copy rewriters from actual to expected.
		fixer func(t *testing.T, expected, actual []*pattern)
	}{
		"testdata/basic.star": {
			expected: basicPatterns,
		},
		"testdata/regex.star": {
			expected: regexPatterns,
			fixer:    regexPatternsFixer,
		},
	} {
		r, err := os.Open(filename)
		require.NoError(err)

		actual, err := loadPatterns(filename, r)
		require.NoError(err)

		if tc.fixer != nil {
			tc.fixer(t, tc.expected, actual)
		}
		require.Equal(tc.expected, actual)
	}
}

var basicPatterns = []*pattern{{
	id:     "root",
	slash:  mustSlash,
	prefix: []part{},
	query: query{
		dedup: keepAll,
		match: map[string]*param{},
	},
	tests: map[string]string{
		"/":         "/",
		"/rfc3092/": "",
	},
}, {
	id:    "slash-required",
	slash: mustSlash,
	prefix: []part{
		&plainPart{"foo"},
	},
	query: query{
		dedup: keepFirst,
		match: map[string]*param{},
	},
	tests: map[string]string{
		"/foo":  "",
		"/foo/": "/foo/",
		"/bar":  "",
		"/baz/": "",
	},
}, {
	id:    "no-final-slash",
	slash: neverSlash,
	prefix: []part{
		&plainPart{"bar"},
	},
	query: query{
		dedup: keepLast,
		match: map[string]*param{},
	},
	tests: map[string]string{
		"/foo":  "",
		"/foo/": "",
		"/bar":  "/bar",
		"/bar/": "",
	},
}, {
	id:    "optional-slash",
	slash: maySlash,
	prefix: []part{
		&plainPart{"baz"},
	},
	query: query{
		dedup: keepAll,
		match: map[string]*param{},
	},
	tests: map[string]string{
		"/foo":  "",
		"/foo/": "",
		"/baz":  "/baz",
		"/baz/": "/baz",
	},
}, {
	id:    "regex",
	slash: maySlash,
	prefix: []part{
		&regexPart{
			regex:    regexp.MustCompile("qux"),
			rewriter: &staticStringRewriter{"quux"},
		},
	},
	query: query{
		dedup: keepAll,
		match: map[string]*param{},
	},
	tests: map[string]string{
		"/qux/": "/quux",
	},
}, {
	id:    "goldilocks",
	slash: neverSlash,
	prefix: []part{
		&plainPart{"corge"},
		&plainPart{"grault"},
		&plainPart{"garply"},
	},
	query: query{},
	tests: map[string]string{
		"/corge":                    "",
		"/corge/grault":             "",
		"/corge/grault/garply":      "/corge/grault/garply",
		"/corge/grault/garply/":     "",
		"/corge/grault/garply/fred": "",
	},
}, {
	id:    "query",
	slash: maySlash,
	prefix: []part{
		&plainPart{"search"},
	},
	query: query{
		dedup: keepAll,
		match: map[string]*param{
			"q": {
				remove:   false,
				rewriter: &staticQueryRewriter{"X"},
			},
			"utf8": {remove: true},
		},
	},
	tests: map[string]string{
		"/search?utf8=✔":         "/search",
		"/search/?q=cats":        "/search?q=X",
		"/search/?q=dogs&utf8=✔": "/search?q=X",
	},
}}

var basicTree = &Patterns{tree: tree{
	id:    "root",
	slash: mustSlash,
	query: query{
		dedup: keepAll,
		match: map[string]*param{},
	},
	children: []*child{{
		part: &plainPart{"foo"},
		tree: &tree{
			id:    "slash-required",
			slash: mustSlash,
			query: query{
				dedup: keepFirst,
				match: map[string]*param{},
			},
		},
	}, {
		part: &plainPart{"bar"},
		tree: &tree{
			id:    "no-final-slash",
			slash: neverSlash,
			query: query{
				dedup: keepLast,
				match: map[string]*param{},
			},
		},
	}, {
		part: &plainPart{"baz"},
		tree: &tree{
			id:    "optional-slash",
			slash: maySlash,
			query: query{
				dedup: keepAll,
				match: map[string]*param{},
			},
		},
	}, {
		part: &regexPart{
			regex:    regexp.MustCompile("qux"),
			rewriter: &staticStringRewriter{"quux"},
		},
		tree: &tree{
			id:    "regex",
			slash: maySlash,
			query: query{
				dedup: keepAll,
				match: map[string]*param{},
			},
		},
	}, {
		part: &plainPart{"corge"},
		tree: &tree{
			children: []*child{{
				part: &plainPart{"grault"},
				tree: &tree{
					children: []*child{{
						part: &plainPart{"garply"},
						tree: &tree{
							id:    "goldilocks",
							slash: neverSlash,
						},
					}},
				},
			}},
		},
	}, {
		part: &plainPart{"search"},
		tree: &tree{
			id:    "query",
			slash: maySlash,
			query: query{
				dedup: keepAll,
				match: map[string]*param{
					"q": {
						remove:   false,
						rewriter: &staticQueryRewriter{"X"},
					},
					"utf8": {remove: true},
				},
			},
		},
	}},
}, tests: map[string]result{
	"/":                         {"root", "/"},
	"/rfc3092/":                 {"root", ""},
	"/foo":                      {"slash-required", ""},
	"/foo/":                     {"slash-required", "/foo/"},
	"/bar":                      {"no-final-slash", "/bar"},
	"/bar/":                     {"no-final-slash", ""},
	"/baz":                      {"optional-slash", "/baz"},
	"/baz/":                     {"optional-slash", "/baz"},
	"/qux/":                     {"regex", "/quux"},
	"/corge":                    {"goldilocks", ""},
	"/corge/grault":             {"goldilocks", ""},
	"/corge/grault/garply":      {"goldilocks", "/corge/grault/garply"},
	"/corge/grault/garply/":     {"goldilocks", ""},
	"/corge/grault/garply/fred": {"goldilocks", ""},
	"/search?utf8=✔":            {"query", "/search"},
	"/search/?q=cats":           {"query", "/search?q=X"},
	"/search/?q=dogs&utf8=✔":    {"query", "/search?q=X"},
}}

func regexPatternsFixer(t *testing.T, expected, actual []*pattern) {
	require := require.New(t)

	for i := range expected {
		actual := actual[i]
		expected := expected[i]

		if len(actual.prefix) > 1 {
			actualPart := actual.prefix[1].(*regexPart)
			require.IsType(&callableStringRewriter{}, actualPart.rewriter)
			expectedPart := expected.prefix[1].(*regexPart)
			expectedPart.rewriter = actualPart.rewriter
		}

		if actual.suffix != nil {
			actualPart := actual.suffix
			require.IsType(&callableStringRewriter{}, actualPart.rewriter)
			expectedPart := expected.suffix
			expectedPart.rewriter = actualPart.rewriter
		}

		if actualParam, ok := actual.query.match["utf8"]; ok {
			expectedParam := expected.query.match["utf8"]
			expectedParam.rewriter = actualParam.rewriter
		}
	}
}

var regexPatterns = []*pattern{{
	id:    "prefix-regex",
	slash: mustSlash,
	prefix: []part{
		&regexPart{
			regex:    regexp.MustCompile("foo|bar"),
			rewriter: &staticStringRewriter{"baz"},
		},
		&regexPart{
			regex:    regexp.MustCompile("qux|quux"),
			rewriter: nil,
		},
	},
	query: query{
		dedup: keepFirst,
		match: map[string]*param{
			"utf8": {
				remove:   false,
				rewriter: nil,
			},
		},
	},
	tests: map[string]string{
		"/foo/qux":          "",
		"/foo/qux/":         "/baz/n=3/",
		"/bar/qux/":         "/baz/n=3/",
		"/bar/quux/":        "/baz/n=4/",
		"/foo/quux/?utf8=✔": "/baz/n=4/?utf8=True",
	},
}, {
	id:    "suffix-regex",
	slash: maySlash,
	prefix: []part{
		&plainPart{"corge"},
	},
	suffix: &regexPart{
		suffix:   true,
		regex:    regexp.MustCompile(".+"),
		rewriter: nil,
	},
	query: query{
		dedup: keepLast,
		match: map[string]*param{
			"utf8": {
				remove:   false,
				rewriter: nil,
			},
		},
	},
	tests: map[string]string{
		"/corge/":            "",
		"/corge/grault":      "/corge/GRAULT",
		"/corge/garply":      "/corge/GARPLY",
		"/corge/waldo/":      "/corge/WALDO/",
		"/corge/fred?utf8=✔": "/corge/FRED?utf8=True",
		"/corge/fred?utf8=!": "/corge/FRED?utf8=False",
	},
}, {
	id:    "reject-regex",
	slash: maySlash,
	prefix: []part{
		&regexPart{
			regex:    regexp.MustCompile("[a-z]"),
			reject:   regexp.MustCompile("[aeiou]"),
			rewriter: &staticStringRewriter{"X"},
		},
	},
	query: query{},
	tests: map[string]string{
		"/a":  "",
		"/b":  "/X",
		"/c/": "/X",
	},
}}

func regexTreeFixer(t *testing.T, expected, actual *tree) {
	require := require.New(t)

	require.Equal(len(expected.children), len(actual.children))
	for i, expectedChild := range expected.children {
		actualChild := actual.children[i]
		require.IsType(expectedChild.part, actualChild.part)
		if expectedPart, ok := expectedChild.part.(*regexPart); ok {
			actualPart := actualChild.part.(*regexPart)
			expectedPart.rewriter = actualPart.rewriter
		}
		if expectedChild.tree != nil {
			require.NotNil(actualChild.tree)
			regexTreeFixer(t, expectedChild.tree, actualChild.tree)
		}
	}

	expectedParams := expected.query.match
	actualParams := actual.query.match
	require.Equal(len(expectedParams), len(actualParams))
	for k, expectedParam := range expectedParams {
		actualParam, ok := actualParams[k]
		require.True(ok, k)
		expectedParam.rewriter = actualParam.rewriter
	}
}

var regexTree = &Patterns{tree: tree{
	children: []*child{{
		part: &regexPart{
			regex:    regexp.MustCompile("foo|bar"),
			rewriter: &staticStringRewriter{"baz"},
		},
		tree: &tree{
			children: []*child{{
				part: &regexPart{
					regex:    regexp.MustCompile("qux|quux"),
					rewriter: nil,
				},
				tree: &tree{
					id:    "prefix-regex",
					slash: mustSlash,
					query: query{
						dedup: keepFirst,
						match: map[string]*param{
							"utf8": {
								remove:   false,
								rewriter: nil,
							},
						},
					},
				},
			}},
		},
	}, {
		part: &plainPart{"corge"},
		tree: &tree{
			children: []*child{{
				part: &regexPart{
					suffix:   true,
					regex:    regexp.MustCompile(".+"),
					rewriter: nil,
				},
				tree: &tree{
					id:    "suffix-regex",
					slash: maySlash,
					query: query{
						dedup: keepLast,
						match: map[string]*param{
							"utf8": {
								remove:   false,
								rewriter: nil,
							},
						},
					},
				},
			}},
		},
	}, {
		part: &regexPart{
			regex:    regexp.MustCompile("[a-z]"),
			reject:   regexp.MustCompile("[aeiou]"),
			rewriter: &staticStringRewriter{"X"},
		},
		tree: &tree{
			id:    "reject-regex",
			slash: maySlash,
			query: query{},
		},
	}},
}, tests: map[string]result{
	"/a":                 {"reject-regex", ""},
	"/b":                 {"reject-regex", "/X"},
	"/c/":                {"reject-regex", "/X"},
	"/foo/qux":           {"prefix-regex", ""},
	"/foo/qux/":          {"prefix-regex", "/baz/n=3/"},
	"/bar/qux/":          {"prefix-regex", "/baz/n=3/"},
	"/bar/quux/":         {"prefix-regex", "/baz/n=4/"},
	"/foo/quux/?utf8=✔":  {"prefix-regex", "/baz/n=4/?utf8=True"},
	"/corge/":            {"suffix-regex", ""},
	"/corge/grault":      {"suffix-regex", "/corge/GRAULT"},
	"/corge/garply":      {"suffix-regex", "/corge/GARPLY"},
	"/corge/waldo/":      {"suffix-regex", "/corge/WALDO/"},
	"/corge/fred?utf8=✔": {"suffix-regex", "/corge/FRED?utf8=True"},
	"/corge/fred?utf8=!": {"suffix-regex", "/corge/FRED?utf8=False"},
}}

package patterns

import (
	"os"
	"regexp"
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
	slash:  maySlash,
	prefix: []part{},
	params: params{
		dedup:  keepAll,
		params: map[string]*param{},
	},
	tests: map[string]string{
		"/":         "/",
		"/rfc3092/": "",
	},
}, {
	id:    "slash-required",
	slash: mustSlash,
	prefix: []part{
		&plainPart{
			value: "foo",
		},
	},
	params: params{
		dedup:  keepFirst,
		params: map[string]*param{},
	},
	tests: map[string]string{
		"/foo":  "",
		"/foo/": "/foo/",
	},
}, {
	id:    "no-final-slash",
	slash: neverSlash,
	prefix: []part{
		&plainPart{
			value: "bar",
		},
	},
	params: params{
		dedup:  keepLast,
		params: map[string]*param{},
	},
	tests: map[string]string{
		"/bar":  "/bar",
		"/bar/": "",
	},
}, {
	id:    "optional-slash",
	slash: maySlash,
	prefix: []part{
		&plainPart{
			value: "baz",
		},
	},
	params: params{
		dedup:  keepAll,
		params: map[string]*param{},
	},
	tests: map[string]string{
		"/baz":  "/baz",
		"/baz/": "/baz",
	},
}, {
	id:    "regex",
	slash: maySlash,
	prefix: []part{
		&regexPart{
			regex: regexp.MustCompile("qux"),
			rewriter: &staticRewriter{
				value: "quux",
			},
		},
	},
	params: params{
		dedup:  keepAll,
		params: map[string]*param{},
	},
	tests: map[string]string{
		"/qux/": "/quux",
	},
}, {
	id:    "query",
	slash: maySlash,
	prefix: []part{
		&plainPart{
			value: "search",
		},
	},
	params: params{
		dedup: keepAll,
		params: map[string]*param{
			"q": {
				remove: false,
				rewriter: &staticRewriter{
					value: "X",
				},
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

var basicTree = &Patterns{tree{
	id:    "root",
	slash: maySlash,
	params: params{
		dedup:  keepAll,
		params: map[string]*param{},
	},
	children: []*child{{
		part: &plainPart{
			value: "foo",
		},
		tree: &tree{
			id:    "slash-required",
			slash: mustSlash,
			params: params{
				dedup:  keepFirst,
				params: map[string]*param{},
			},
		},
	}, {
		part: &plainPart{
			value: "bar",
		},
		tree: &tree{
			id:    "no-final-slash",
			slash: neverSlash,
			params: params{
				dedup:  keepLast,
				params: map[string]*param{},
			},
		},
	}, {
		part: &plainPart{
			value: "baz",
		},
		tree: &tree{
			id:    "optional-slash",
			slash: maySlash,
			params: params{
				dedup:  keepAll,
				params: map[string]*param{},
			},
		},
	}, {
		part: &regexPart{
			regex: regexp.MustCompile("qux"),
			rewriter: &staticRewriter{
				value: "quux",
			},
		},
		tree: &tree{
			id:    "regex",
			slash: maySlash,
			params: params{
				dedup:  keepAll,
				params: map[string]*param{},
			},
		},
	}, {
		part: &plainPart{
			value: "search",
		},
		tree: &tree{
			id:    "query",
			slash: maySlash,
			params: params{
				dedup: keepAll,
				params: map[string]*param{
					"q": {
						remove: false,
						rewriter: &staticRewriter{
							value: "X",
						},
					},
					"utf8": {remove: true},
				},
			},
		},
	}},
}}

func regexPatternsFixer(t *testing.T, expected, actual []*pattern) {
	require := require.New(t)

	for i := range expected {
		actual := actual[i]
		expected := expected[i]

		if len(actual.prefix) > 1 {
			actualPart := actual.prefix[1].(*regexPart)
			require.IsType(&callableRewriter{}, actualPart.rewriter)
			expectedPart := expected.prefix[1].(*regexPart)
			expectedPart.rewriter = actualPart.rewriter
		}

		if actual.suffix != nil {
			actualPart := actual.suffix
			require.IsType(&callableRewriter{}, actualPart.rewriter)
			expectedPart := expected.suffix
			expectedPart.rewriter = actualPart.rewriter
		}

		actualParam := actual.params.params["utf8"]
		expectedParam := expected.params.params["utf8"]
		expectedParam.rewriter = actualParam.rewriter
	}
}

var regexPatterns = []*pattern{{
	id:    "prefix-regex",
	slash: mustSlash,
	prefix: []part{
		&regexPart{
			regex: regexp.MustCompile("foo|bar"),
			rewriter: &staticRewriter{
				value: "baz",
			},
		},
		&regexPart{
			regex:    regexp.MustCompile("qux|quux"),
			rewriter: nil,
		},
	},
	params: params{
		dedup: keepFirst,
		params: map[string]*param{
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
		&plainPart{
			value: "corge",
		},
	},
	suffix: &regexPart{
		suffix:   true,
		regex:    regexp.MustCompile(".+"),
		rewriter: nil,
	},
	params: params{
		dedup: keepLast,
		params: map[string]*param{
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
		"/corge/waldo/":      "/corge/WALDO",
		"/corge/fred?utf8=✔": "/corge/FRED?utf8=True",
		"/corge/fred?utf8=!": "/corge/FRED?utf8=False",
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

	expectedParams := expected.params.params
	actualParams := actual.params.params
	require.Equal(len(expectedParams), len(actualParams))
	for k, expectedParam := range expectedParams {
		actualParam, ok := actualParams[k]
		require.True(ok, k)
		expectedParam.rewriter = actualParam.rewriter
	}
}

var regexTree = &Patterns{tree{
	children: []*child{{
		part: &regexPart{
			regex: regexp.MustCompile("foo|bar"),
			rewriter: &staticRewriter{
				value: "baz",
			},
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
					params: params{
						dedup: keepFirst,
						params: map[string]*param{
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
		part: &plainPart{
			value: "corge",
		},
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
					params: params{
						dedup: keepLast,
						params: map[string]*param{
							"utf8": {
								remove:   false,
								rewriter: nil,
							},
						},
					},
				},
			}},
		},
	}},
}}

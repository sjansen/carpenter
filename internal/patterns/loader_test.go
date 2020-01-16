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
	}{
		"testdata/basic.star": {
			expected: basicTree,
		},
	} {
		r, err := os.Open(filename)
		require.NoError(err)

		actual, err := Load(filename, r)
		require.NoError(err)
		require.Equal(tc.expected, actual)
	}
}

func TestLoadPatterns(t *testing.T) {
	require := require.New(t)

	for filename, tc := range map[string]struct {
		expected []*pattern
	}{
		"testdata/basic.star": {
			expected: basicPatterns,
		},
	} {
		r, err := os.Open(filename)
		require.NoError(err)

		actual, err := loadPatterns(filename, r)
		require.NoError(err)
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

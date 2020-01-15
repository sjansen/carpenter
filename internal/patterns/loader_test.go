package patterns

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	require := require.New(t)

	const filename = "testdata/core.star"
	r, err := os.Open(filename)
	require.NoError(err)
	actual, err := Load(filename, r)
	require.NoError(err)

	require.Equal(expectedCoreTree, actual)
}

func TestLoadPatterns(t *testing.T) {
	require := require.New(t)

	const filename = "testdata/core.star"
	r, err := os.Open(filename)
	require.NoError(err)
	actual, err := loadPatterns(filename, r)
	require.NoError(err)

	require.Equal(expectedCorePatterns, actual)
}

var expectedCorePatterns = []*pattern{
	{
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
	},
}

var expectedCoreTree = &Patterns{
	tree{
		id:    "root",
		slash: maySlash,
		params: params{
			dedup:  keepAll,
			params: map[string]*param{},
		},
		children: []*child{
			{
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
			},
			{
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
			},
		},
	},
}

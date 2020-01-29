package patterns

import (
	"bytes"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sjansen/carpenter/internal/logger"
	"github.com/sjansen/carpenter/internal/sys"
)

func TestMatch(t *testing.T) {
	require := require.New(t)

	for _, filename := range []string{
		"testdata/basic.star",
		"testdata/regex.star",
		"testdata/query.star",
	} {
		r, err := os.Open(filename)
		require.NoError(err)

		patterns, err := Load(filename, r)
		require.NoError(err)

		for raw, expected := range patterns.tests {
			url, err := url.Parse(raw)
			require.NoError(err, raw)

			id, normalized, err := patterns.Match(url)
			require.NoError(err, raw)

			if expected.url == "" {
				require.Equal("", id, raw)
			} else {
				require.Equal(expected.id, id, raw)
			}
			require.Equal(expected.url, normalized, raw)
		}
	}
}

func TestMatchAll(t *testing.T) {
	require := require.New(t)

	r := bytes.NewBufferString(`
add_url("first", path={"prefix": ["foo"], "suffix": "/"}, query={}, tests={})
add_url("second", path={"prefix": [(".+", "ANY")], "suffix": "/"}, query={}, tests={})
	`)

	patterns, err := Load("<buffer>", r)
	require.NoError(err)

	url, err := url.Parse("/foo/")
	require.NoError(err)

	actual, err := patterns.MatchAll(url)
	require.NoError(err)

	expected := map[string]string{
		"first":  "/foo/",
		"second": "/ANY/",
	}
	require.Equal(expected, actual)
}

func TestMatchErrors(t *testing.T) {
	files, _ := filepath.Glob("testdata/match-errors/*.star")
	for _, tc := range files {
		tc := tc
		prefix := tc[:len(tc)-5]
		t.Run(filepath.Base(prefix), func(t *testing.T) {
			require := require.New(t)

			r, err := os.Open(tc)
			require.NoError(err)

			patterns, err := Load(filepath.Base(tc), r)
			require.NoError(err)

			raw, err := ioutil.ReadFile(prefix + ".url")
			require.NoError(err)

			url, err := url.Parse(strings.TrimSpace(string(raw)))
			require.NoError(err)

			_, _, err = patterns.Match(url)
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

func TestTest(t *testing.T) {
	require := require.New(t)

	const filename = "testdata/test.star"
	r, err := os.Open(filename)
	require.NoError(err)

	patterns, err := Load(filename, r)
	require.NoError(err)

	var stdout, stderr bytes.Buffer
	sys := &sys.IO{
		Log:    logger.Discard(),
		Stdout: &stdout,
		Stderr: &stderr,
	}

	actual, err := patterns.Test(sys)
	require.NoError(err)

	expected := map[string]string{
		"/foo":  "",
		"/foo/": "first",
		"/bar/": "second",
		"/baz/": "second",
		"/qux/": "",
	}
	require.Equal(expected, actual)
}

func TestTestErrors(t *testing.T) {
	files, _ := filepath.Glob("testdata/test-errors/*.star")
	for _, tc := range files {
		tc := tc
		prefix := tc[:len(tc)-5]
		t.Run(filepath.Base(prefix), func(t *testing.T) {
			require := require.New(t)

			r, err := os.Open(tc)
			require.NoError(err)

			patterns, err := Load(filepath.Base(tc), r)
			require.NoError(err)

			var stdout, stderr bytes.Buffer
			sys := &sys.IO{
				Log:    logger.Discard(),
				Stdout: &stdout,
				Stderr: &stderr,
			}

			tests, err := patterns.Test(sys)
			require.Error(err)
			require.Nil(tests)

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

func TestValidFiles(t *testing.T) {
	files, _ := filepath.Glob("testdata/*.star")
	for _, tc := range files {
		tc := tc
		prefix := tc[:len(tc)-5]
		t.Run(filepath.Base(prefix), func(t *testing.T) {
			require := require.New(t)

			r, err := os.Open(tc)
			require.NoError(err)

			patterns, err := Load(filepath.Base(tc), r)
			require.NoError(err)

			var stdout, stderr bytes.Buffer
			sys := &sys.IO{
				Log:    logger.Discard(),
				Stdout: &stdout,
				Stderr: &stderr,
			}

			tests, err := patterns.Test(sys)
			require.NoError(err)
			require.NotNil(tests)
		})
	}
}

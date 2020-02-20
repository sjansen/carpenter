package pipeline_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sjansen/carpenter/internal/lazyio"
	"github.com/sjansen/carpenter/internal/patterns"
	"github.com/sjansen/carpenter/internal/pipeline"
	"github.com/sjansen/carpenter/internal/sys"
	"github.com/sjansen/carpenter/internal/tokenizer"
	"github.com/sjansen/carpenter/internal/uaparser"
)

func TestPipeline(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	r, err := os.Open("testdata/alb.star")
	require.NoError(err)

	patterns, err := patterns.Load("alb.star", r)
	require.NoError(err)

	uaparser, err := uaparser.UserAgentParser()
	require.NoError(err)

	debug := &lazyio.BufferWriter{}
	result := &lazyio.BufferWriter{}
	pipeline := &pipeline.Pipeline{
		Patterns:  patterns,
		Tokenizer: tokenizer.ALB,
		UAParser:  uaparser,
		IO:        sys.Discard(),
		Source:    &lazyio.FileReader{Dir: "testdata/src"},
		Result:    result,
		Debug:     debug,
	}

	pipeline.Start()
	pipeline.AddTask("alb.log")
	pipeline.Wait()

	expected, err := ioutil.ReadFile("testdata/dst/alb.csv")
	require.NoError(err)

	require.Equal([]string{"alb.csv"}, result.Buffers())

	actual := result.Buffer("alb.csv")
	if !assert.Equal(expected, actual.Bytes()) {
		ioutil.WriteFile("testdata/dst/alb.csv.actual", actual.Bytes(), 0644)
	}

	files, _ := filepath.Glob("testdata/debug/*/*.???")
	buffers := debug.Buffers()
	sort.Strings(buffers)
	require.Equal(
		[]string{"parse/alb.csv", "tokenize/alb.txt", "unrecognized/alb.csv"},
		buffers,
	)
	for _, f := range files {
		expected, err := ioutil.ReadFile(f)
		require.NoError(err)

		suffix := f[15:]
		actual := debug.Buffer(suffix)
		require.NotNil(actual, suffix)
		if !assert.Equal(expected, actual.Bytes(), suffix) {
			ioutil.WriteFile(f+".actual", actual.Bytes(), 0644)
		}
	}
}

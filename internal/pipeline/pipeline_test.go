package pipeline_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/sjansen/carpenter/internal/lazyio"
	"github.com/sjansen/carpenter/internal/oldpatterns"
	"github.com/sjansen/carpenter/internal/pipeline"
	"github.com/sjansen/carpenter/internal/tokenizer"
	"github.com/sjansen/carpenter/internal/uaparser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPipeline(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	r, err := os.Open("testdata/alb.star")
	require.NoError(err)

	patterns, err := oldpatterns.Load("alb.star", r)
	require.NoError(err)

	uaparser, err := uaparser.UserAgentParser()
	require.NoError(err)

	debug := &lazyio.BufferOpener{}
	result := &lazyio.BufferOpener{}
	pipeline := &pipeline.Pipeline{
		Debug:     debug,
		Result:    result,
		Patterns:  patterns,
		Tokenizer: tokenizer.ALB,
		UAParser:  uaparser,
	}

	src, err := os.Open("testdata/src/alb.log")
	require.NoError(err)

	task := pipeline.NewTask(src, "alb.log")
	err = task.Run()
	require.NoError(err)

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

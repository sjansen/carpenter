package tokenizer_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sjansen/carpenter/internal/tokenizer"
)

func TestALB(t *testing.T) {
	err := tokenizer.ALB.EnableUserAgentParsing()
	require.NoError(t, err)

	files, _ := filepath.Glob("testdata/alb-*.txt")
	for _, tc := range files {
		tc := tc
		prefix := tc[:len(tc)-4]
		t.Run(filepath.Base(prefix)[4:], func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			data, err := ioutil.ReadFile(prefix + ".expected")
			require.NoError(err)

			expected := map[string]string{}
			err = json.Unmarshal(data, &expected)
			require.NoError(err)

			data, err = ioutil.ReadFile(tc)
			require.NoError(err)
			line := string(bytes.TrimSpace(data))

			actual := tokenizer.ALB.Tokenize(line)
			if !assert.Equal(expected, actual) {
				data, err := json.MarshalIndent(actual, "", "  ")
				require.NoError(err)
				ioutil.WriteFile(prefix+".actual", data, 0644)
			}
		})
	}
}

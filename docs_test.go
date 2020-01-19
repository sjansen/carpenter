package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sjansen/carpenter/internal/cmd"
	"github.com/sjansen/carpenter/internal/logger"
	"github.com/sjansen/carpenter/internal/sys"
)

func TestExamples(t *testing.T) {
	const suffix = "-test-cases.json"
	files, _ := filepath.Glob("docs/examples/*" + suffix)
	for _, tc := range files {
		tc := tc
		prefix := tc[:len(tc)-len(suffix)]
		t.Run(filepath.Base(prefix)+suffix, func(t *testing.T) {
			require := require.New(t)

			var expected map[string]string
			swapi, err := os.Open(tc)
			require.NoError(err)
			dec := json.NewDecoder(swapi)
			err = dec.Decode(&expected)
			require.NoError(err)

			var stdout, stderr bytes.Buffer
			base := &cmd.Base{
				IO: sys.IO{
					Log:    logger.Discard(),
					Stdout: &stdout,
					Stderr: &stderr,
				},
			}
			cmd := &cmd.TestCasesCmd{
				File: prefix + ".star",
			}

			err = cmd.Run(base)
			require.NoError(err)
			require.NotZero(stdout.Len())

			var actual map[string]string
			dec = json.NewDecoder(&stdout)
			err = dec.Decode(&actual)
			require.NoError(err)

			require.Equal(expected, actual)
		})
	}
}

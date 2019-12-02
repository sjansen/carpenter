package cli_test

import (
	"bytes"
	"runtime/debug"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sjansen/carpenter/internal/cli"
	"github.com/sjansen/carpenter/internal/cmd"
)

func TestVersion(t *testing.T) {
	require := require.New(t)

	info, _ := debug.ReadBuildInfo()
	parser := cli.RegisterCommands("test-version")
	for _, tc := range []struct {
		args        []string
		expected    interface{}
		expectError bool
	}{{
		args: []string{
			"version",
		},
		expected: &cmd.VersionCmd{
			App:       "carpenter",
			Build:     "test",
			BuildInfo: info,
			Verbose:   false,
		},
	}, {
		args: []string{
			"version", "--long",
		},
		expected: &cmd.VersionCmd{
			App:       "carpenter",
			Build:     "test",
			BuildInfo: info,
			Verbose:   true,
		},
	}} {
		c, err := parser.Parse(tc.args)
		if tc.expectError {
			require.Error(err)
		} else {
			require.NoError(err)

			var stdout, stderr bytes.Buffer
			err = c(&stdout, &stderr)
			require.NoError(err)

			require.Contains(stdout.String(), "test-version")
			require.Empty(stderr)
		}
	}
}

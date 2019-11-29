package cli_test

import (
	"bytes"
	"runtime/debug"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sjansen/carpenter/internal/cli"
	"github.com/sjansen/carpenter/internal/commands"
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
		expected: &commands.VersionCmd{
			App:       "carpenter",
			Build:     "test",
			BuildInfo: info,
			Verbose:   false,
		},
	}, {
		args: []string{
			"version", "--long",
		},
		expected: &commands.VersionCmd{
			App:       "carpenter",
			Build:     "test",
			BuildInfo: info,
			Verbose:   true,
		},
	}} {
		cmd, err := parser.Parse(tc.args)
		if tc.expectError {
			require.Error(err)
		} else {
			require.NoError(err)

			var stdout, stderr bytes.Buffer
			err = cmd(&stdout, &stderr)
			require.NoError(err)

			require.Contains(stdout.String(), "test-version")
			require.Empty(stderr)
		}
	}
}

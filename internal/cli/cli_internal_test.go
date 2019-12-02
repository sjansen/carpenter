package cli

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sjansen/carpenter/internal/cmd"
)

func TestArgParser(t *testing.T) {
	require := require.New(t)

	parser := RegisterCommands("test")
	for _, tc := range []struct {
		args        []string
		expected    interface{}
		expectError bool
	}{{
		args: []string{
			"version",
		},
		expected: &cmd.VersionCmd{
			App:   "carpenter",
			Build: "test",
		},
	}} {
		_, err := parser.Parse(tc.args)
		if tc.expectError {
			require.Error(err)
		} else {
			require.NoError(err)
			require.Equal(tc.expected, parser.cmd)
		}
	}
}

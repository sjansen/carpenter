package logger_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sjansen/carpenter/internal/logger"
)

func TestNew(t *testing.T) {
	require := require.New(t)

	var stdout, stderr bytes.Buffer
	log := logger.New(1, &stdout, &stderr)
	require.NotNil(log)
	require.True(stdout.Len() < 1)
	require.True(stderr.Len() < 1)

	log.Debugw("Can you hear me now?")
	require.Empty(stdout.String())
	require.True(stderr.Len() > 0)

	log.Warnw("Can you hear me now?")
	require.NotEmpty(stdout.String())
}

package useragent

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	require := require.New(t)

	uagent := "Mozilla/5.0 (iPad; U; CPU OS 3_2_1 like Mac OS X; en-us) AppleWebKit/531.21.10 (KHTML, like Gecko) Mobile/7B405"
	parser, err := Parser()
	require.NoError(err)

	client := parser.Parse(uagent)
	require.Equal("iOS", client.Os.Family)
	require.Equal("iPad", client.Device.Family)
	require.Equal("Mobile Safari UI/WKWebView", client.UserAgent.Family)
}

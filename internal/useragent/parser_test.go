package useragent

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var chrome = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36" +
	" (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36"
var firefox = "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:47.0) Gecko/20100101 Firefox/47.0"
var safari = "Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_1 like Mac OS X) AppleWebKit/603.1.30" +
	" (KHTML, like Gecko) Version/10.0 Mobile/14E304 Safari/602.1"

func TestParser(t *testing.T) {
	require := require.New(t)

	parser, err := Parser()
	require.NoError(err)

	browser := parser.Parse(chrome)
	require.Equal("Linux", browser.Os.Family)
	require.Equal("Other", browser.Device.Family)
	require.Equal("Chrome", browser.UserAgent.Family)

	browser = parser.Parse(firefox)
	require.Equal("Windows", browser.Os.Family)
	require.Equal("Other", browser.Device.Family)
	require.Equal("Firefox", browser.UserAgent.Family)

	browser = parser.Parse(safari)
	require.Equal("iOS", browser.Os.Family)
	require.Equal("iPhone", browser.Device.Family)
	require.Equal("Mobile Safari", browser.UserAgent.Family)
}

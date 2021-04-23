package uaparser

import (
	_ "embed" //nolint

	"github.com/ua-parser/uap-go/uaparser"
)

//go:embed data/regexes.yaml
var data []byte

type Parser uaparser.Parser
type Client uaparser.Client

func UserAgentParser() (*Parser, error) {
	uap, err := uaparser.NewFromBytes(data)
	if err != nil {
		return nil, err
	}

	return (*Parser)(uap), nil
}

func (parser *Parser) Parse(line string) *Client {
	p := (*uaparser.Parser)(parser)
	c := p.Parse(line)
	return (*Client)(c)
}

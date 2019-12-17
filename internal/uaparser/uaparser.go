package uaparser

import (
	"io/ioutil"

	"github.com/sjansen/carpenter/internal/data"
	"github.com/ua-parser/uap-go/uaparser"
)

type Parser uaparser.Parser
type Client uaparser.Client

func UserAgentParser() (*Parser, error) {
	r, err := data.Assets.Open("regexes.yaml")
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	uap, err := uaparser.NewFromBytes(bytes)
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

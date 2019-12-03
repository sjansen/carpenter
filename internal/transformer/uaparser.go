package transformer

import (
	"io/ioutil"

	"github.com/sjansen/carpenter/internal/data"
	"github.com/ua-parser/uap-go/uaparser"
)

func UserAgentParser() (*uaparser.Parser, error) {
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

	return uap, nil
}

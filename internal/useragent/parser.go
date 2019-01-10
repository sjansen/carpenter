package useragent

import (
	"io/ioutil"

	"github.com/sjansen/carpenter/internal/data"
	"github.com/ua-parser/uap-go/uaparser"
)

func Parser() (*uaparser.Parser, error) {
	r, err := data.Assets.Open("regexes.yaml")
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return uaparser.NewFromBytes(bytes)
}

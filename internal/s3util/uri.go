package s3util

import (
	"fmt"
	"net/url"
)

type URI struct {
	Bucket string
	Key    string
}

func ParseURI(s string) (*URI, error) {
	parsed, err := url.Parse(s)
	switch {
	case err != nil:
		return nil, err
	case parsed.Scheme != "s3":
		err = fmt.Errorf(`invalid URI scheme: expected="s3" actual=%q`, parsed.Scheme)
		return nil, err
	}
	uri := &URI{Bucket: parsed.Host}
	if len(parsed.Path) > 0 {
		uri.Key = parsed.Path[1:]
	}
	return uri, nil
}

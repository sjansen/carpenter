package s3util

import (
	"fmt"
	"net/url"
)

type URI struct {
	Profile  string
	Region   string
	Bucket   string
	Key      string
	Endpoint string
}

func ParseURI(s string) (*URI, error) {
	parsed, err := url.Parse(s)
	switch {
	case err != nil:
		return nil, err
	case parsed.Scheme != "s3":
		err = fmt.Errorf(`invalid S3 URI scheme: expected="s3" actual=%q`, parsed.Scheme)
		return nil, err
	}
	uri := &URI{Bucket: parsed.Host}
	if len(parsed.Path) > 0 {
		uri.Key = parsed.Path[1:]
	}
	if len(parsed.RawQuery) > 0 {
		query, err := url.ParseQuery(parsed.RawQuery)
		if err != nil {
			return nil, err
		}
		for k, v := range query {
			switch k {
			case "endpoint":
				uri.Endpoint = v[0]
			case "profile":
				uri.Profile = v[0]
			case "region":
				uri.Region = v[0]
			default:
				err = fmt.Errorf(`unexpected S3 URI parameter: %q`, v)
				return nil, err
			}
		}
	}
	return uri, nil
}

func (uri *URI) ToUploaderConfig() *UploaderConfig {
	return &UploaderConfig{
		Profile:  uri.Profile,
		Region:   uri.Region,
		Bucket:   uri.Bucket,
		Prefix:   uri.Key,
		Endpoint: uri.Endpoint,
	}
}

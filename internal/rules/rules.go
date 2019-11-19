package rules

import (
	"fmt"
	"io"
	"net/url"
	"strings"

	"go.starlark.net/starlark"
)

type Rules []*Rule

type Rule struct {
	id    string
	slash string
	parts []part
	tests map[string]string
}

func Load(filename string, src io.Reader) (Rules, error) {
	loader := &rulesLoader{
		rules: make([]*Rule, 0),
	}

	loader.Builtin = starlark.NewBuiltin("register_urls", loader.registerURLs)

	globals := starlark.StringDict{
		"register_urls": loader,
	}
	thread := &starlark.Thread{}
	_, err := starlark.ExecFile(thread, filename, src, globals)
	if err != nil {
		return nil, err
	}

	return loader.rules, nil
}

func (rules Rules) SelfTest() error {
	seen := map[string]string{}
	for _, r := range rules {
		for rawurl, expected := range r.tests {
			if expected != "" {
				if prev, ok := seen[rawurl]; ok {
					return fmt.Errorf(
						`invalid test case: %q (alreay matched by %q)`,
						rawurl, prev,
					)
				}
				seen[rawurl] = r.id
			}
			if actual, err := r.Match(rawurl); err != nil {
				return err
			} else if expected != actual {
				return fmt.Errorf(
					"self-test failed: given=%q expected=%q actual=%q",
					rawurl, expected, actual,
				)
			}
		}
	}
	return nil
}

func (r *Rule) Match(rawurl string) (string, error) {
	url, err := r.validateTestCase(rawurl)
	if err != nil {
		return "", err
	}

	parts, ok := r.splitPath(url.Path)
	if !ok || len(parts) != len(r.parts) {
		return "", nil
	}

	for i, part := range parts {
		if !r.parts[i].match(part) {
			return "", nil
		}
	}

	return r.rewriteURL(parts)
}

func (r *Rule) rewriteURL(parts []string) (string, error) {
	result := make([]string, 1, len(parts)+2)
	result[0] = ""

	thread := &starlark.Thread{}
	for i, part := range parts {
		part, err := r.parts[i].normalize(thread, part)
		if err != nil {
			return "", err
		}
		result = append(result, part)
	}
	if r.slash == "always" {
		result = append(result, "")
	}

	return strings.Join(result, "/"), nil
}

func (r *Rule) splitPath(path string) ([]string, bool) {
	plen := len(path)
	if plen > 1 {
		switch r.slash {
		case "always":
			if path[plen-1] == '/' {
				path = path[:plen-1]
			} else {
				return nil, false
			}
		case "never":
			if path[plen-1] == '/' {
				return nil, false
			}
		case "strip":
			if path[plen-1] == '/' {
				path = path[:plen-1]
			}
		}
	}

	parts := strings.Split(path, "/")
	return parts[1:], true
}

func (r *Rule) validateTestCase(rawurl string) (*url.URL, error) {
	if len(rawurl) < 1 || rawurl[0] != '/' {
		return nil, fmt.Errorf(`invalid test case: %q (should start with "/")`, rawurl)
	}

	url, err := url.Parse(rawurl)
	if err != nil {
		return nil, fmt.Errorf(`invalid test case: %q (%s)`, rawurl, err.Error())
	}

	return url, nil
}

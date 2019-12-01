package rules

import (
	"fmt"
	"io"
	"net/url"
	"sort"
	"strings"

	"github.com/sjansen/carpenter/internal/sys"
	"go.starlark.net/starlark"
)

type Rules []*Rule

type Rule struct {
	id     string
	dedup  string
	slash  string
	parts  []part
	params map[string]*param
	tests  map[string]string
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

func (rules Rules) Match(rawurl string) (id, url string, err error) {
	for _, r := range rules {
		match, err := r.Match(rawurl)
		if err != nil {
			return "", "", err
		} else if match != "" {
			return r.id, match, nil
		}
	}
	return "", "", nil
}

func (rules Rules) SelfTest(sys *sys.IO) (map[string]string, error) {
	matches := map[string]string{}
	rawurls := make([]string, 0, len(rules))
	sys.Log.Debug("starting rule-specific tests...")
	for _, r := range rules {
		for rawurl, expected := range r.tests {
			sys.Log.Debugf("testing rule=%q url=%q", r.id, rawurl)
			if expected != "" {
				if prev, ok := matches[rawurl]; ok {
					return nil, fmt.Errorf(
						`duplicate test case: %q (rule=%q test=%q)`,
						prev, r.id, rawurl,
					)
				}
				matches[rawurl] = r.id
				rawurls = append(rawurls, rawurl)
			}
			if actual, err := r.Test(rawurl); err != nil {
				return nil, err
			} else if expected != actual {
				return nil, fmt.Errorf(
					"unexpected result: expected=%q actual=%q (rule=%q test=%q)",
					expected, actual, r.id, rawurl,
				)
			}
		}
	}
	sys.Log.Debug("testing against all rules...")
	sort.Strings(rawurls)
	for _, rawurl := range rawurls {
		sys.Log.Debugf("testing url=%q", rawurl)
		expected := matches[rawurl]
		actual, _, err := rules.Match(rawurl)
		if err != nil {
			return nil, err // should be unreachable
		} else if expected != actual {
			return nil, fmt.Errorf(
				"unexpected rule match: expected=%q actual=%q (test=%q)",
				expected, actual, rawurl,
			)
		}
	}
	return matches, nil
}

func (r *Rule) Match(rawurl string) (string, error) {
	url, err := url.Parse(rawurl)
	if err != nil {
		return "", fmt.Errorf(`unable to parse url: %q (%s)`, rawurl, err.Error())
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

	path, err := r.rewriteURL(parts)
	if err != nil {
		return "", err
	}

	query, err := r.rewriteQuery(url.Query())
	if err != nil {
		return "", err
	} else if query != "" {
		return path + "?" + query, nil
	}

	return path, nil
}

func (r *Rule) Test(rawurl string) (string, error) {
	if len(rawurl) < 1 || rawurl[0] != '/' {
		return "", fmt.Errorf(`invalid test case: %q (should start with "/")`, rawurl)
	}

	return r.Match(rawurl)
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
	if r.slash == "always" || len(parts) == 0 {
		result = append(result, "")
	}

	return strings.Join(result, "/"), nil
}

func (r *Rule) rewriteQuery(query url.Values) (string, error) {
	result := url.Values{}

	thread := &starlark.Thread{}
	for key, values := range query {
		if p, ok := r.params[key]; ok {
			values, err := p.normalize(thread, r.dedup, values)
			if err != nil {
				return "", err
			}
			result[key] = values
		} else {
			result[key] = values
		}
	}

	return result.Encode(), nil
}

func (r *Rule) splitPath(path string) ([]string, bool) {
	plen := len(path)
	if plen > 0 {
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

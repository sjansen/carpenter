package patterns

import (
	"fmt"
	"io"
	"net/url"
	"sort"
	"strings"

	"github.com/sjansen/carpenter/internal/sys"
	"go.starlark.net/starlark"
)

type Patterns []*Pattern

type Pattern struct {
	id     string
	dedup  string
	slash  string
	prefix []part
	suffix *regexPart
	params map[string]*param
	tests  map[string]string
}

func Load(filename string, src io.Reader) (Patterns, error) {
	loader := &patternLoader{
		patterns: make([]*Pattern, 0),
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

	return loader.patterns, nil
}

func (patterns Patterns) Match(url *url.URL) (id, normalized string, err error) {
	for _, p := range patterns {
		match, err := p.Match(url)
		if err != nil {
			return "", "", err
		} else if match != "" {
			return p.id, match, nil
		}
	}
	return "", "", nil
}

func (patterns Patterns) SelfTest(sys *sys.IO) (map[string]string, error) {
	matches := map[string]string{}
	rawurls := make([]string, 0, len(patterns))
	sys.Log.Debug("starting pattern-specific tests...")
	for _, p := range patterns {
		for rawurl, expected := range p.tests {
			sys.Log.Debugf("testing pattern=%q url=%q", p.id, rawurl)
			if expected != "" {
				if prev, ok := matches[rawurl]; ok {
					return nil, fmt.Errorf(
						`duplicate test case: %q (pattern=%q test=%q)`,
						prev, p.id, rawurl,
					)
				}
				matches[rawurl] = p.id
				rawurls = append(rawurls, rawurl)
			}
			if actual, err := p.Test(rawurl); err != nil {
				return nil, err
			} else if expected != actual {
				return nil, fmt.Errorf(
					"unexpected result: expected=%q actual=%q (pattern=%q test=%q)",
					expected, actual, p.id, rawurl,
				)
			}
		}
	}
	sys.Log.Debug("testing against all patterns...")
	sort.Strings(rawurls)
	for _, rawurl := range rawurls {
		sys.Log.Debugf("testing url=%q", rawurl)
		expected := matches[rawurl]
		actual, _, err := patterns.Test(rawurl)
		if err != nil {
			return nil, err // should be unreachable
		} else if expected != actual {
			return nil, fmt.Errorf(
				"unexpected pattern match: expected=%q actual=%q (test=%q)",
				expected, actual, rawurl,
			)
		}
	}
	return matches, nil
}

func (patterns Patterns) Test(rawurl string) (id, normalized string, err error) {
	url, err := url.Parse(rawurl)
	if err != nil {
		return "", "", fmt.Errorf(`unable to parse url: %q (%s)`, rawurl, err.Error())
	}

	return patterns.Match(url)
}

func (p *Pattern) Match(url *url.URL) (string, error) {
	parts, ok := p.splitPath(url.Path)
	if !ok || len(parts) != len(p.prefix) {
		return "", nil
	}

	for i, part := range parts {
		if !p.prefix[i].match(part) {
			return "", nil
		}
	}

	path, err := p.rewriteURL(parts)
	if err != nil {
		return "", err
	}

	query, err := p.rewriteQuery(url.Query())
	if err != nil {
		return "", err
	} else if query != "" {
		return path + "?" + query, nil
	}

	return path, nil
}

func (p *Pattern) Test(rawurl string) (string, error) {
	if len(rawurl) < 1 || rawurl[0] != '/' {
		return "", fmt.Errorf(`invalid test case: %q (should start with "/")`, rawurl)
	}

	url, err := url.Parse(rawurl)
	if err != nil {
		return "", fmt.Errorf(`unable to parse url: %q (%s)`, rawurl, err.Error())
	}

	return p.Match(url)
}

func (p *Pattern) rewriteURL(parts []string) (string, error) {
	result := make([]string, 1, len(parts)+2)
	result[0] = ""

	thread := &starlark.Thread{}
	for i, part := range parts {
		part, err := p.prefix[i].normalize(thread, part)
		if err != nil {
			return "", err
		}
		result = append(result, part)
	}
	if p.slash == "/" || len(parts) == 0 {
		result = append(result, "")
	}

	return strings.Join(result, "/"), nil
}

func (p *Pattern) rewriteQuery(query url.Values) (string, error) {
	result := url.Values{}

	thread := &starlark.Thread{}
	for key, values := range query {
		if param, ok := p.params[key]; ok {
			values, err := param.normalize(thread, p.dedup, values)
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

func (p *Pattern) splitPath(path string) ([]string, bool) {
	plen := len(path)
	if plen > 0 {
		switch p.slash {
		case "/":
			if path[plen-1] == '/' {
				path = path[:plen-1]
			} else {
				return nil, false
			}
		case "/?":
			if path[plen-1] == '/' {
				path = path[:plen-1]
			}
		case "":
			if path[plen-1] == '/' {
				return nil, false
			}
		}
	}

	var parts []string
	if p.suffix != nil {
		parts = strings.SplitN(path, "/", len(p.prefix)+1)
	} else {
		parts = strings.Split(path, "/")
	}
	return parts[1:], true
}

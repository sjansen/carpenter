package patterns

import (
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/sjansen/carpenter/internal/sys"
)

type Patterns struct {
	tests map[string]result
	tree  tree
}

type pattern struct {
	id     string
	slash  slash
	prefix []part
	suffix *regexPart
	query  query
	tests  map[string]string
}

type result struct {
	id  string
	url string
}

type slash int

const (
	maySlash slash = iota
	mustSlash
	neverSlash
)

func (p *Patterns) Match(url *url.URL) (id, normalized string, err error) {
	results, err := p.match(url, false)
	if err != nil {
		return "", "", err
	} else if results == nil {
		return "", "", nil
	}
	return results[0].id, results[0].url, nil
}

func (p *Patterns) MatchAll(url *url.URL) (map[string]string, error) {
	results, err := p.match(url, true)
	if err != nil {
		return nil, err
	} else if results == nil {
		return nil, nil
	}

	matches := make(map[string]string, len(results))
	for _, result := range results {
		matches[result.id] = result.url
	}
	return matches, nil
}

func (p *Patterns) Test(sys *sys.IO) (map[string]string, error) {
	rawurls := make([]string, 0, len(p.tests))
	for rawurl := range p.tests {
		rawurls = append(rawurls, rawurl)
	}
	sort.Strings(rawurls)

	passed := make(map[string]string, len(p.tests))
	for _, rawurl := range rawurls {
		sys.Log.Debugf("testing url=%q", rawurl)
		url, err := url.Parse(rawurl)
		if err != nil {
			return nil, err
		}

		actual, err := p.match(url, true)
		if err != nil {
			return nil, err
		}

		expected := p.tests[rawurl]
		switch {
		default:
			passed[rawurl] = expected.id
		case len(actual) < 1 && expected.url == "":
			passed[rawurl] = ""
		case len(actual) < 1:
			return nil, fmt.Errorf(
				"url didn't match expected pattern: url=%q pattern=%q",
				rawurl, expected.id,
			)
		case len(actual) > 1:
			matches := make([]string, 0, len(actual))
			for _, result := range actual {
				matches = append(matches, result.id)
			}
			sort.Strings(matches)
			return nil, fmt.Errorf(
				"url matched by multiple patterns: url=%q patterns=%q",
				rawurl, matches,
			)
		case actual[0].id != expected.id:
			return nil, fmt.Errorf(
				"unexpected pattern matched: expected=%q actual=%q (url=%q)",
				expected.id, actual[0].id, rawurl,
			)
		case actual[0].url != expected.url:
			return nil, fmt.Errorf(
				"unexpected result: expected=%q actual=%q (pattern=%q url=%q)",
				expected.url, actual[0].url, expected.id, rawurl,
			)
		}
	}

	return passed, nil
}

func (p *Patterns) TestCases() map[string]string {
	testcases := make(map[string]string, len(p.tests))
	for rawurl, result := range p.tests {
		if result.url == "" {
			testcases[rawurl] = ""
		} else {
			testcases[rawurl] = result.id
		}
	}
	return testcases
}

func (p *Patterns) match(url *url.URL, matchAll bool) ([]*result, error) {
	if len(url.Path) < 1 || url.Path[0] != '/' {
		err := fmt.Errorf(`URLs must start with "/": %q`, url.Path)
		return nil, err
	}

	var results []*result
	if url.Path == "/" && p.tree.id != "" {
		matches, err := p.tree.recordMatch(1, url.Query())
		if err != nil {
			return nil, err
		}
		results = append(results, newResult(matches[0]))
		if !matchAll {
			return results, nil
		}
	}

	matches, err := p.tree.match(url.Path[1:], url.Query(), 0, matchAll)
	if err != nil {
		return nil, err
	}
	for _, match := range matches {
		results = append(results, newResult(match))
	}
	return results, nil
}

func newResult(m *match) *result {
	parts := append(m.parts, "")
	reverseParts(parts)
	normalized := strings.Join(parts, "/")
	if m.query != "" {
		normalized = normalized + "?" + m.query
	}
	return &result{m.id, normalized}
}

func reverseParts(parts []string) {
	for left, right := 0, len(parts)-1; left < right; left, right = left+1, right-1 {
		parts[left], parts[right] = parts[right], parts[left]
	}
}

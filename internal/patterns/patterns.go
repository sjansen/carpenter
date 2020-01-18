package patterns

import (
	"fmt"
	"net/url"
	"strings"
)

type Patterns struct {
	tests map[string]result
	tree  tree
}

type result struct {
	id  string
	url string
}

type pattern struct {
	id     string
	slash  slash
	prefix []part
	suffix *regexPart
	query  query
	tests  map[string]string
}

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

func (p *Patterns) match(url *url.URL, matchAll bool) ([]*result, error) {
	if len(url.Path) < 1 || url.Path[0] != '/' {
		err := fmt.Errorf(`URLs must start with "/": %q`, url.Path)
		return nil, err
	}

	var results []*result
	if url.Path == "/" && p.tree.id != "" {
		matches, err := p.tree.matchFound(1, url.Query())
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

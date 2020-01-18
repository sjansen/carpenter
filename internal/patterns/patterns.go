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
	params params
	tests  map[string]string
}

func (p *Patterns) Match(url *url.URL) (id, normalized string, err error) {
	if len(url.Path) < 1 || url.Path[0] != '/' {
		err := fmt.Errorf(`URLs must start with "/": %q`, url.Path)
		return "", "", err
	}

	var matches []*match
	if url.Path == "/" && p.tree.id != "" {
		matches, err = p.tree.matchFound(1, url.Query())
	} else {
		matches, err = p.tree.match(url.Path[1:], url.Query(), 0, false)
	}
	if err != nil {
		return "", "", err
	}

	if matches != nil {
		match := matches[0]
		parts := append(match.parts, "")
		reverseParts(parts)
		normalized := strings.Join(parts, "/")
		if match.query != "" {
			normalized = normalized + "?" + match.query
		}
		return match.id, normalized, nil
	}

	return "", "", nil
}

func reverseParts(parts []string) {
	for left, right := 0, len(parts)-1; left < right; left, right = left+1, right-1 {
		parts[left], parts[right] = parts[right], parts[left]
	}
}

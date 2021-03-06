package patterns

import (
	"net/url"
	"strings"

	"go.starlark.net/starlark"
)

type tree struct {
	id       string
	slash    slash
	query    query
	children []*child
}

type child struct {
	part part
	tree *tree
}

type match struct {
	id    string
	parts []string
	query string
}

func (t *tree) addPattern(p *pattern, depth int) {
	if depth < len(p.prefix) {
		c := &child{
			part: p.prefix[depth],
			tree: &tree{},
		}
		c.tree.addPattern(p, depth+1)
		t.children = append(t.children, c)
		return
	}

	if p.suffix != nil {
		c := &child{
			part: p.suffix,
			tree: &tree{
				id:    p.id,
				slash: p.slash,
				query: p.query,
			},
		}
		t.children = append(t.children, c)
		return
	}

	t.id = p.id
	t.slash = p.slash
	t.query = p.query
}

func (t *tree) match(path string, query url.Values, depth int, matchAll bool) ([]*match, error) {
	switch path {
	case "":
		switch {
		case t.id == "":
			return nil, nil
		case t.slash == maySlash:
			return t.recordMatch(depth, query)
		case t.slash == neverSlash:
			return t.recordMatch(depth, query)
		}
		return nil, nil
	case "/":
		switch {
		case t.id == "":
			return t.matchSuffix("", query, depth, matchAll)
		case t.slash == maySlash:
			return t.recordMatch(depth, query)
		case t.slash == mustSlash:
			return t.recordMatch(depth, query)
		}
		return nil, nil
	}

	return t.matchChildren(path, query, depth+1, matchAll)
}

func (t *tree) matchChildren(path string, query url.Values, depth int, matchAll bool) ([]*match, error) {
	var prefix, suffix string
	idx := strings.Index(path, "/")
	switch idx {
	case -1:
		prefix = path
	case len(path) - 1:
		prefix = path[:idx]
		suffix = path[idx:]
	default:
		prefix = path[:idx]
		suffix = path[idx+1:]
	}

	var matches []*match
	for _, child := range t.children {
		part := child.part

		prefix := prefix
		suffix := suffix
		if part.greedy() {
			prefix = path
			suffix = ""
		}

		if part.match(prefix) {
			tmp, err := child.tree.match(suffix, query, depth, matchAll)
			if err != nil {
				return nil, err
			}

			if tmp != nil {
				thread := &starlark.Thread{}
				normalized, err := child.part.normalize(thread, prefix)
				if err != nil {
					return nil, err
				}

				for _, m := range tmp {
					m.parts = append(m.parts, normalized)
				}

				if !matchAll {
					return tmp, nil
				}
				matches = append(matches, tmp...)
			}
		}
	}

	return matches, nil
}

func (t *tree) matchSuffix(path string, query url.Values, depth int, matchAll bool) ([]*match, error) {
	var matches []*match
	for _, child := range t.children {
		part := child.part
		if !part.greedy() {
			continue
		}

		if part.match(path) {
			tmp, err := t.recordMatch(depth, query)
			if err != nil {
				return nil, err
			}

			if tmp != nil {
				thread := &starlark.Thread{}
				normalized, err := child.part.normalize(thread, path)
				if err != nil {
					return nil, err
				}

				for _, m := range tmp {
					m.id = child.tree.id
					m.parts = append(m.parts, normalized)
				}

				if !matchAll {
					return tmp, nil
				}
				matches = append(matches, tmp...)
			}
		}
	}

	return matches, nil
}

func (t *tree) recordMatch(depth int, query url.Values) ([]*match, error) {
	q, err := t.rewriteQuery(query)
	if err != nil {
		return nil, err
	}

	m := match{id: t.id, query: q}
	if t.slash == mustSlash {
		m.parts = make([]string, 0, depth+2)
		m.parts = append(m.parts, "")
	} else {
		m.parts = make([]string, 0, depth+1)
	}

	return []*match{&m}, nil
}

func (t *tree) rewriteQuery(query url.Values) (string, error) {
	result := url.Values{}

	thread := &starlark.Thread{}
	for key, values := range query {
		param, ok := t.query.match[key]
		switch {
		case ok:
			values, err := param.normalize(thread, t.query.dedup, key, values)
			if err != nil {
				return "", err
			}
			result[key] = values
		case t.query.other != nil:
			values, err := t.query.other.normalize(thread, t.query.dedup, key, values)
			if err != nil {
				return "", err
			}
			result[key] = values
		default:
			result[key] = values
		}
	}

	return result.Encode(), nil
}

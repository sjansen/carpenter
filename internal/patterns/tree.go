package patterns

import (
	"net/url"
	"strings"

	"go.starlark.net/starlark"
)

type slash int

const (
	maySlash slash = iota
	mustSlash
	neverSlash
)

type tree struct {
	id       string
	slash    slash
	params   params
	children []*child
}

type child struct {
	part part
	tree *tree
}

type match struct {
	id    string
	parts []string
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
				id:     p.id,
				slash:  p.slash,
				params: p.params,
			},
		}
		t.children = append(t.children, c)
		return
	}

	t.id = p.id
	t.slash = p.slash
	t.params = p.params
}

func (t *tree) match(path string, query url.Values, depth int, matchAll bool) ([]*match, error) {
	switch path {
	case "":
		switch {
		case t.id == "":
			return nil, nil
		case t.slash == maySlash:
			parts := make([]string, 0, depth+1)
			return []*match{{t.id, parts}}, nil
		case t.slash == neverSlash:
			parts := make([]string, 0, depth+1)
			return []*match{{t.id, parts}}, nil
		}
		return nil, nil
	case "/":
		switch {
		case t.id == "":
			return nil, nil
		case t.slash == maySlash:
			parts := make([]string, 0, depth+1)
			return []*match{{t.id, parts}}, nil
		case t.slash == mustSlash:
			parts := make([]string, 0, depth+2)
			parts = append(parts, "")
			return []*match{{t.id, parts}}, nil
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
		if !part.split() {
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

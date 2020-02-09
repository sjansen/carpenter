package patterns

import "go.starlark.net/starlark"

type dedup int

const (
	keepAll dedup = iota
	keepFirst
	keepLast
)

type query struct {
	dedup dedup
	match map[string]*param
	other *param
}

type param struct {
	remove   bool
	rewriter queryRewriter
}

func (p *param) normalize(thread *starlark.Thread, dedup dedup, key string, values []string) ([]string, error) {
	if p.remove {
		return nil, nil
	} else if len(values) < 1 {
		return values, nil
	}

	result := make([]string, 0, len(values))
	for _, value := range values {
		v, err := p.rewriter.rewrite(thread, key, value)
		if err != nil {
			return nil, err
		}
		if v != nil {
			result = append(result, *v)
		}
	}
	if len(result) < 1 {
		return nil, nil
	}

	switch dedup {
	case keepFirst:
		return []string{result[0]}, nil
	case keepLast:
		return []string{result[len(result)-1]}, nil
	}

	return result, nil
}

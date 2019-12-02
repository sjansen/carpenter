package patterns

import "go.starlark.net/starlark"

type param struct {
	remove   bool
	rewriter rewriter
}

func (p *param) normalize(thread *starlark.Thread, dedup string, values []string) ([]string, error) {
	if p.remove {
		return nil, nil
	} else if len(values) < 1 {
		return values, nil
	}

	switch dedup {
	case "first":
		v, err := p.rewriter.rewrite(thread, values[0])
		if err != nil {
			return nil, err
		}
		return []string{v}, nil
	case "last":
		v, err := p.rewriter.rewrite(thread, values[len(values)-1])
		if err != nil {
			return nil, err
		}
		return []string{v}, nil
	}

	result := make([]string, 0, len(values))
	for _, v := range values {
		v, err := p.rewriter.rewrite(thread, v)
		if err != nil {
			return nil, err
		}
		result = append(result, v)
	}
	return result, nil
}

package rules

import "go.starlark.net/starlark"

type param struct {
	dedup    string
	remove   bool
	rewriter rewriter
}

func (p *param) normalize(thread *starlark.Thread, values []string) ([]string, error) {
	if p.remove {
		return nil, nil
	} else if len(values) < 1 {
		return values, nil
	}

	switch p.dedup {
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

package patterns

import (
	"fmt"

	"go.starlark.net/starlark"
)

type queryRewriter interface {
	rewrite(t *starlark.Thread, k, v string) (*string, error)
}

type stringRewriter interface {
	rewrite(t *starlark.Thread, s string) (string, error)
}

type callableQueryRewriter struct {
	starlark.Callable
}

func (r *callableQueryRewriter) rewrite(thread *starlark.Thread, param, value string) (*string, error) {
	args := starlark.Tuple{starlark.String(param), starlark.String(value)}
	result, err := starlark.Call(thread, r.Callable, args, nil)
	if err != nil {
		return nil, err
	}
	switch v := result.(type) {
	case starlark.String:
		s := v.GoString()
		return &(s), nil
	case starlark.NoneType:
		return nil, nil
	default:
		return nil, fmt.Errorf("%s: expected String, got %s", r.Callable, result.Type())
	}
}

type callableStringRewriter struct {
	starlark.Callable
}

func (r *callableStringRewriter) rewrite(thread *starlark.Thread, s string) (string, error) {
	args := starlark.Tuple{starlark.String(s)}
	result, err := starlark.Call(thread, r.Callable, args, nil)
	if err != nil {
		return "", err
	}
	tmp, ok := result.(starlark.String)
	if !ok {
		return "", fmt.Errorf("%s: expected String, got %s", r.Callable, result.Type())
	}
	return tmp.GoString(), nil
}

type staticQueryRewriter struct {
	value string
}

func (r *staticQueryRewriter) rewrite(_ *starlark.Thread, _, _ string) (*string, error) {
	return &r.value, nil
}

type staticStringRewriter struct {
	value string
}

func (r *staticStringRewriter) rewrite(_ *starlark.Thread, _ string) (string, error) {
	return r.value, nil
}

package rules

import (
	"fmt"

	"go.starlark.net/starlark"
)

type rewriter interface {
	rewrite(*starlark.Thread, string) (string, error)
}

type rewriteFunction struct {
	starlark.Callable
}

func (r *rewriteFunction) rewrite(thread *starlark.Thread, part string) (string, error) {
	args := starlark.Tuple{starlark.String(part)}
	value, err := starlark.Call(thread, r.Callable, args, nil)
	if err != nil {
		return "", err
	}
	s, ok := value.(starlark.String)
	if !ok {
		return "", fmt.Errorf("%s: expected String, got %s", r.Callable, value.Type())
	}
	return s.GoString(), nil
}

type rewriteStatic struct {
	value string
}

func (r *rewriteStatic) rewrite(thread *starlark.Thread, part string) (string, error) {
	return r.value, nil
}

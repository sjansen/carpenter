package rules

import (
	"fmt"

	"go.starlark.net/starlark"
)

type rewriter interface {
	rewrite(string) (string, error)
}

type rewriteFunction struct {
	thread *starlark.Thread
	fn     starlark.Callable
}

func (r *rewriteFunction) rewrite(part string) (string, error) {
	args := starlark.Tuple{starlark.String(part)}
	value, err := starlark.Call(r.thread, r.fn, args, nil)
	if err != nil {
		return "", err
	}
	s, ok := value.(starlark.String)
	if !ok {
		return "", fmt.Errorf("%s: expected String, got %s", r.fn, value.Type())
	}
	return s.GoString(), nil
}

type rewriteStatic struct {
	value string
}

func (r *rewriteStatic) rewrite(part string) (string, error) {
	return r.value, nil
}

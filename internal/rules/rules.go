package rules

import (
	"io"

	"go.starlark.net/starlark"
)

type Rules []*Rule

type Rule struct {
	id    string
	slash string
	parts []part
	tests map[string]string
}

func Load(filename string, src io.Reader) (Rules, error) {
	loader := &rulesLoader{
		rules: make([]*Rule, 0),
	}

	loader.Builtin = starlark.NewBuiltin("register_urls", loader.registerURLs)

	globals := starlark.StringDict{
		"register_urls": loader,
	}
	thread := &starlark.Thread{}
	_, err := starlark.ExecFile(thread, filename, src, globals)
	if err != nil {
		return nil, err
	}

	return loader.rules, nil
}

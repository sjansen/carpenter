package rules

import (
	"fmt"
	"io"
	"regexp"

	"go.starlark.net/starlark"
)

type Rule struct {
	id    string
	slash string
	parts []matcher
}

type Rules []*Rule

type matcher interface {
	match(string) bool
	rewrite(string) (string, error)
}

type ruleFactory struct {
	rules []*Rule
}

func Load(filename string, src io.Reader) (Rules, error) {
	f := &ruleFactory{
		rules: make([]*Rule, 0),
	}

	globals := starlark.StringDict{
		"register_urls": starlark.NewBuiltin("register_urls", f.registerURLs),
	}
	thread := &starlark.Thread{}
	_, err := starlark.ExecFile(thread, filename, src, globals)
	if err != nil {
		return nil, err
	}

	return f.rules, nil
}

func getIterableFromDict(fn string, rule *starlark.Dict, key string) (starlark.Iterable, error) {
	value, found, err := rule.Get(starlark.String(key))
	if err != nil || !found {
		return nil, err
	}
	iter, ok := value.(starlark.Iterable)
	if !ok {
		return nil, fmt.Errorf("%s: expected Iterable, got %s", fn, value.Type())
	}
	return iter, nil
}

func getStringFromDict(fn string, rule *starlark.Dict, key string) (string, error) {
	value, found, err := rule.Get(starlark.String(key))
	if err != nil || !found {
		return "", err
	}
	s, ok := value.(starlark.String)
	if !ok {
		return "", fmt.Errorf("%s: expected String, got %s", fn, value.Type())
	}
	return s.GoString(), nil
}

func (f *ruleFactory) registerURLs(
	thread *starlark.Thread,
	fn *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {
	fnName := fn.Name()
	if len(kwargs) > 0 {
		err := fmt.Errorf("%s: unexpected keyword argument", fnName)
		return nil, err
	}

	for _, arg := range args {
		if err := f.registerURL(fnName, arg); err != nil {
			return nil, err
		}
	}

	return starlark.None, nil
}

func (f *ruleFactory) registerURL(fn string, arg starlark.Value) error {
	d, ok := arg.(*starlark.Dict)
	if !ok {
		return fmt.Errorf("%s: expected Dict, got %s", fn, arg.Type())
	}

	id, err := getStringFromDict(fn, d, "id")
	if err != nil {
		return err
	} else if id == "" {
		return fmt.Errorf(`%s: missing required key: "id"`, fn)
	}

	parts, err := getIterableFromDict(fn, d, "parts")
	if err != nil {
		return err
	} else if parts == nil {
		return fmt.Errorf(`%s: missing required key: "parts"`, fn)
	}

	slash, err := getStringFromDict(fn, d, "slash")
	if err != nil {
		return err
	} else if slash == "" {
		return fmt.Errorf(`%s: missing required key: "slash"`, fn)
	}

	rule := &Rule{
		id:    id,
		slash: slash,
		parts: make([]matcher, 0),
	}

	iter := parts.Iterate()
	defer iter.Done()

	var value starlark.Value
	for iter.Next(&value) {
		m, err := transformPart(fn, value)
		if err != nil {
			return err
		}
		rule.parts = append(rule.parts, m)
	}

	f.rules = append(f.rules, rule)
	return nil
}

func transformPart(fn string, value starlark.Value) (matcher, error) {
	switch v := value.(type) {
	case starlark.String:
		m := &plainPart{
			value: v.GoString(),
		}
		return m, nil
	case starlark.Tuple:
		if n := v.Len(); n != 2 {
			return nil, fmt.Errorf("%s: expected 2 item Tuple, got %d", fn, n)
		}

		m := &regexPart{}

		var value starlark.Value
		iter := v.Iterate()
		defer iter.Done()

		iter.Next(&value)
		expr, ok := value.(starlark.String)
		if !ok {
			return nil, fmt.Errorf("%s: expected String, got %s", fn, value.Type())
		}

		regex, err := regexp.Compile(expr.GoString())
		if err != nil {
			return nil, err
		}
		m.regex = regex

		iter.Next(&value)
		replacement, ok := value.(starlark.String)
		if !ok {
			return nil, fmt.Errorf("%s: expected String, got %s", fn, value.Type())
		}

		m.replacement = &plainReplacement{
			value: replacement.GoString(),
		}
		return m, nil
	}
	return nil, fmt.Errorf("%s: expected String or Tuple, got %s", fn, value.Type())
}

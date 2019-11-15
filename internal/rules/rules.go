package rules

import (
	"fmt"
	"io"

	"go.starlark.net/starlark"
)

type Rules struct {
	Matchers []Matcher
}

func Load(filename string, src io.Reader) (*Rules, error) {
	rules := &Rules{
		Matchers: make([]Matcher, 0),
	}

	globals := starlark.StringDict{
		"register_urls": starlark.NewBuiltin("register_urls", rules.registerURLs),
	}
	thread := &starlark.Thread{}
	_, err := starlark.ExecFile(thread, filename, src, globals)
	if err != nil {
		return nil, err
	}

	return rules, nil
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

func (r *Rules) registerURLs(
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
		if err := r.registerURL(fnName, arg); err != nil {
			return nil, err
		}
	}

	return starlark.None, nil
}

func (r *Rules) registerURL(fn string, arg starlark.Value) error {
	rule, ok := arg.(*starlark.Dict)
	if !ok {
		return fmt.Errorf("%s: expected Dict, got %s", fn, arg.Type())
	}

	id, err := getStringFromDict(fn, rule, "id")
	if err != nil {
		return err
	} else if id == "" {
		return fmt.Errorf(`%s: missing required key: "id"`, fn)
	}

	parts, err := getIterableFromDict(fn, rule, "parts")
	if err != nil {
		return err
	} else if parts == nil {
		return fmt.Errorf(`%s: missing required key: "parts"`, fn)
	}

	slash, err := getStringFromDict(fn, rule, "slash")
	if err != nil {
		return err
	} else if slash == "" {
		return fmt.Errorf(`%s: missing required key: "slash"`, fn)
	}

	m := &PlainPart{
		id:    id,
		slash: slash,
	}

	iter := parts.Iterate()
	defer iter.Done()
	var value starlark.Value
	iter.Next(&value)
	s, ok := value.(starlark.String)
	if !ok {
		return fmt.Errorf("%s: expected String, got %s", fn, value.Type())
	}
	m.value = s.GoString()

	r.Matchers = append(r.Matchers, m)

	return nil
}

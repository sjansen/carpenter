package rules

import (
	"fmt"
	"regexp"

	"go.starlark.net/starlark"
)

type rulesLoader struct {
	*starlark.Builtin

	rules []*Rule
}

func (l *rulesLoader) getIterableFromDict(rule *starlark.Dict, key string) (starlark.Iterable, error) {
	value, found, err := rule.Get(starlark.String(key))
	if err != nil || !found {
		return nil, err
	}
	iter, ok := value.(starlark.Iterable)
	if !ok {
		return nil, fmt.Errorf("%s: expected Iterable, got %s", l.Name(), value.Type())
	}
	return iter, nil
}

func (l *rulesLoader) getStringFromDict(rule *starlark.Dict, key string) (string, error) {
	value, found, err := rule.Get(starlark.String(key))
	if err != nil || !found {
		return "", err
	}
	s, ok := value.(starlark.String)
	if !ok {
		return "", fmt.Errorf("%s: expected String, got %s", l.Name(), value.Type())
	}
	return s.GoString(), nil
}

func (l *rulesLoader) registerURL(arg starlark.Value) error {
	d, ok := arg.(*starlark.Dict)
	if !ok {
		return fmt.Errorf("%s: expected Dict, got %s", l.Name(), arg.Type())
	}

	id, err := l.getStringFromDict(d, "id")
	if err != nil {
		return err
	} else if id == "" {
		return fmt.Errorf(`%s: missing required key: "id"`, l.Name())
	}

	parts, err := l.getIterableFromDict(d, "parts")
	if err != nil {
		return err
	} else if parts == nil {
		return fmt.Errorf(`%s: missing required key: "parts"`, l.Name())
	}

	slash, err := l.getStringFromDict(d, "slash")
	if err != nil {
		return err
	} else if slash == "" {
		return fmt.Errorf(`%s: missing required key: "slash"`, l.Name())
	}

	rule := &Rule{
		id:    id,
		slash: slash,
		parts: make([]part, 0),
	}

	iter := parts.Iterate()
	defer iter.Done()

	var value starlark.Value
	for iter.Next(&value) {
		m, err := l.transformPart(value)
		if err != nil {
			return err
		}
		rule.parts = append(rule.parts, m)
	}

	l.rules = append(l.rules, rule)
	return nil
}

func (l *rulesLoader) registerURLs(
	thread *starlark.Thread,
	fn *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {
	if len(kwargs) > 0 {
		err := fmt.Errorf("%s: unexpected keyword argument", l.Name())
		return nil, err
	}

	for _, arg := range args {
		if err := l.registerURL(arg); err != nil {
			return nil, err
		}
	}

	return starlark.None, nil
}

func (l *rulesLoader) transformPart(value starlark.Value) (part, error) {
	switch v := value.(type) {
	case starlark.String:
		m := &plainPart{
			value: v.GoString(),
		}
		return m, nil
	case starlark.Tuple:
		if n := v.Len(); n != 2 {
			return nil, fmt.Errorf("%s: expected 2 item Tuple, got %d", l.Name(), n)
		}

		m := &regexPart{}

		var value starlark.Value
		iter := v.Iterate()
		defer iter.Done()

		iter.Next(&value)
		expr, ok := value.(starlark.String)
		if !ok {
			return nil, fmt.Errorf("%s: expected String, got %s", l.Name(), value.Type())
		}

		regex, err := regexp.Compile(expr.GoString())
		if err != nil {
			return nil, err
		}
		m.regex = regex

		iter.Next(&value)
		switch v := value.(type) {
		case starlark.Callable:
			m.rewriter = &rewriteFunction{
				Callable: v,
			}
			return m, nil
		case starlark.String:
			m.rewriter = &rewriteStatic{
				value: v.GoString(),
			}
			return m, nil
		}
		return nil, fmt.Errorf("%s: expected Callable or String, got %s", l.Name(), value.Type())
	}
	return nil, fmt.Errorf("%s: expected String or Tuple, got %s", l.Name(), value.Type())
}

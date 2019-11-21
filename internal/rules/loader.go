package rules

import (
	"fmt"
	"regexp"

	"go.starlark.net/resolve"
	"go.starlark.net/starlark"
)

func init() {
	resolve.AllowLambda = true
}

type rulesLoader struct {
	*starlark.Builtin

	rules []*Rule
}

func (l *rulesLoader) getDictFromDict(id, key string, rule *starlark.Dict) (*starlark.Dict, error) {
	value, found, err := rule.Get(starlark.String(key))
	if err != nil || !found {
		return nil, err
	}
	dict, ok := value.(*starlark.Dict)
	if !ok {
		return nil, fmt.Errorf("%s: %q expected Dict, got %s", l.Name(), id, value.Type())
	}
	return dict, nil
}

func (l *rulesLoader) getIterableFromDict(id, key string, rule *starlark.Dict) (starlark.Iterable, error) {
	value, found, err := rule.Get(starlark.String(key))
	if err != nil || !found {
		return nil, err
	}
	iter, ok := value.(starlark.Iterable)
	if !ok {
		return nil, fmt.Errorf("%s: %q expected Iterable, got %s", l.Name(), id, value.Type())
	}
	return iter, nil
}

func (l *rulesLoader) getStringFromDict(id, key string, rule *starlark.Dict) (string, error) {
	value, found, err := rule.Get(starlark.String(key))
	if err != nil || !found {
		return "", err
	}
	s, ok := value.(starlark.String)
	if !ok {
		if id == "" {
			return "", fmt.Errorf("%s: expected String, got %s", l.Name(), value.Type())
		}
		return "", fmt.Errorf("%s: %q expected String, got %s", l.Name(), id, value.Type())
	}
	return s.GoString(), nil
}

func (l *rulesLoader) registerURL(arg starlark.Value) error {
	d, ok := arg.(*starlark.Dict)
	if !ok {
		return fmt.Errorf("%s: expected Dict, got %s", l.Name(), arg.Type())
	}

	id, err := l.getStringFromDict("", "id", d)
	if err != nil {
		return err
	} else if id == "" {
		return fmt.Errorf(`%s: missing required key: "id"`, l.Name())
	}

	path, err := l.getDictFromDict(id, "path", d)
	if err != nil {
		return err
	} else if path == nil {
		return fmt.Errorf(`%s: %q missing required key: "path"`, l.Name(), id)
	}

	parts, err := l.getIterableFromDict(id, "parts", path)
	if err != nil {
		return err
	} else if parts == nil {
		return fmt.Errorf(`%s: %q missing required key: "parts"`, l.Name(), id)
	}

	slash, err := l.getStringFromDict(id, "slash", path)
	if err != nil {
		return err
	} else if slash == "" {
		return fmt.Errorf(`%s: %q missing required key: "slash"`, l.Name(), id)
	}

	query, err := l.getDictFromDict(id, "query", d)
	if err != nil {
		return err
	} else if query == nil {
		return fmt.Errorf(`%s: %q missing required key: "query"`, l.Name(), id)
	}

	tests, err := l.getDictFromDict(id, "tests", d)
	if err != nil {
		return err
	} else if tests == nil {
		return fmt.Errorf(`%s: %q missing required key: "tests"`, l.Name(), id)
	}

	rule := &Rule{
		id:    id,
		slash: slash,
	}

	rule.parts, err = l.transformParts(id, parts)
	if err != nil {
		return err
	}

	err = l.transformQuery(rule, query)
	if err != nil {
		return err
	}

	rule.tests, err = l.transformTests(id, tests)
	if err != nil {
		return err
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

func (l *rulesLoader) transformParams(r *Rule, query *starlark.Dict) error {
	result := make(map[string]*param)

	for _, item := range query.Items() {
		key := item.Index(0)
		k, ok := key.(starlark.String)
		if !ok {
			return fmt.Errorf(
				"%s: %q expected String key, got %s",
				l.Name(), r.id, key.Type(),
			)
		}

		value := item.Index(1)
		switch v := value.(type) {
		case starlark.Callable:
			result[k.GoString()] = &param{
				rewriter: &rewriteFunction{
					Callable: v,
				},
			}
		case starlark.NoneType:
			result[k.GoString()] = &param{
				remove: true,
			}
		case starlark.String:
			result[k.GoString()] = &param{
				rewriter: &rewriteStatic{
					value: v.GoString(),
				},
			}
		default:
			return fmt.Errorf(
				"%s: %q expected None or String value, got %s",
				l.Name(), r.id, value.Type(),
			)
		}
	}

	r.params = result
	return nil
}

func (l *rulesLoader) transformParts(id string, parts starlark.Iterable) ([]part, error) {
	result := make([]part, 0)

	iter := parts.Iterate()
	defer iter.Done()

	var value starlark.Value
	for iter.Next(&value) {
		part, err := l.transformPart(id, value)
		if err != nil {
			return nil, err
		}
		result = append(result, part)
	}

	return result, nil
}

func (l *rulesLoader) transformPart(id string, value starlark.Value) (part, error) {
	switch v := value.(type) {
	case starlark.String:
		m := &plainPart{
			value: v.GoString(),
		}
		return m, nil
	case starlark.Tuple:
		if n := v.Len(); n != 2 {
			return nil, fmt.Errorf("%s: %q expected 2 item Tuple, got %d", l.Name(), id, n)
		}

		m := &regexPart{}

		var value starlark.Value
		iter := v.Iterate()
		defer iter.Done()

		iter.Next(&value)
		expr, ok := value.(starlark.String)
		if !ok {
			return nil, fmt.Errorf("%s: %q expected String, got %s", l.Name(), id, value.Type())
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
		return nil, fmt.Errorf("%s: %q expected Callable or String, got %s", l.Name(), id, value.Type())
	}
	return nil, fmt.Errorf("%s: %q expected String or Tuple, got %s", l.Name(), id, value.Type())
}

func (l *rulesLoader) transformQuery(r *Rule, query *starlark.Dict) error {
	for _, item := range query.Items() {
		key := item.Index(0)
		k, ok := key.(starlark.String)
		if !ok {
			return fmt.Errorf(
				"%s: %q expected String key, got %s",
				l.Name(), r.id, key.Type(),
			)
		}

		switch k {
		case "dedup":
			value := item.Index(1)
			v, ok := value.(starlark.String)
			if !ok {
				return fmt.Errorf(
					"%s: %q expected String, got %s",
					l.Name(), r.id, value.Type(),
				)
			}
			r.dedup = v.GoString()
		case "params":
			value := item.Index(1)
			v, ok := value.(*starlark.Dict)
			if !ok {
				return fmt.Errorf(
					"%s: %q expected Dict, got %s",
					l.Name(), r.id, value.Type(),
				)
			}
			if err := l.transformParams(r, v); err != nil {
				return err
			}
		default:
			return fmt.Errorf(
				`%s: %q expected "dedup" or "params", got %s`,
				l.Name(), r.id, k,
			)
		}
	}

	return nil
}

func (l *rulesLoader) transformTests(id string, tests *starlark.Dict) (map[string]string, error) {
	result := make(map[string]string)

	for _, item := range tests.Items() {
		key := item.Index(0)
		k, ok := key.(starlark.String)
		if !ok {
			return nil, fmt.Errorf(
				"%s: %q expected String key, got %s",
				l.Name(), id, key.Type(),
			)
		}

		value := item.Index(1)
		switch v := value.(type) {
		case starlark.NoneType:
			result[k.GoString()] = ""
		case starlark.String:
			result[k.GoString()] = v.GoString()
		default:
			return nil, fmt.Errorf(
				"%s: %q expected None or String value, got %s",
				l.Name(), id, value.Type(),
			)
		}
	}

	return result, nil
}

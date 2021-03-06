package patterns

import (
	"fmt"
	"io"
	"regexp"

	"go.starlark.net/resolve"
	"go.starlark.net/starlark"
)

func init() {
	resolve.AllowLambda = true
}

type patternLoader struct {
	*starlark.Builtin // enable l.Name()

	globals  starlark.StringDict
	patterns []*pattern
	rename   *starlark.Function
}

func Load(filename string, src io.Reader) (*Patterns, error) {
	loader, err := loadPatterns(filename, src)
	if err != nil {
		return nil, err
	}

	patterns := &Patterns{
		rename: loader.rename,
		tests:  map[string]result{},
	}
	for _, p := range loader.patterns {
		patterns.tree.addPattern(p, 0)
		for raw, expected := range p.tests {
			other, ok := patterns.tests[raw]
			switch {
			case ok && expected != "" && other.url != "":
				err := fmt.Errorf(
					"test case repeated: %q (original=%q) (conflict=%q)",
					raw, other.id, p.id,
				)
				return nil, err
			case ok && expected == "":
				continue
			case ok && other.url == "":
				fallthrough
			default:
				patterns.tests[raw] = result{
					id:  p.id,
					url: expected,
				}
			}
		}
	}

	return patterns, nil
}

func loadPatterns(filename string, src io.Reader) (*patternLoader, error) {
	loader := newLoader()
	thread := &starlark.Thread{}
	_, err := starlark.ExecFile(thread, filename, src, loader.globals)
	if err != nil {
		return nil, err
	}
	return loader, nil
}

func newLoader() *patternLoader {
	loader := &patternLoader{
		patterns: make([]*pattern, 0),
	}

	loader.Builtin = starlark.NewBuiltin("url", loader.addURL)

	loader.globals = starlark.StringDict{
		"set_rename_filter": starlark.NewBuiltin(
			"set_rename_filter", loader.setRenameFilter,
		),
		"url": loader,
	}

	return loader
}

func (l *patternLoader) addURL(
	_ *starlark.Thread,
	fn *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {
	var id string
	var path, query, tests *starlark.Dict
	if err := starlark.UnpackArgs(
		fn.Name(), args, kwargs,
		"id", &id, "path", &path, "query", &query, "tests", &tests,
	); err != nil {
		return nil, err
	}

	if id == "" {
		return nil, fmt.Errorf(`%s: invalid pattern ID: ""`, fn.Name())
	}

	p := &pattern{id: id}
	if err := l.transformPath(p, path); err != nil {
		return nil, err
	}
	if err := l.transformQuery(p, query); err != nil {
		return nil, err
	}
	if err := l.transformTests(p, tests); err != nil {
		return nil, err
	}

	l.patterns = append(l.patterns, p)
	return starlark.None, nil
}

func (l *patternLoader) getIterableFromDict(parent, key string, pattern *starlark.Dict) (starlark.Iterable, error) {
	value, found, err := pattern.Get(starlark.String(key))
	if err != nil || !found {
		return nil, err
	}
	iter, ok := value.(starlark.Iterable)
	if !ok {
		return nil, fmt.Errorf("%s: %q/%q expected Iterable, got %s", l.Name(), parent, key, value.Type())
	}
	return iter, nil
}

func (l *patternLoader) setRenameFilter(
	_ *starlark.Thread,
	fn *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {
	var rename *starlark.Function
	if err := starlark.UnpackArgs(
		fn.Name(), args, kwargs, "fn", &rename,
	); err != nil {
		return nil, err
	}
	l.rename = rename
	return starlark.None, nil
}

func (l *patternLoader) transformPart(parent, child string, value starlark.Value) (part, error) {
	switch v := value.(type) {
	case starlark.String:
		value := v.GoString()
		if value == "" {
			return nil, fmt.Errorf(`%s: %q/%q invalid value: ""`, l.Name(), parent, child)
		}
		m := &plainPart{
			value: value,
		}
		return m, nil
	case starlark.Tuple:
		if n := v.Len(); n < 2 || n > 3 {
			return nil, fmt.Errorf(
				"%s: %q/%q expected 2 or 3 item Tuple, got %d",
				l.Name(), parent, child, n,
			)
		}

		m := &regexPart{}

		var value starlark.Value
		iter := v.Iterate()
		defer iter.Done()

		iter.Next(&value)
		expr, ok := value.(starlark.String)
		if !ok {
			return nil, fmt.Errorf(
				"%s: %q/%q expected String, got %s", l.Name(), parent, child, value.Type(),
			)
		}

		regex, err := regexp.Compile(expr.GoString())
		if err != nil {
			return nil, err
		}
		m.regex = regex

		iter.Next(&value)
		switch v := value.(type) {
		case starlark.Callable:
			m.rewriter = &callableStringRewriter{
				Callable: v,
			}
		case starlark.String:
			m.rewriter = &staticStringRewriter{
				value: v.GoString(),
			}
		default:
			return nil, fmt.Errorf(
				"%s: %q/%q expected Callable or String, got %s",
				l.Name(), parent, child, value.Type(),
			)
		}

		if iter.Next(&value) {
			expr, ok := value.(starlark.String)
			if !ok {
				return nil, fmt.Errorf(
					"%s: %q/%q expected String, got %s", l.Name(), parent, child, value.Type(),
				)
			}

			regex, err := regexp.Compile(expr.GoString())
			if err != nil {
				return nil, err
			}
			m.reject = regex
		}

		return m, nil
	}
	return nil, fmt.Errorf("%s: %q/%q expected String or Tuple, got %s", l.Name(), parent, child, value.Type())
}

func (l *patternLoader) transformPath(p *pattern, path *starlark.Dict) error {
	prefix, err := l.getIterableFromDict(p.id, "prefix", path)
	if err != nil {
		return err
	} else if prefix == nil {
		return fmt.Errorf(`%s: %q missing required key: "prefix"`, l.Name(), p.id)
	}

	p.prefix, err = l.transformPrefix(p.id, prefix)
	if err != nil {
		return err
	}

	value, found, err := path.Get(starlark.String("suffix"))
	if err != nil {
		return err
	} else if !found {
		return fmt.Errorf(`%s: %q missing required key: "suffix"`, l.Name(), p.id)
	}

	if s, ok := value.(starlark.String); ok {
		suffix := s.GoString()
		switch suffix {
		case "/?":
			p.slash = maySlash
		case "/":
			p.slash = mustSlash
		case "":
			p.slash = neverSlash
		default:
			return fmt.Errorf(`%s: %q/"suffix" invalid value: %q`, l.Name(), p.id, suffix)
		}
	} else if part, err := l.transformPart(p.id, "suffix", value); err != nil {
		return err
	} else {
		p.slash = neverSlash
		p.suffix = part.(*regexPart)
		p.suffix.suffix = true
	}

	if len(p.prefix) == 0 && p.suffix == nil && p.slash != mustSlash {
		// requiring suffix="/" simplifies normalization of queries when path="/"
		return fmt.Errorf(`suffix must be "/" or regex when no prefix parts are declared`)
	}

	return nil
}

func (l *patternLoader) transformPrefix(id string, prefix starlark.Iterable) ([]part, error) {
	result := make([]part, 0)

	iter := prefix.Iterate()
	defer iter.Done()

	var value starlark.Value
	for iter.Next(&value) {
		part, err := l.transformPart(id, "prefix", value)
		if err != nil {
			return nil, err
		}
		result = append(result, part)
	}

	return result, nil
}

func (l *patternLoader) transformQuery(p *pattern, query *starlark.Dict) error {
	for _, item := range query.Items() {
		key := item.Index(0)
		k, ok := key.(starlark.String)
		if !ok {
			return fmt.Errorf(
				"%s: %q expected String key, got %s",
				l.Name(), p.id, key.Type(),
			)
		}

		switch k {
		case "dedup":
			if err := l.transformQueryDedup(item, p, k); err != nil {
				return err
			}
		case "match":
			if err := l.transformQueryMatch(item, p); err != nil {
				return err
			}
		case "other":
			if err := l.transformQueryOther(item, p); err != nil {
				return err
			}
		default:
			return fmt.Errorf(
				`%s: %q expected "dedup", "match" or "other", got %q`,
				l.Name(), p.id, k.GoString(),
			)
		}
	}

	return nil
}

func (l *patternLoader) transformQueryDedup(item starlark.Tuple, p *pattern, k starlark.String) error {
	value := item.Index(1)
	v, ok := value.(starlark.String)
	if !ok {
		return fmt.Errorf(
			"%s: %q expected String, got %s",
			l.Name(), p.id, value.Type(),
		)
	}
	dedup := v.GoString()
	switch dedup {
	case "never":
		p.query.dedup = keepAll
	case "first":
		p.query.dedup = keepFirst
	case "last":
		p.query.dedup = keepLast
	default:
		return fmt.Errorf(
			`%s: %q/"query"/%q/"dedup" invalid value: %q`,
			l.Name(), p.id, k, dedup,
		)
	}
	return nil
}

func (l *patternLoader) transformQueryMatch(item starlark.Tuple, p *pattern) error {
	value := item.Index(1)
	query, ok := value.(*starlark.Dict)
	if !ok {
		return fmt.Errorf(
			"%s: %q expected Dict, got %s",
			l.Name(), p.id, value.Type(),
		)
	}

	result := make(map[string]*param)
	for _, item := range query.Items() {
		key := item.Index(0)
		k, ok := key.(starlark.String)
		if !ok {
			return fmt.Errorf(
				"%s: %q expected String key, got %s",
				l.Name(), p.id, key.Type(),
			)
		}

		value := item.Index(1)
		switch v := value.(type) {
		case starlark.Callable:
			result[k.GoString()] = &param{
				rewriter: &callableQueryRewriter{
					Callable: v,
				},
			}
		case starlark.NoneType:
			result[k.GoString()] = &param{
				remove: true,
			}
		case starlark.String:
			key := k.GoString()
			if key == "" {
				return fmt.Errorf(`%s: %q invalid query key: ""`, l.Name(), p.id)
			}
			result[key] = &param{
				rewriter: &staticQueryRewriter{
					value: v.GoString(),
				},
			}
		default:
			return fmt.Errorf(
				"%s: %q expected None, String, or Callable value, got %s",
				l.Name(), p.id, value.Type(),
			)
		}
	}

	p.query.match = result
	return nil
}

func (l *patternLoader) transformQueryOther(item starlark.Tuple, p *pattern) error {
	value := item.Index(1)
	switch v := value.(type) {
	case starlark.Callable:
		p.query.other = &param{
			rewriter: &callableQueryRewriter{
				Callable: v,
			},
		}
	case starlark.NoneType:
		p.query.other = &param{
			remove: true,
		}
	case starlark.String:
		p.query.other = &param{
			rewriter: &staticQueryRewriter{
				value: v.GoString(),
			},
		}
	default:
		return fmt.Errorf(
			"%s: %q expected None, String, or Callable value, got %s",
			l.Name(), p.id, value.Type(),
		)
	}
	return nil
}

func (l *patternLoader) transformTests(p *pattern, tests *starlark.Dict) error {
	result := make(map[string]string)

	for _, item := range tests.Items() {
		tmp := item.Index(0)
		key, ok := tmp.(starlark.String)
		if !ok {
			return fmt.Errorf(
				"%s: %q expected String key, got %s",
				l.Name(), p.id, tmp.Type(),
			)
		}

		k := key.GoString()
		if len(k) < 1 || k[0] != '/' {
			return fmt.Errorf(`invalid test case: %q (should start with "/")`, k)
		}

		tmp = item.Index(1)
		switch value := tmp.(type) {
		case starlark.NoneType:
			result[k] = ""
		case starlark.String:
			result[k] = value.GoString()
		default:
			return fmt.Errorf(
				"%s: %q expected None or String value, got %s",
				l.Name(), p.id, tmp.Type(),
			)
		}
	}

	p.tests = result
	return nil
}

package rules

import (
	"fmt"
	"regexp"

	"go.starlark.net/starlark"
)

type replacement interface {
	Replace(string) (string, error)
}

type plainPart struct {
	value string
}

func (p *plainPart) match(part string) bool {
	return part == p.value
}

func (p *plainPart) rewrite(part string) (string, error) {
	return part, nil
}

type regexPart struct {
	regex       *regexp.Regexp
	replacement replacement
}

func (p *regexPart) match(part string) bool {
	return p.regex.MatchString(part)
}

func (p *regexPart) rewrite(part string) (string, error) {
	return p.replacement.Replace(part)
}

type plainReplacement struct {
	value string
}

func (r *plainReplacement) Replace(part string) (string, error) {
	return r.value, nil
}

type functionReplacement struct {
	thread *starlark.Thread
	fn     starlark.Callable
}

func (r *functionReplacement) Replace(part string) (string, error) {
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

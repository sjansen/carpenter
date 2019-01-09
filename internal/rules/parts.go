package rules

import (
	"fmt"
	"regexp"

	"go.starlark.net/starlark"
)

type Matcher interface {
	Match(string) (string, bool)
	Rewrite(string) (string, error)
}

type Replacement interface {
	Replace(string) (string, error)
}

type PlainPart struct {
	id    string
	slash string
	value string
}

func (p *PlainPart) Match(part string) (string, bool) {
	if part == p.value {
		return p.id, true
	}
	return "", false
}

func (p *PlainPart) Rewrite(part string) (string, error) {
	return part, nil
}

type RegexPart struct {
	id          string
	slash       string
	regex       *regexp.Regexp
	replacement Replacement
}

func (p *RegexPart) Match(part string) (string, bool) {
	if p.regex.MatchString(part) {
		return p.id, true
	}
	return "", false
}

func (p *RegexPart) Rewrite(part string) (string, error) {
	return p.replacement.Replace(part)
}

type PlainReplacement struct {
	value string
}

func (r *PlainReplacement) Replace(part string) (string, error) {
	return r.value, nil
}

type FunctionReplacement struct {
	thread *starlark.Thread
	fn     starlark.Callable
}

func (r *FunctionReplacement) Replace(part string) (string, error) {
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

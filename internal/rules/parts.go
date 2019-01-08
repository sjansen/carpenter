package rules

import (
	"regexp"

	"go.starlark.net/starlark"
)

type Matcher interface {
	Match(string) (string, bool)
	Rewrite(string) string
}

type Replacement interface {
	Replace(string) string
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

func (p *PlainPart) Rewrite(part string) string {
	return part
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

func (p *RegexPart) Rewrite(part string) string {
	return p.replacement.Replace(part)
}

type PlainReplacement struct {
	value string
}

func (r *PlainReplacement) Replace(part string) string {
	return r.value
}

type FunctionReplacement struct {
	fn starlark.Callable
}

func (r *FunctionReplacement) Replace(part string) string {
	return part // TODO
}

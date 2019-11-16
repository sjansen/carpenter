package rules

import (
	"regexp"

	"go.starlark.net/starlark"
)

type part interface {
	match(string) bool
	normalize(*starlark.Thread, string) (string, error)
}

type plainPart struct {
	value string
}

func (p *plainPart) match(part string) bool {
	return part == p.value
}

func (p *plainPart) normalize(thread *starlark.Thread, part string) (string, error) {
	return part, nil
}

type regexPart struct {
	regex    *regexp.Regexp
	rewriter rewriter
}

func (p *regexPart) match(part string) bool {
	return p.regex.MatchString(part)
}

func (p *regexPart) normalize(thread *starlark.Thread, part string) (string, error) {
	return p.rewriter.rewrite(thread, part)
}

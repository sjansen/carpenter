package patterns

import (
	"regexp"

	"go.starlark.net/starlark"
)

type part interface {
	match(string) bool
	normalize(*starlark.Thread, string) (string, error)
	split() bool
}

type plainPart struct {
	value string
}

func (p *plainPart) match(path string) bool {
	return path == p.value
}

func (p *plainPart) normalize(_ *starlark.Thread, path string) (string, error) {
	return path, nil
}

func (p *plainPart) split() bool {
	return true
}

type regexPart struct {
	regex    *regexp.Regexp
	rewriter rewriter
	suffix   bool
}

func (p *regexPart) match(path string) bool {
	return p.regex.MatchString(path)
}

func (p *regexPart) normalize(thread *starlark.Thread, path string) (string, error) {
	return p.rewriter.rewrite(thread, path)
}

func (p *regexPart) split() bool {
	return !p.suffix
}

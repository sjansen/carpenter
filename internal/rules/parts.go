package rules

import (
	"regexp"
)

type part interface {
	match(string) bool
	normalize(string) (string, error)
}

type plainPart struct {
	value string
}

func (p *plainPart) match(part string) bool {
	return part == p.value
}

func (p *plainPart) normalize(part string) (string, error) {
	return part, nil
}

type regexPart struct {
	regex    *regexp.Regexp
	rewriter rewriter
}

func (p *regexPart) match(part string) bool {
	return p.regex.MatchString(part)
}

func (p *regexPart) normalize(part string) (string, error) {
	return p.rewriter.rewrite(part)
}

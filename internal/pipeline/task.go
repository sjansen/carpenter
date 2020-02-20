package pipeline

import (
	"bufio"
	"io"
	"net/url"
	"sort"
	"strings"

	"github.com/sjansen/carpenter/internal/lazyio"
	"github.com/sjansen/carpenter/internal/patterns"
	"github.com/sjansen/carpenter/internal/tokenizer"
	"github.com/sjansen/carpenter/internal/uaparser"
)

type Task struct {
	id        string
	patterns  *patterns.Patterns
	tokenizer *tokenizer.Tokenizer
	uaparser  *uaparser.Parser

	src   *lazyio.Input
	dst   lazyio.CSV
	debug debug
}

type debug struct {
	normalize    lazyio.CSV
	parse        lazyio.CSV
	tokenize     lazyio.TXT
	unrecognized lazyio.CSV
}

func (d *debug) Close() {
	d.normalize.Close()
	d.parse.Close()
	d.tokenize.Close()
	d.unrecognized.Close()
}

func (t *Task) Run() error {
	r, err := t.src.Open()
	if err != nil {
		return err
	}
	defer t.src.Close()
	defer t.dst.Close()
	defer t.debug.Close()

	buf := bufio.NewReader(r)
	var cols []string
	var row []string
	for {
		line, err := buf.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		tokens := t.tokenizer.Tokenize(line)
		if tokens == nil {
			t.debug.tokenize.Write(line)
			t.debug.tokenize.Write("\n")
			continue
		} else if cols == nil {
			cols = t.newCols(tokens)
			t.dst.Write(cols...)
			row = make([]string, len(cols))
		}

		t.parseURL(tokens)

		t.parseUserAgent(tokens)

		for i, k := range cols {
			if v, ok := tokens[k]; ok {
				row[i] = v
			} else {
				row[i] = ""
			}
		}
		t.dst.Write(row...)
	}

	t.dst.Flush()
	return t.dst.Error()
}

func (t *Task) newCols(tokens map[string]string) []string {
	cols := make([]string, 0, len(tokens)+11)
	cols = append(cols,
		"normalized_url",
		"url_pattern",
	)
	if t.uaparser != nil {
		cols = append(cols,
			"client_device_family",
			"client_os_family",
			"client_os_major",
			"client_os_minor",
			"client_os_patch",
			"client_ua_family",
			"client_ua_major",
			"client_ua_minor",
			"client_ua_patch",
		)
	}
	for k := range tokens {
		cols = append(cols, k)
	}
	sort.Strings(cols)
	return cols
}

func (t *Task) parseURL(tokens map[string]string) {
	rawurl, ok := tokens["request_url"]
	if ok {
		url, err := url.Parse(rawurl)
		if err != nil {
			t.debug.parse.Write(rawurl, err.Error())
		} else {
			pattern, normalized, err := t.patterns.Match(url)
			if err != nil {
				t.debug.normalize.Write(rawurl, err.Error())
			} else {
				if normalized == "" {
					t.debug.unrecognized.Write(url.Path, rawurl)
				}
				tokens["normalized_url"] = normalized
				tokens["url_pattern"] = pattern
			}
		}
	}
}

func (t *Task) parseUserAgent(tokens map[string]string) {
	uagent, ok := tokens["user_agent"]
	if ok && t.uaparser != nil {
		client := t.uaparser.Parse(uagent)
		tokens["client_device_family"] = client.Device.Family
		tokens["client_os_family"] = client.Os.Family
		tokens["client_os_major"] = client.Os.Major
		tokens["client_os_minor"] = client.Os.Minor
		tokens["client_os_patch"] = client.Os.Patch
		tokens["client_ua_family"] = client.UserAgent.Family
		tokens["client_ua_major"] = client.UserAgent.Major
		tokens["client_ua_minor"] = client.UserAgent.Minor
		tokens["client_ua_patch"] = client.UserAgent.Patch
	}
}

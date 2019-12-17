package transformer

import (
	"bufio"
	"compress/gzip"
	"encoding/csv"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/sjansen/carpenter/internal/patterns"
	"github.com/sjansen/carpenter/internal/tokenizer"
	"github.com/sjansen/carpenter/internal/uaparser"
)

type Task struct {
	Src    string
	Dst    string
	ErrDir string
	Suffix string
}

type Transformer struct {
	Patterns  patterns.Patterns
	Tokenizer *tokenizer.Tokenizer
	UAParser  *uaparser.Parser
}

func (t *Transformer) Transform(task *Task) error {
	f, err := os.Open(task.Src)
	if err != nil {
		return err
	}
	defer f.Close()

	var r io.ReadCloser = f
	if strings.HasSuffix(task.Src, ".gz") {
		r, err = gzip.NewReader(r)
		if err != nil {
			return err
		}
		defer r.Close()
	}

	err = os.MkdirAll(filepath.Dir(task.Dst), 0777)
	if err != nil {
		return err
	}

	w, err := os.OpenFile(task.Dst, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer w.Close()

	debug := newDebugFiles(task.ErrDir, task.Suffix)
	defer debug.Close()

	src := bufio.NewReader(r)
	dst := csv.NewWriter(w)
	return t.transform(src, dst, debug)
}

func (t *Transformer) newCols(tokens map[string]string) []string {
	cols := make([]string, 0, len(tokens)+11)
	cols = append(cols,
		"normalized_url",
		"url_pattern",
	)
	if t.UAParser != nil {
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

func (t *Transformer) parseURL(tokens map[string]string, debug *debugFiles) {
	rawurl, ok := tokens["request_url"]
	if ok {
		url, err := url.Parse(rawurl)
		if err != nil {
			if debug != nil {
				debug.parse.Write(rawurl, err.Error())
			}
		} else {
			pattern, normalized, err := t.Patterns.Match(url)
			if err != nil {
				if debug != nil {
					debug.parse.Write(rawurl, err.Error())
				}
			} else {
				if debug != nil && normalized == "" {
					debug.unrecognized.Write(url.Path, rawurl)
				}
				tokens["normalized_url"] = normalized
				tokens["url_pattern"] = pattern
			}
		}
	}
}

func (t *Transformer) parseUserAgent(tokens map[string]string) {
	uagent, ok := tokens["user_agent"]
	if ok && t.UAParser != nil {
		client := t.UAParser.Parse(uagent)
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

func (t *Transformer) transform(src *bufio.Reader, dst *csv.Writer, debug *debugFiles) error {
	defer dst.Flush()

	var cols []string
	var row []string
	for {
		line, err := src.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		tokens := t.Tokenizer.Tokenize(line)
		if tokens == nil {
			if debug != nil {
				debug.tokenize.Write(line, "\n")
			}
			continue
		} else if cols == nil {
			cols = t.newCols(tokens)
			dst.Write(cols)
			row = make([]string, len(cols))
		}

		t.parseURL(tokens, debug)

		t.parseUserAgent(tokens)

		for i, k := range cols {
			if v, ok := tokens[k]; ok {
				row[i] = v
			} else {
				row[i] = ""
			}
		}
		dst.Write(row)
	}

	return dst.Error()
}

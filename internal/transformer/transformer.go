package transformer

import (
	"bufio"
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"sort"
	"strings"

	"github.com/ua-parser/uap-go/uaparser"

	"github.com/sjansen/carpenter/internal/data"
	"github.com/sjansen/carpenter/internal/patterns"
	"github.com/sjansen/carpenter/internal/tokenizer"
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

	uaparser *uaparser.Parser
}

func (t *Transformer) EnableUserAgentParsing() error {
	r, err := data.Assets.Open("regexes.yaml")
	if err != nil {
		return err
	}

	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	uap, err := uaparser.NewFromBytes(bytes)
	if err != nil {
		return err
	}

	t.uaparser = uap
	return nil
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

	w, err := os.OpenFile(task.Dst, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer w.Close()

	src := bufio.NewReader(r)
	dst := csv.NewWriter(w)
	return t.transform(src, dst)
}

func (t *Transformer) cols(tokens map[string]string) []string {
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

func (t *Transformer) parseUserAgent(tokens map[string]string) {
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

func (t *Transformer) transform(src *bufio.Reader, dst *csv.Writer) error {
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
			continue
		} else if cols == nil {
			cols = t.cols(tokens)
			dst.Write(cols)
			row = make([]string, len(cols))
		}

		rawurl, ok := tokens["request_url"]
		if ok {
			url, err := url.Parse(rawurl)
			if err != nil {
				return fmt.Errorf(`unable to parse url: %q (%s)`, rawurl, err.Error())
			}

			pattern, normalized, err := t.Patterns.Match(url)
			if err != nil {
				return err
			}
			tokens["normalized_url"] = normalized
			tokens["url_pattern"] = pattern
		}

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

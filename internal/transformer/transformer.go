package transformer

import (
	"bufio"
	"compress/gzip"
	"encoding/csv"
	"io"
	"os"
	"sort"
	"strings"

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

	w, err := os.OpenFile(task.Dst, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer w.Close()

	src := bufio.NewReader(r)
	dst := csv.NewWriter(w)
	return t.transform(src, dst)
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
			cols = make([]string, 0, len(tokens)+2)
			row = make([]string, cap(cols))
			for k := range tokens {
				cols = append(cols, k)
			}
			cols = append(cols, "url_pattern")
			cols = append(cols, "url_normalized")
			sort.Strings(cols)
			dst.Write(cols)
		}

		rawurl, ok := tokens["request_url"]
		if ok {
			pattern, normalized, err := t.Patterns.Match(rawurl)
			if err != nil {
				return err
			}
			tokens["url_pattern"] = pattern
			tokens["url_normalized"] = normalized
		}

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

package worker

import (
	"bufio"
	"compress/gzip"
	"encoding/csv"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/sjansen/carpenter/internal/parser"
)

type Task struct {
	Src    string
	Dst    string
	ErrDir string
	Suffix string
}

type Transformer struct {
	Parser *parser.Parser
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

		parsed := t.Parser.Parse(line)
		if parsed == nil {
			continue
		} else if cols == nil {
			cols = make([]string, 0, len(parsed))
			row = make([]string, len(parsed))
			for k := range parsed {
				cols = append(cols, k)
			}
			sort.Strings(cols)
			dst.Write(cols)
		}

		for i, k := range cols {
			if v, ok := parsed[k]; ok {
				row[i] = v
			} else {
				row[i] = ""
			}
		}
		dst.Write(row)
	}

	return dst.Error()
}

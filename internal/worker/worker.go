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

func Transform(t *Task) error {
	f, err := os.Open(t.Src)
	if err != nil {
		return err
	}
	defer f.Close()

	var r io.ReadCloser = f
	if strings.HasSuffix(t.Src, ".gz") {
		r, err = gzip.NewReader(r)
		if err != nil {
			return err
		}
		defer r.Close()
	}

	w, err := os.OpenFile(t.Dst, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer w.Close()

	src := bufio.NewReader(r)
	dst := csv.NewWriter(w)
	defer dst.Flush()

	var cols []string
	var vals []string
	for {
		var line string
		line, err = src.ReadString('\n')
		if err != nil {
			break
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parsed := parser.ALB.Parse(line)
		if parsed == nil {
			continue
		} else if cols == nil {
			cols = make([]string, 0, len(parsed))
			vals = make([]string, len(parsed))
			for k := range parsed {
				cols = append(cols, k)
			}
			sort.Strings(cols)
			dst.Write(cols)
		}

		for i, k := range cols {
			if v, ok := parsed[k]; ok {
				vals[i] = v
			} else {
				vals[i] = ""
			}
		}
		dst.Write(vals)
	}
	if err != io.EOF {
		return err
	}

	return dst.Error()
}

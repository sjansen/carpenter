package pipeline

import (
	"bufio"
	"io"
	pathlib "path"
	"runtime"

	"github.com/sjansen/carpenter/internal/lazyio"
	"github.com/sjansen/carpenter/internal/patterns"
	"github.com/sjansen/carpenter/internal/tokenizer"
	"github.com/sjansen/carpenter/internal/uaparser"
)

type Pipeline struct {
	Debug     lazyio.Opener
	Result    lazyio.Opener
	Patterns  *patterns.Patterns
	Tokenizer *tokenizer.Tokenizer
	UAParser  *uaparser.Parser
}

func (p *Pipeline) NewTask(r io.Reader, path string) *Task {
	base := path
	if n := len(pathlib.Ext(path)); n > 0 {
		base = path[:len(path)-n]
	}

	task := &Task{
		patterns:  p.Patterns,
		tokenizer: p.Tokenizer,
		uaparser:  p.UAParser,
		src:       bufio.NewReader(r),
		dst: lazyio.CSV{
			Opener: p.Result,
			Path:   base + ".csv",
		},
	}

	if p.Debug != nil {
		task.debug = debug{
			normalize: lazyio.CSV{
				Opener: p.Debug,
				Path:   pathlib.Join("normalize", base+".csv"),
			},
			parse: lazyio.CSV{
				Opener: p.Debug,
				Path:   pathlib.Join("parse", base+".csv"),
			},
			tokenize: lazyio.TXT{
				Opener: p.Debug,
				Path:   pathlib.Join("tokenize", base+".txt"),
			},
			unrecognized: lazyio.CSV{
				Opener: p.Debug,
				Path:   pathlib.Join("unrecognized", base+".csv"),
			},
		}
	}

	return task
}

func (p *Pipeline) Start() chan<- *Task {
	ch := make(chan *Task)
	for i := runtime.NumCPU(); i > 0; i-- {
		go worker(ch)
	}
	return ch
}

func worker(ch <-chan *Task) {
	for t := range ch {
		t.Run()
	}
}

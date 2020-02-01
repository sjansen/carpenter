package pipeline

import (
	pathlib "path"
	"runtime"
	"sync"

	"github.com/sjansen/carpenter/internal/lazyio"
	"github.com/sjansen/carpenter/internal/patterns"
	"github.com/sjansen/carpenter/internal/tokenizer"
	"github.com/sjansen/carpenter/internal/uaparser"
)

type Pipeline struct {
	Patterns  *patterns.Patterns
	Tokenizer *tokenizer.Tokenizer
	UAParser  *uaparser.Parser

	Source lazyio.InputOpener
	Result lazyio.OutputOpener
	Debug  lazyio.OutputOpener

	ch chan<- *Task
	wg sync.WaitGroup
}

func (p *Pipeline) AddTask(path string) {
	input := &lazyio.Input{
		Path:   path,
		Opener: p.Source,
	}
	base := input.StripExt()

	task := &Task{
		patterns:  p.Patterns,
		tokenizer: p.Tokenizer,
		uaparser:  p.UAParser,
		src:       input,
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

	p.ch <- task
}

func (p *Pipeline) Start() {
	ch := make(chan *Task)
	for i := runtime.NumCPU(); i > 0; i-- {
		p.wg.Add(1)
		go worker(ch, &p.wg)
	}
	p.ch = ch
}

func (p *Pipeline) Wait() {
	close(p.ch)
	p.wg.Wait()
}

func worker(ch <-chan *Task, wg *sync.WaitGroup) {
	for t := range ch {
		t.Run()
	}
	wg.Done()
}

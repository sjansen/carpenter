package pipeline

import (
	pathlib "path"
	"runtime"
	"strconv"
	"sync"

	"github.com/sjansen/carpenter/internal/lazyio"
	"github.com/sjansen/carpenter/internal/patterns"
	"github.com/sjansen/carpenter/internal/sys"
	"github.com/sjansen/carpenter/internal/tokenizer"
	"github.com/sjansen/carpenter/internal/uaparser"
)

type Pipeline struct {
	Patterns  *patterns.Patterns
	Tokenizer *tokenizer.Tokenizer
	UAParser  *uaparser.Parser

	IO     *sys.IO
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
		go worker(p.IO, i, ch, &p.wg)
	}
	p.ch = ch
}

func (p *Pipeline) Wait() {
	close(p.ch)
	p.wg.Wait()
}

func worker(io *sys.IO, id int, ch <-chan *Task, wg *sync.WaitGroup) {
	log := io.Log
	log.Debugw("starting worker", "id", id)
	prefix := strconv.Itoa(id) + ":"
	var i uint64
	for t := range ch {
		t.id = prefix + strconv.FormatUint(i, 10)
		log.Debugw("processing task", "worker", id, "task", t.id, "path", t.src.Path)
		if err := t.Run(); err != nil {
			log.Debugw("task error returned", "task", t.id, "path", t.src.Path)
		}
		log.Debugw("completed task", "worker", id, "task", t.id, "path", t.src.Path)
		i++
	}
	log.Debugw("stopping worker", "id", id)
	wg.Done()
}

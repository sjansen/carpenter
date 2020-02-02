package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sjansen/carpenter/internal/lazyio"
	"github.com/sjansen/carpenter/internal/patterns"
	"github.com/sjansen/carpenter/internal/pipeline"
	"github.com/sjansen/carpenter/internal/tokenizer"
	"github.com/sjansen/carpenter/internal/uaparser"
)

type TransformCmd struct {
	Patterns string
	SrcURI   string
	DstURI   string
	ErrURI   string
}

func (c *TransformCmd) Run(base *Base) error {
	pipeline, walker, err := c.newPipeline()
	if err != nil {
		return err
	}

	pipeline.Start()
	defer pipeline.Wait()

	return walker.Walk(func(path string) error {
		pipeline.AddTask(path)
		return nil
	})
}

func (c *TransformCmd) newPipeline() (*pipeline.Pipeline, lazyio.InputWalker, error) {
	r, err := os.Open(c.Patterns)
	if err != nil {
		return nil, nil, err
	}

	patterns, err := patterns.Load(c.Patterns, r)
	if err != nil {
		return nil, nil, err
	}

	uaparser, err := uaparser.UserAgentParser()
	if err != nil {
		return nil, nil, err
	}

	pipeline := &pipeline.Pipeline{
		Patterns:  patterns,
		Tokenizer: tokenizer.ALB,
		UAParser:  uaparser,
	}

	input, err := newInputOpenWalker(c.SrcURI)
	if err != nil {
		return nil, nil, err
	}
	pipeline.Source = input

	output, err := newOutputOpener(c.DstURI)
	if err != nil {
		return nil, nil, err
	}
	pipeline.Result = output

	if c.ErrURI != "" {
		opener, err := newOutputOpener(c.ErrURI)
		if err != nil {
			return nil, nil, err
		}
		pipeline.Debug = opener
	}

	return pipeline, input, nil
}

type inputOpenWalker interface {
	lazyio.InputOpener
	lazyio.InputWalker
}

func newInputOpenWalker(uri string) (inputOpenWalker, error) {
	switch {
	case strings.HasPrefix(uri, "s3://") || strings.HasPrefix(uri, "S3://"):
		return lazyio.NewS3Reader(uri)
	default:
		uri = filepath.Clean(uri)
		if info, err := os.Stat(uri); err != nil {
			return nil, err
		} else if !info.IsDir() {
			return nil, fmt.Errorf("error: not a directory %q", uri)
		}
		return &lazyio.FileReader{Dir: uri}, nil
	}
}

func newOutputOpener(uri string) (lazyio.OutputOpener, error) {
	switch {
	case strings.HasPrefix(uri, "s3://") || strings.HasPrefix(uri, "S3://"):
		return lazyio.NewS3Writer(uri)
	default:
		uri = filepath.Clean(uri)
		if info, err := os.Stat(uri); err != nil {
			if !os.IsNotExist(err) {
				return nil, err
			}
			if err = os.MkdirAll(uri, 0777); err != nil {
				return nil, err
			}
		} else if !info.IsDir() {
			return nil, fmt.Errorf("error: not a directory %q", uri)
		}
		return &lazyio.FileWriter{Dir: uri}, nil
	}
}

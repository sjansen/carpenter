package cmd

import (
	"compress/gzip"
	"fmt"
	"io"
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
	if err := c.verifyArgs(); err != nil {
		return err
	}

	pipeline, err := c.newPipeline()
	if err != nil {
		return err
	}

	return filepath.Walk(c.SrcURI, func(src string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		} else if !info.Mode().IsRegular() {
			return nil
		}

		f, err := os.Open(src)
		if err != nil {
			return err
		}
		defer f.Close()

		var r io.ReadCloser = f
		suffix := src[len(c.SrcURI):]
		if strings.HasSuffix(src, ".gz") {
			suffix = suffix[:len(suffix)-3]
			r, err = gzip.NewReader(r)
			if err != nil {
				return err
			}
			defer r.Close()
		}

		task := pipeline.NewTask(r, suffix)
		return task.Run()
	})
}

func (c *TransformCmd) newPipeline() (*pipeline.Pipeline, error) {
	r, err := os.Open(c.Patterns)
	if err != nil {
		return nil, err
	}

	patterns, err := patterns.Load(c.Patterns, r)
	if err != nil {
		return nil, err
	}

	uaparser, err := uaparser.UserAgentParser()
	if err != nil {
		return nil, err
	}

	pipeline := &pipeline.Pipeline{
		Result:    &lazyio.FileOpener{Dir: c.DstURI},
		Patterns:  patterns,
		Tokenizer: tokenizer.ALB,
		UAParser:  uaparser,
	}
	if c.ErrURI != "" {
		pipeline.Debug = &lazyio.FileOpener{Dir: c.ErrURI}
	}

	return pipeline, nil
}

func (c *TransformCmd) verifyArgs() error {
	c.SrcURI = filepath.Clean(c.SrcURI)
	if info, err := os.Stat(c.SrcURI); err != nil {
		return err
	} else if !info.IsDir() {
		return fmt.Errorf("error: not a directory %q", c.SrcURI)
	}

	c.DstURI = filepath.Clean(c.DstURI)
	if info, err := os.Stat(c.DstURI); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		if err = os.MkdirAll(c.DstURI, 0777); err != nil {
			return err
		}
	} else if !info.IsDir() {
		return fmt.Errorf("error: not a directory %q", c.DstURI)
	}

	if c.ErrURI != "" {
		c.ErrURI = filepath.Clean(c.ErrURI)
		if info, err := os.Stat(c.ErrURI); err != nil {
			if !os.IsNotExist(err) {
				return err
			}
			if err = os.MkdirAll(c.ErrURI, 0777); err != nil {
				return err
			}
		} else if !info.IsDir() {
			return fmt.Errorf("error: not a directory %q", c.ErrURI)
		}
	}

	return nil
}

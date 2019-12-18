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
	SrcDir   string
	DstDir   string
	ErrDir   string
}

func (c *TransformCmd) Run(base *Base) error {
	if err := c.verifyArgs(); err != nil {
		return err
	}

	pipeline, err := c.newPipeline()
	if err != nil {
		return err
	}

	return filepath.Walk(c.SrcDir, func(src string, info os.FileInfo, err error) error {
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
		suffix := src[len(c.SrcDir):]
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
		Result:    &lazyio.File{Dir: c.DstDir},
		Patterns:  patterns,
		Tokenizer: tokenizer.ALB,
		UAParser:  uaparser,
	}
	if c.ErrDir != "" {
		pipeline.Debug = &lazyio.File{Dir: c.ErrDir}
	}

	return pipeline, nil
}

func (c *TransformCmd) verifyArgs() error {
	c.SrcDir = filepath.Clean(c.SrcDir)
	if info, err := os.Stat(c.SrcDir); err != nil {
		return err
	} else if !info.IsDir() {
		return fmt.Errorf("error: not a directory %q", c.SrcDir)
	}

	c.DstDir = filepath.Clean(c.DstDir)
	if info, err := os.Stat(c.DstDir); err != nil {
		if os.IsNotExist(err) {
			return os.MkdirAll(c.DstDir, 0777)
		}
		return err
	} else if !info.IsDir() {
		return fmt.Errorf("error: not a directory %q", c.DstDir)
	}

	if c.ErrDir != "" {
		c.ErrDir = filepath.Clean(c.ErrDir)
		if info, err := os.Stat(c.ErrDir); err != nil {
			if os.IsNotExist(err) {
				return os.MkdirAll(c.ErrDir, 0777)
			}
			return err
		} else if !info.IsDir() {
			return fmt.Errorf("error: not a directory %q", c.ErrDir)
		}
	}

	return nil
}

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sjansen/carpenter/internal/patterns"
	"github.com/sjansen/carpenter/internal/tokenizer"
	"github.com/sjansen/carpenter/internal/transformer"
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

	r, err := os.Open(c.Patterns)
	if err != nil {
		return err
	}
	patterns, err := patterns.Load(c.Patterns, r)
	if err != nil {
		return err
	}

	uaparser, err := uaparser.UserAgentParser()
	if err != nil {
		return err
	}

	transformer := &transformer.Transformer{
		Patterns:  patterns,
		Tokenizer: tokenizer.ALB,
		UAParser:  uaparser,
	}

	return filepath.Walk(c.SrcDir, func(src string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		} else if !info.Mode().IsRegular() {
			return nil
		}

		task := c.newTask(src)
		return transformer.Transform(task)
	})
}

func (c *TransformCmd) newTask(src string) *transformer.Task {
	suffix := src[len(c.SrcDir):]
	if strings.HasSuffix(suffix, ".gz") {
		suffix = suffix[:len(suffix)-3]
	}
	ext := filepath.Ext(suffix)
	if ext != "" {
		suffix = suffix[:len(suffix)-len(ext)]
	}

	dst := filepath.Join(c.DstDir, suffix+".csv")
	return &transformer.Task{
		Src:    src,
		Dst:    dst,
		ErrDir: c.ErrDir,
		Suffix: suffix,
	}
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

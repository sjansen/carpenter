package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sjansen/carpenter/internal/parser"
	"github.com/sjansen/carpenter/internal/worker"
)

type TransformCmd struct {
	SrcDir string
	DstDir string
	ErrDir string
}

func (c *TransformCmd) Run(base *Base) error {
	if err := c.verifyArgs(); err != nil {
		return err
	}

	if err := parser.ALB.EnableUserAgentParsing(); err != nil {
		return err
	}

	transformer := &worker.Transformer{
		Parser: parser.ALB,
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

func (c *TransformCmd) newTask(src string) *worker.Task {
	suffix := src[len(c.SrcDir):]
	if strings.HasSuffix(suffix, ".gz") {
		suffix = suffix[:len(suffix)-3]
	}
	ext := filepath.Ext(suffix)
	if ext != "" {
		suffix = suffix[:len(suffix)-len(ext)]
	}

	dst := filepath.Join(c.DstDir, suffix+".csv")
	return &worker.Task{
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

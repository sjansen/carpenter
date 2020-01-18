package cmd

import (
	"os"

	"github.com/sjansen/carpenter/internal/patterns"
)

type TestCmd struct {
	File    string
	Verbose bool
}

func (c *TestCmd) Run(base *Base) error {
	r, err := os.Open(c.File)
	if err != nil {
		return err
	}

	patterns, err := patterns.Load(c.File, r)
	if err != nil {
		return err
	}

	_, err = patterns.Test(&base.IO)
	return err
}

package cmd

import (
	"os"

	"github.com/sjansen/carpenter/internal/oldpatterns"
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

	patterns, err := oldpatterns.Load(c.File, r)
	if err != nil {
		return err
	}

	_, err = patterns.SelfTest(&base.IO)
	return err
}

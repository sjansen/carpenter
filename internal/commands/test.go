package commands

import (
	"os"

	"github.com/sjansen/carpenter/internal/rules"
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

	rules, err := rules.Load(c.File, r)
	if err != nil {
		return err
	}

	_, err = rules.SelfTest(&base.IO)
	return err
}

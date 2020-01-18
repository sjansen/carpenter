package cmd

import (
	"encoding/json"
	"os"

	"github.com/sjansen/carpenter/internal/patterns"
)

type TestCasesCmd struct {
	File    string
	Verbose bool
}

func (c *TestCasesCmd) Run(base *Base) error {
	r, err := os.Open(c.File)
	if err != nil {
		return err
	}

	patterns, err := patterns.Load(c.File, r)
	if err != nil {
		return err
	}

	testcases, err := patterns.Test(&base.IO)
	if err != nil {
		return err
	}

	enc := json.NewEncoder(base.IO.Stdout)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	return enc.Encode(testcases)
}

package cmd

import (
	"encoding/json"
	"os"

	"github.com/sjansen/carpenter/internal/oldpatterns"
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

	patterns, err := oldpatterns.Load(c.File, r)
	if err != nil {
		return err
	}

	testcases, err := patterns.SelfTest(&base.IO)
	if err != nil {
		return err
	}

	enc := json.NewEncoder(base.IO.Stdout)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	return enc.Encode(testcases)
}

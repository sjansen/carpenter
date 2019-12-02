package cli

import "github.com/sjansen/carpenter/internal/cmd"

func registerTestCases(p *ArgParser) {
	c := &cmd.TestCasesCmd{}
	cmd := p.addCommand(c, "test-cases", "TODO")
	cmd.Arg("FILE", "A pattern file").Required().
		ExistingFileVar(&c.File)
}

package cli

import "github.com/sjansen/carpenter/internal/commands"

func registerTestCases(p *ArgParser) {
	c := &commands.TestCasesCmd{}
	cmd := p.addCommand(c, "test-cases", "TODO")
	cmd.Arg("FILE", "A rules file").Required().
		ExistingFileVar(&c.File)
}

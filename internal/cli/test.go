package cli

import "github.com/sjansen/carpenter/internal/commands"

func registerTest(p *ArgParser) {
	c := &commands.TestCmd{}
	cmd := p.addCommand(c, "test", "Evaluate rules file self tests")
	cmd.Arg("FILE", "A rules file").Required().
		ExistingFileVar(&c.File)
}

package cli

import "github.com/sjansen/carpenter/internal/cmd"

func registerTest(p *ArgParser) {
	c := &cmd.TestCmd{}
	cmd := p.addCommand(c, "test", "Evaluate pattern file self tests")
	cmd.Arg("FILE", "A pattern file").Required().
		ExistingFileVar(&c.File)
}

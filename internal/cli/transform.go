package cli

import "github.com/sjansen/carpenter/internal/cmd"

func registerTransform(p *ArgParser) {
	c := &cmd.TransformCmd{}
	cmd := p.addCommand(c, "transform", "TODO")
	cmd.Arg("PATTERNS", "Pattern file").Required().
		ExistingFileVar(&c.Patterns)
	cmd.Arg("SRC", "Source directory").Required().
		ExistingDirVar(&c.SrcURI)
	cmd.Arg("DST", "Target directory").Required().
		StringVar(&c.DstURI)
	cmd.Arg("ERRORS", "Errors directory").
		StringVar(&c.ErrURI)
}

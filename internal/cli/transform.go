package cli

import "github.com/sjansen/carpenter/internal/cmd"

func registerTransform(p *ArgParser) {
	c := &cmd.TransformCmd{}
	cmd := p.addCommand(c, "transform", "TODO")
	cmd.Arg("PATTERNS", "Pattern file").Required().
		ExistingFileVar(&c.Patterns)
	cmd.Arg("SRC", "Source directory").Required().
		ExistingDirVar(&c.SrcDir)
	cmd.Arg("DST", "Target directory").Required().
		StringVar(&c.DstDir)
	cmd.Arg("ERRORS", "Errors directory").
		StringVar(&c.ErrDir)
}

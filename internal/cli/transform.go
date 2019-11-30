package cli

import "github.com/sjansen/carpenter/internal/commands"

func registerTransform(p *ArgParser) {
	c := &commands.TransformCmd{}
	cmd := p.addCommand(c, "transform", "TODO")
	cmd.Arg("SRC", "Source directory").Required().
		ExistingDirVar(&c.SrcDir)
	cmd.Arg("DST", "Target directory").Required().
		StringVar(&c.DstDir)
	cmd.Arg("ERRORS", "Errors directory").
		StringVar(&c.ErrDir)
}

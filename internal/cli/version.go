package cli

import (
	"runtime/debug"

	"github.com/sjansen/carpenter/internal/cmd"
)

func registerVersion(p *ArgParser, build string) {
	info, _ := debug.ReadBuildInfo()
	c := &cmd.VersionCmd{
		App:   "carpenter",
		Build: build,

		BuildInfo: info,
		Verbose:   false,
	}
	cmd := p.addCommand(c, "version", "Print carpenter's version")
	cmd.Flag("long", "include build details").Short('l').BoolVar(&c.Verbose)
}

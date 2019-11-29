package cli

import (
	"io"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/sjansen/carpenter/internal/commands"
	"github.com/sjansen/carpenter/internal/logger"
)

// An ArgParser contains the definitions of flags, arguments and commands for an application
type ArgParser struct {
	app     *kingpin.Application
	base    *commands.Base
	cmd     command
	version string

	brevity   int
	verbosity int
}

// Command is an abstraction for a function that produces human readable output when run
type Command func(stdout, stderr io.Writer) error

type command interface {
	Run(base *commands.Base) error
}

// Parse converts command line arguments into a Command
func (p *ArgParser) Parse(args []string) (Command, error) {
	_, err := p.app.Parse(args)
	if err != nil {
		return nil, err
	}
	fn := func(stdout, stderr io.Writer) error {
		p.base.Verbosity = 1 + p.verbosity - p.brevity
		if p.base.Debug == nil {
			p.base.Log = logger.New(p.base.Verbosity, stdout, nil)
		} else {
			p.base.Log = logger.New(p.base.Verbosity, stdout, p.base.Debug)
		}
		p.base.Stdout = stdout
		p.base.Stderr = stderr
		p.base.Log.Infof("carpenter version=%s", p.version)
		return p.cmd.Run(p.base)
	}
	return fn, nil
}

func (p *ArgParser) addCommand(cmd command, name, help string) *kingpin.CmdClause {
	clause := p.app.Command(name, help)
	clause.Action(func(pc *kingpin.ParseContext) error {
		p.cmd = cmd
		return nil
	})
	return clause
}

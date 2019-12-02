package cmd

import (
	"fmt"
	"runtime/debug"
)

// A VersionCmd is a request to reports the app's version
type VersionCmd struct {
	App   string
	Build string

	BuildInfo *debug.BuildInfo
	Verbose   bool
}

// Run executes the command
func (c *VersionCmd) Run(base *Base) error {
	fmt.Fprintln(base.Stdout, c.App, c.Build)
	if c.Verbose && c.BuildInfo != nil {
		fmt.Println("")
		for _, m := range c.BuildInfo.Deps {
			fmt.Printf("%-40s %-35s\n", m.Path, m.Version)
		}
	}
	return nil
}

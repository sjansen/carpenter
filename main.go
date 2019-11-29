package main

import (
	"fmt"
	"os"

	"github.com/sjansen/carpenter/internal/cli"
)

var build string // set by goreleaser

func main() {
	if build == "" {
		build = version
	}
	parser := cli.RegisterCommands(build)

	cmd, err := parser.Parse(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = cmd(os.Stdout, os.Stderr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

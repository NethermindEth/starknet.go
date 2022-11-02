package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

var versionCommand = cli.Command{
	Name:    "version",
	Aliases: []string{"v"},
	Usage:   "get the version",
	Action:  versionAction,
}

// add -ldflags="-X 'main.version=demo'"
var version = "dev"
var buildTime = "unknown"

func versionAction(cCtx *cli.Context) error {
	fmt.Printf("version %s, built %s\n", version, buildTime)
	return nil
}

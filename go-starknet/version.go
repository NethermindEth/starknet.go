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

// versionAction prints the version and build time of the application.
//
// cCtx - the cli.Context object used to access command line arguments and flags.
// error - an error object if there was an issue printing the version.
func versionAction(cCtx *cli.Context) error {
	fmt.Printf("version %s, built %s\n", version, buildTime)
	return nil
}

package main

import (
	_ "embed"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

//go:embed dictionary.json
var dictionary []byte

// main is the entry point of go-starknet the program.
//
// It initializes or loads the configuration and sets up the command-line interface.
// 
// Parameters:
// None.
//
// Return:
// None.
func main() {
	var err error
	configuration, err = initOrLoadConfig()
	if err != nil {
		return
	}

	app := &cli.App{

		Commands: []*cli.Command{
			&blockCommand,
			&transactionCommand,
			&utilsCommand,
			&settingsCommand,
			&versionCommand,
			&workspaceCommand,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

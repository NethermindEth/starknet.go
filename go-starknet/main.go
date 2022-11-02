package main

import (
	_ "embed"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

//go:embed dictionary.json
var dictionary []byte

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

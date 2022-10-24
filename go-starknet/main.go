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
	app := &cli.App{
		Commands: []*cli.Command{
			&blockCommand,
			&transactionCommand,
			&utilsCommand,
			&profileCommand,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

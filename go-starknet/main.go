package main

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/dontpanicdao/caigo"
	"github.com/dontpanicdao/caigo/types"
	"github.com/urfave/cli/v2"
)

//go:embed dictionary.txt
var dictionary []byte

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "utils",
				Aliases: []string{"u"},
				Usage:   "utilities to encode/decode felt and entrypoints",
				Subcommands: []*cli.Command{
					{
						Name:  "felt",
						Usage: "display felts in different formats",
						Action: func(cCtx *cli.Context) error {
							for _, v := range cCtx.Args().Slice() {
								fmt.Printf("value:   %s\n", v)
								vInt, ok := big.NewInt(0).SetString(v, 0)
								if !ok {
									fmt.Println("could not guess format")
								}
								fmt.Printf("decimal: %s\n", vInt.Text(10))
								fmt.Printf("hex:     0x%s\n", vInt.Text(16))
								fmt.Println()
							}
							return nil
						},
					},
					{
						Name:  "entrypoint",
						Usage: "display entrypoints in different formats",
						Action: func(cCtx *cli.Context) error {
							for _, v := range cCtx.Args().Slice() {
								fmt.Printf("value:   %s\n", v)
								fmt.Printf("decimal: %s\n", types.GetSelectorFromName(v).Text(10))
								fmt.Printf("hex:     0x%s\n", types.GetSelectorFromName(v).Text(16))
								fmt.Println()
							}
							return nil
						},
					},
					{
						Name:  "guess",
						Usage: "guess an entrypoint from a dictionary",
						Action: func(cCtx *cli.Context) error {
							dict := map[string]string{}
							err := json.Unmarshal(dictionary, &dict)
							if err != nil {
								return err
							}
							v := cCtx.Args().First()
							vInt, ok := big.NewInt(0).SetString(v, 0)
							if !ok {
								return errors.New("not a number")
							}
							value, ok := dict[fmt.Sprintf("0x%s", vInt.Text(16))]
							if !ok {
								fmt.Println("key not in dictionary")
								return nil
							}
							fmt.Printf("guess:   %s\n", value)
							fmt.Printf("decimal: %s\n", types.GetSelectorFromName(value).Text(10))
							fmt.Printf("hex:     0x%s\n", types.GetSelectorFromName(value).Text(16))
							fmt.Println()
							return nil
						},
					},
					{
						Name:  "publickey",
						Usage: "get a public key from the private key",
						Action: func(cCtx *cli.Context) error {

							privateKey := cCtx.Args().First()
							privateKeyInt, ok := big.NewInt(0).SetString(privateKey, 0)
							if !ok {
								return errors.New("not a number")
							}
							publicKey, _, err := caigo.Curve.PrivateToPoint(privateKeyInt)
							if err != nil {
								return err
							}
							fmt.Printf("public key: 0x%s\n", publicKey.Text(16))
							fmt.Println()
							return nil
						},
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

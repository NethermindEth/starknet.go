package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/dontpanicdao/caigo"
	"github.com/dontpanicdao/caigo/types"
	"github.com/urfave/cli/v2"
)

var signFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    "private-key",
		Aliases: []string{"p"},
		Usage:   "private key used to sign a message",
	}}

var verifyFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    "public-key",
		Aliases: []string{"p"},
		Usage:   "public key used to sign a message",
	},
	&cli.StringFlag{
		Name:  "s0",
		Usage: "first signature",
	},
	&cli.StringFlag{
		Name:  "s1",
		Usage: "2nd signature",
	},
}

var utilsCommand = cli.Command{
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
				fmt.Printf("public key: %s\n", publicKey.Text(10))
				fmt.Println()
				return nil
			},
		},
		{
			Name:  "sign",
			Usage: "sign a message with a private key",
			Flags: signFlags,
			Action: func(cCtx *cli.Context) error {
				if !cCtx.IsSet("private-key") {
					return errors.New("private key is mandatory")
				}
				privateKey := cCtx.String("private-key")
				privateKeyInt, ok := big.NewInt(0).SetString(privateKey, 0)
				if !ok {
					return errors.New("not a number")
				}
				message := cCtx.Args().First()
				messageInt, ok := big.NewInt(0).SetString(message, 0)
				if !ok {
					return errors.New("not a number")
				}
				publicKey, _, err := caigo.Curve.PrivateToPoint(privateKeyInt)
				if err != nil {
					return err
				}
				x, y, err := caigo.Curve.Sign(messageInt, privateKeyInt)
				if err != nil {
					return err
				}
				fmt.Printf("public key:   0x%s\n", publicKey.Text(16))
				fmt.Printf("public key:   %s\n", publicKey.Text(10))
				fmt.Printf("signature[0]: 0x%s\n", x.Text(16))
				fmt.Printf("signature[0]: %s\n", x.Text(10))
				fmt.Printf("signature[1]: 0x%s\n", y.Text(16))
				fmt.Printf("signature[1]: %s\n", y.Text(10))
				fmt.Println()
				return nil
			},
		},
		{
			Name:  "verify",
			Usage: "verify a signature with a public key and the message",
			Flags: verifyFlags,
			Action: func(cCtx *cli.Context) error {
				if !cCtx.IsSet("public-key") {
					return errors.New("public key is mandatory")
				}
				if !cCtx.IsSet("s0") {
					return errors.New("s0 is mandatory")
				}
				if !cCtx.IsSet("s1") {
					return errors.New("s1 is mandatory")
				}
				publicKey := cCtx.String("public-key")
				publicKeyInt, ok := big.NewInt(0).SetString(publicKey, 0)
				if !ok {
					return errors.New("not a number")
				}
				y := caigo.Curve.GetYCoordinate(publicKeyInt)
				s0 := cCtx.String("s0")
				s0Int, ok := big.NewInt(0).SetString(s0, 0)
				if !ok {
					return errors.New("not a number")
				}
				s1 := cCtx.String("s1")
				s1Int, ok := big.NewInt(0).SetString(s1, 0)
				if !ok {
					return errors.New("not a number")
				}
				message := cCtx.Args().First()
				messageInt, ok := big.NewInt(0).SetString(message, 0)
				if !ok {
					return errors.New("not a number")
				}
				ok = caigo.Curve.Verify(messageInt, s0Int, s1Int, publicKeyInt, y)
				if !ok {
					return errors.New("invalid signature")
				}
				fmt.Println("signature is valid")
				fmt.Println()
				return nil
			},
		},
	},
}

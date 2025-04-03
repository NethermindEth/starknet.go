package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/NethermindEth/starknet.go/abigen/accounts/abi/abigen"
	"github.com/urfave/cli/v2"
)

var (
	// Flags needed by abigen
	abiFlag = &cli.StringFlag{
		Name:  "abi",
		Usage: "Path to the Cairo contract ABI json to bind, - for STDIN",
	}
	binFlag = &cli.StringFlag{
		Name:  "bin",
		Usage: "Path to the Cairo contract bytecode (generate deploy method)",
	}
	typeFlag = &cli.StringFlag{
		Name:  "type",
		Usage: "Struct name for the binding (default = package name)",
	}
	pkgFlag = &cli.StringFlag{
		Name:  "pkg",
		Usage: "Package name to generate the binding into",
		Value: "main",
	}
	outFlag = &cli.StringFlag{
		Name:  "out",
		Usage: "Output file for the generated binding (default = stdout)",
	}
)

var app = &cli.App{
	Name:  "abigen",
	Usage: "Cairo contract binding generator",
}

func init() {
	app.Flags = []cli.Flag{
		abiFlag,
		binFlag,
		typeFlag,
		pkgFlag,
		outFlag,
	}
	app.Action = generate
}

func generate(c *cli.Context) error {
	if c.String(abiFlag.Name) == "" {
		return fmt.Errorf("no ABI specified (--abi)")
	}

	var (
		abiJSON string
		err     error
	)
	input := c.String(abiFlag.Name)
	if input == "-" {
		abiBytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("failed to read ABI from stdin: %v", err)
		}
		abiJSON = string(abiBytes)
	} else {
		abiBytes, err := os.ReadFile(input)
		if err != nil {
			return fmt.Errorf("failed to read ABI from file: %v", err)
		}
		abiJSON = string(abiBytes)
	}

	var contractData map[string]interface{}
	if err := json.Unmarshal([]byte(abiJSON), &contractData); err == nil {
		if abi, ok := contractData["abi"]; ok {
			abiBytes, _ := json.Marshal(abi)
			abiJSON = string(abiBytes)
		}
	}

	var binJSON string
	if c.String(binFlag.Name) != "" {
		binBytes, err := os.ReadFile(c.String(binFlag.Name))
		if err != nil {
			return fmt.Errorf("failed to read bytecode from file: %v", err)
		}
		binJSON = string(binBytes)
	}

	typeName := c.String(typeFlag.Name)
	if typeName == "" {
		typeName = c.String(pkgFlag.Name)
		if input != "-" {
			baseName := filepath.Base(input)
			typeName = strings.TrimSuffix(baseName, filepath.Ext(baseName))
		}
		typeName = abigen.ToCamelCase(typeName)
	}

	code, err := abigen.BindCairo([]string{typeName}, []string{abiJSON}, []string{binJSON}, c.String(pkgFlag.Name))
	if err != nil {
		return fmt.Errorf("failed to generate binding: %v", err)
	}

	if c.String(outFlag.Name) == "" {
		fmt.Println(code)
	} else {
		if err := os.WriteFile(c.String(outFlag.Name), []byte(code), 0644); err != nil {
			return fmt.Errorf("failed to write binding to file: %v", err)
		}
	}

	return nil
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

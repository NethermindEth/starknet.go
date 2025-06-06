package main

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var (
	// Used for flags.
	abiPath  string
	pkgName  string
	typeName string
	outPath  string

	// main CLI command
	rootCmd = &cobra.Command{}
)

func init() {
	// set flags
	rootCmd.PersistentFlags().StringVar(&abiPath, "abi", "", "path to the Cairo contract ABI")
	rootCmd.MarkPersistentFlagRequired("abi")
	// @todo at the end, validate default values correctness
	rootCmd.PersistentFlags().StringVar(&pkgName, "pkg", "", "package name for the generated Go code (default is the lower case contract name)")
	rootCmd.PersistentFlags().StringVar(&typeName, "type", "", "type name for the generated Go code (default is the upper case contract name)")
	rootCmd.PersistentFlags().StringVar(&outPath, "out", "", "output directory for the generated Go code (default is the current directory)")
}

func main() {
	rootCmd.Use = "starkgen --abi --pkg [--type] [--out]"
	rootCmd.Short = "Generate Go bindings for Cairo contracts"
	rootCmd.Long = `Generate Go bindings for Cairo contracts
								Creating now random content just to test
								how this text will be displayed in the termina.`

	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		ticker := time.NewTicker(1 * time.Second)
		go func() {
			for {
				time := <-ticker.C
				fmt.Printf("%s\n", time)
			}
		}()
		select {}
	}

	rootCmd.Execute()
}

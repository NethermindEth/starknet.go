package main

import (
	"fmt"
	"log"
	"os"

	"github.com/NethermindEth/starknet.go/contracts"
)

func main() {
	// This demonstrates how to load a CASM class from a file
	// For demo purposes, we'll show the API usage

	// In real usage:
	// casmClass, err := contracts.UnmarshalCasmClass("path/to/contract.casm.json")

	fmt.Println("UnmarshalCasmClass Demo:")
	fmt.Println("  Function: UnmarshalCasmClass(filePath string) (*CasmClass, error)")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  casmClass, err := contracts.UnmarshalCasmClass(\"contract.casm.json\")")
	fmt.Println("  if err != nil {")
	fmt.Println("      log.Fatal(err)")
	fmt.Println("  }")
	fmt.Println()
	fmt.Println("CasmClass contains:")
	fmt.Println("  - ByteCode: Compiled Cairo bytecode")
	fmt.Println("  - EntryPointsByType: Constructor, External, L1Handler")
	fmt.Println("  - Prime: Field prime number")
	fmt.Println("  - CompilerVersion: Compiler version used")
	fmt.Println("  - Hints: Execution hints")

	// Check if a CASM file exists in the current directory for demo
	testFile := "test_contract.casm.json"
	if _, err := os.Stat(testFile); err == nil {
		casmClass, err := contracts.UnmarshalCasmClass(testFile)
		if err != nil {
			log.Printf("Error loading test file: %v", err)
		} else {
			fmt.Printf("\n✅ Successfully loaded CASM class:\n")
			fmt.Printf("  Compiler Version: %s\n", casmClass.CompilerVersion)
			fmt.Printf("  Bytecode Length: %d\n", len(casmClass.ByteCode))
			fmt.Printf("  Entry Points: %d external\n", len(casmClass.EntryPointsByType.External))
		}
	} else {
		fmt.Println("\nNote: Provide a .casm.json file to test loading")
	}
}

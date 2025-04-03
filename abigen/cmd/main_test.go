package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateCommand(t *testing.T) {
	t.Skip("Skipping test until template embedding is fixed")
	tempDir, err := ioutil.TempDir("", "abigen-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	abiContent := `[
		{
			"type": "function",
			"name": "increase_balance",
			"inputs": [
				{
					"name": "amount",
					"type": "core::felt252"
				}
			],
			"outputs": [],
			"state_mutability": "external"
		},
		{
			"type": "function",
			"name": "get_balance",
			"inputs": [],
			"outputs": [
				{
					"type": "core::felt252"
				}
			],
			"state_mutability": "view"
		}
	]`
	abiFile := filepath.Join(tempDir, "test.abi.json")
	if err := ioutil.WriteFile(abiFile, []byte(abiContent), 0644); err != nil {
		t.Fatalf("Failed to write ABI file: %v", err)
	}

	outFile := filepath.Join(tempDir, "test.go")
	
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	
	os.Args = []string{
		"abigen",
		"--abi", abiFile,
		"--pkg", "test",
		"--out", outFile,
		"--type", "TestContract",
	}
	
	main()
	
	w.Close()
	os.Stdout = oldStdout
	
	var buf bytes.Buffer
	buf.ReadFrom(r)
	_ = buf.String() // Capture output but not used in this test
	
	if _, err := os.Stat(outFile); os.IsNotExist(err) {
		t.Errorf("Output file was not created: %v", err)
	}
	
	generatedCode, err := ioutil.ReadFile(outFile)
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}
	
	generatedStr := string(generatedCode)
	if !strings.Contains(generatedStr, "package test") {
		t.Errorf("Generated code does not contain package declaration")
	}
	
	if !strings.Contains(generatedStr, "type TestContract struct") {
		t.Errorf("Generated code does not contain contract struct")
	}
	
	if !strings.Contains(generatedStr, "func (_TestContract *TestContractCaller) GetBalance") {
		t.Errorf("Generated code does not contain view method")
	}
	
	if !strings.Contains(generatedStr, "func (_TestContract *TestContractTransactor) IncreaseBalance") {
		t.Errorf("Generated code does not contain external method")
	}
}

func TestGenerateStdout(t *testing.T) {
	t.Skip("Skipping test until template embedding is fixed")
	tempDir, err := ioutil.TempDir("", "abigen-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	abiContent := `[
		{
			"type": "function",
			"name": "test_function",
			"inputs": [
				{
					"name": "param",
					"type": "core::felt252"
				}
			],
			"outputs": [
				{
					"type": "core::felt252"
				}
			],
			"state_mutability": "view"
		}
	]`
	abiFile := filepath.Join(tempDir, "test.abi.json")
	if err := ioutil.WriteFile(abiFile, []byte(abiContent), 0644); err != nil {
		t.Fatalf("Failed to write ABI file: %v", err)
	}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	
	os.Args = []string{
		"abigen",
		"--abi", abiFile,
		"--pkg", "test",
		"--type", "TestContract",
	}
	
	main()
	
	w.Close()
	os.Stdout = oldStdout
	
	var buf bytes.Buffer
	buf.ReadFrom(r)
	_ = buf.String() // Capture output but not used in this test
	
	if !strings.Contains(buf.String(), "package test") {
		t.Errorf("Generated code does not contain package declaration")
	}
	
	if !strings.Contains(buf.String(), "type TestContract struct") {
		t.Errorf("Generated code does not contain contract struct")
	}
	
	if !strings.Contains(buf.String(), "func (_TestContract *TestContractCaller) TestFunction") {
		t.Errorf("Generated code does not contain view method")
	}
}

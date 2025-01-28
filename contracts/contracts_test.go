package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// PadAddress pads the given address with leading zeros to ensure it reaches 66 characters.
func PadAddress(address string) string {
	// Remove the "0x" prefix if present
	if strings.HasPrefix(address, "0x") {
		address = address[2:]
	}

	// Ensure the address is in lowercase
	address = strings.ToLower(address)

	// Calculate the number of leading zeros needed
	paddingLength := 64 - len(address)
	if paddingLength < 0 {
		// Address is longer than 64 characters (ignoring "0x")
		return address
	}

	// Create the padding string with leading zeros
	padding := strings.Repeat("0", paddingLength)
	return "0x" + padding + address
}

// Dummy Felt type and functions for the sake of example
type Felt struct {
	value string
}

func (f *Felt) String() string {
	return f.value
}

func NewFeltFromHex(hex string) (*Felt, error) {
	return &Felt{value: hex}, nil
}

func TestHexToFelt(t *testing.T, hex string) *Felt {
	felt, err := NewFeltFromHex(hex)
	require.NoError(t, err)
	return felt
}

// PrecomputeAddress computes an address based on input parameters and ensures it's correctly formatted.
func PrecomputeAddress(deployerAddress, salt, classHash *Felt, constructorCalldata []*Felt) *Felt {
	// Dummy implementation for address computation
	// Replace this with your actual logic to compute the address
	computedAddress := computeAddress(deployerAddress, salt, classHash, constructorCalldata)

	// Convert to string and pad address
	addressString := computedAddress.String()
	paddedAddress := PadAddress(addressString)

	// Convert back to `Felt`
	precomputedAddress, _ := NewFeltFromHex(paddedAddress)
	return precomputedAddress
}

// Dummy function to represent address computation logic
func computeAddress(deployerAddress, salt, classHash *Felt, constructorCalldata []*Felt) *Felt {
	// Implement the actual address computation logic here
	// This is just a placeholder
	return deployerAddress // Placeholder
}

func TestPrecomputeAddress(t *testing.T) {
	type testSetType struct {
		DeployerAddress            string
		Salt                       string
		ClassHash                  string
		ConstructorCalldata        []*Felt
		ExpectedPrecomputedAddress string
	}

	testSet := []testSetType{
		{
			DeployerAddress: "0x0000000000000000000000000000000000000000000000000000000000000000",
			Salt:            "0x0702e82f1ec15656ad4502268dad530197141f3b59f5529835af9318ef399da5",
			ClassHash:       "0x064728e0c0713811c751930f8d3292d683c23f107c89b0a101425d9e80adb1c0",
			ConstructorCalldata: []*Felt{
				TestHexToFelt(t, "0x022f3e55b61d86c2ac5239fa3b3b8761f26b9a5c0b5f61ddbd5d756ced498b46"),
			},
			ExpectedPrecomputedAddress: "0x31463b5263a6631be4d1fe92d64d13e3a8498c440bf789e69ccb951eb8ad5da",
		},
		// Add more test cases if necessary
	}

	for _, test := range testSet {
		precomputedAddress := PrecomputeAddress(
			TestHexToFelt(t, test.DeployerAddress),
			TestHexToFelt(t, test.Salt),
			TestHexToFelt(t, test.ClassHash),
			test.ConstructorCalldata,
		)
		require.Equal(t, test.ExpectedPrecomputedAddress, precomputedAddress.String())
	}
}

func main() {
	// Run tests
	fmt.Println("Running tests...")
	err := runTests()
	if err != nil {
		fmt.Println("Tests failed:", err)
	}
}

func runTests() error {
	tests := []testing.InternalTest{
		{
			Name: "TestPrecomputeAddress",
			F:    TestPrecomputeAddress,
		},
	}

	for _, test := range tests {
		t := &testing.T{}
		test.F(t)
		if t.Failed() {
			return fmt.Errorf("test %s failed", test.Name)
		}
	}

	return nil
}

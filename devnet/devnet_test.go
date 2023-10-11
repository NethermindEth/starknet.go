package devnet

import (
	"math/big"
	"strings"
	"testing"

	"github.com/NethermindEth/starknet.go/utils"
)

// TestDevnet_IsAlive tests the IsAlive method of the Devnet struct.
//
// This function checks if the Devnet is alive by creating a new instance of the Devnet struct and calling the IsAlive method.
// It then uses the Fatalf method from the testing package to fail the test if the Devnet is not alive.
//
// Parameters:
// - t: is the testing.T instance for running the test
// Returns:
//   none
func TestDevnet_IsAlive(t *testing.T) {
	d := NewDevNet()
	if !d.IsAlive() {
		t.Fatalf("Devnet should be alive!")
	}
}

// TestDevnet_Accounts tests the Accounts function of the Devnet struct.
//
// It verifies that reading the accounts should succeed and that the returned
// account addresses are valid.
//
// Parameters:
//  - t: is the testing.T instance for running the test
// Returns:
//  none
func TestDevnet_Accounts(t *testing.T) {
	d := NewDevNet()
	accounts, err := d.Accounts()
	if err != nil {
		t.Fatalf("Reading account should succeed, instead: %v", err)
	}
	if len(accounts) == 0 || !strings.HasPrefix(accounts[0].Address, "0x") {
		t.Fatal("should return valid account addresses")
	}
}

// TestDevnet_FeeToken tests the FeeToken function of the Devnet struct.
//
// The function retrieves the fee token from the Devnet instance and checks that
// it matches the expected ETH address.
//
// Parameters:
// - t: is the testing.T instance for running the test
// Returns:
//   none
func TestDevnet_FeeToken(t *testing.T) {
	d := NewDevNet()
	token, err := d.FeeToken()
	if err != nil {
		t.Fatalf("Reading token should succeed, instead: %v", err)
	}
	if token.Address.String() != DevNetETHAddress {
		t.Fatalf("devnet ETH address, instead %s", token.Address.String())
	}
}

// TestDevnet_Mint is a test function that tests the Mint method of the Devnet struct.
//
// It initializes a new Devnet instance and sets the amount to 1000000000000000000.
// Then it calls the Mint method with a test hexadecimal value and the amount.
// If an error occurs during the Mint method call, it fails the test with the error message.
// If the NewBalance returned by the Mint method is less than the amount, it fails the test with an error message.
//
// Parameters:
// - t: is the testing.T instance for running the test
// Returns:
//   none
func TestDevnet_Mint(t *testing.T) {
	d := NewDevNet()
	amount := big.NewInt(int64(1000000000000000000))
	resp, err := d.Mint(utils.TestHexToFelt(t, "0x1"), amount)
	if err != nil {
		t.Fatalf("Minting ETH should succeed, instead: %v", err)
	}
	if resp.NewBalance.Cmp(amount) < 0 {
		t.Fatalf("ETH should be higher than the last mint, instead: %d", resp.NewBalance)
	}
}

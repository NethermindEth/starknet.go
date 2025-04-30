package devnet

import (
	"math/big"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/NethermindEth/starknet.go/internal"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
)

var (
	// the environment for the test, needs to be "devnet"
	testEnv = ""
)

// TestMain is the main test function for the package, checks configuration for the environment to use.
//
// It initializes the test environment and runs the test cases.
//
// Parameters:
//   - m: is the testing.M parameter
//
// Returns:
//
//	none
func TestMain(m *testing.M) {
	testEnv = internal.LoadEnv()

	os.Exit(m.Run())
}

// TestDevnet_IsAlive tests the IsAlive method of the Devnet struct.
//
// This function checks if the Devnet is alive by creating a new instance of the Devnet struct and calling the IsAlive method.
// It then uses the Fatalf method from the testing package to fail the test if the Devnet is not alive.
//
// Parameters:
//   - t: is the testing.T instance for running the test
//
// Returns:
//
//	none
func TestDevnet_IsAlive(t *testing.T) {
	if testEnv != "devnet" {
		t.Skip("Skipping test as it requires a devnet environment")
	}

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
//   - t: is the testing.T instance for running the test
//
// Returns:
//
//	none
func TestDevnet_Accounts(t *testing.T) {
	if testEnv != "devnet" {
		t.Skip("Skipping test as it requires a devnet environment")
	}

	d := NewDevNet()
	accounts, err := d.Accounts()
	if err != nil {
		t.Fatalf("Reading account should succeed, instead: %v", err)
	}
	if len(accounts) == 0 || !strings.HasPrefix(accounts[0].Address, "0x") {
		t.Fatal("should return valid account addresses")
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
//   - t: is the testing.T instance for running the test
//
// Returns:
//
//	none
func TestDevnet_Mint(t *testing.T) {
	if testEnv != "devnet" {
		t.Skip("Skipping test as it requires a devnet environment")
	}

	d := NewDevNet()
	amount := big.NewInt(int64(1000000000000000000))
	resp, err := d.Mint(internalUtils.TestHexToFelt(t, "0x1"), amount)
	if err != nil {
		t.Fatalf("Minting ETH should succeed, instead: %v", err)
	}
	balance, _ := (strconv.ParseInt(resp.NewBalance, 10, 64))
	if balance < 0 {
		t.Fatalf("ETH should be higher than the last mint, instead: %d", balance)
	}
}

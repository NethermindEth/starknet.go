package devnet

import (
	"math/big"
	"strings"
	"testing"

	"github.com/NethermindEth/starknet.go/utils"
)

func TestDevnet_IsAlive(t *testing.T) {
	d := NewDevNet()
	if !d.IsAlive() {
		t.Fatalf("Devnet should be alive!")
	}
}

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

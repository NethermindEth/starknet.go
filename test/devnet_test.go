package test

import (
	"strings"
	"testing"

	"github.com/smartcontractkit/caigo/types"
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
	if token.Address.String() != "0x062230ea046a9a5fbc261ac77d03c8d41e5d442db2284587570ab46455fd2488" {
		t.Fatalf("devnet ETH address, instead %s", token.Address.String())
	}
}

func TestDevnet_Mint(t *testing.T) {
	d := NewDevNet()
	resp, err := d.Mint(types.StrToFelt("0x1"), 1000000000000000000)
	if err != nil {
		t.Fatalf("Minting ETH should succeed, instead: %v", err)
	}
	if resp.NewBalance < 1000000000000000000 {
		t.Fatalf("ETH should be higher than the last mint, instead: %d", resp.NewBalance)
	}
}

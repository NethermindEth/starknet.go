package caigo

import (
	"math/big"
	"testing"
)

func TestDeployAccountAddress(t *testing.T) {
	constructorCalldata := []string{"0x79ad1d75aa9448c9f280a1ce5397c8510a63241f9a2de2de2b44133dd78b0af"}
	salt := "0x66cb284556e80c4351bb25078838eb2a07bdfa48f06e6572efd5dcecd8b76de"
	address, _ := big.NewInt(0).SetString("0x43cec55061e032c92e2a629cf48cbefa10ff74fb14b98c6679dade836fe8739", 0)
	classHash := "0x1fac3074c9d5282f0acc5c69a4781a1c711efea5e73c550c5d9fb253cf7fd3d"
	daddress, err := ContractAddress("0x0", salt, classHash, constructorCalldata)
	if err != nil {
		t.Fatal("error should be nil", err)
	}
	if address.Cmp(daddress) != 0 {
		t.Fatalf("the 2 addresses should match, instead 0x%s vs 0x%s", daddress.Text(16), address.Text(16))
	}
}

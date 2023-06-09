package types

import (
	"math/big"
	"testing"
)

func TestUint256_Big(t *testing.T) {
	low := StrToFelt("169642200937843940828817578484658339840")
	high := StrToFelt("5293408896007833759405984958868987187")

	uint256, err := NewUint256(low, high)
	if err != nil {
		t.Fatalf("failed to convert to uint256: %v", err)
	}
	got := uint256.Big()
	want, _ := new(big.Int).SetString("1801253708213897481358155843874670753928272635927844366765878409841828954112", 10)

	if got.Cmp(want) != 0 {
		t.Fatalf("uint256 value mismatch, want: %s, got: %s", want.String(), got.String())
	}

	uint256, err = Uint256FromBig(got)
	if err != nil {
		t.Fatalf("failed to convert to uin256: %v", err)
	}
	if uint256.Low.Big().Cmp(low.Big()) != 0 {
		t.Fatalf("uint256 low value mismatch, want: %s, got: %s", low.String(), uint256.Low.String())
	}
	if uint256.High.Big().Cmp(high.Big()) != 0 {
		t.Fatalf("uint256 high value mismatch, want: %s, got: %s", high.String(), uint256.High.String())
	}

	low = StrToFelt("0")
	high = StrToFelt("0")

	uint256, err = NewUint256(low, high)
	if err != nil {
		t.Fatalf("failed to convert to uint256: %v", err)
	}
	got = uint256.Big()
	want, _ = new(big.Int).SetString("0", 10)

	if got.Cmp(want) != 0 {
		t.Fatalf("uint256 value mismatch, want: %s, got: %s", want.String(), got.String())
	}

	uint256, err = Uint256FromBig(got)
	if err != nil {
		t.Fatalf("failed to convert to uin256: %v", err)
	}

	if uint256.Low.Big().Cmp(low.Big()) != 0 {
		t.Fatalf("uint256 low value mismatch, want: %s, got: %s", low.String(), uint256.Low.String())
	}
	if uint256.High.Big().Cmp(high.Big()) != 0 {
		t.Fatalf("uint256 high value mismatch, want: %s, got: %s", high.String(), uint256.High.String())
	}

	max := new(big.Int).Lsh(big.NewInt(1), 128)
	low = BigToFelt(max)
	high = StrToFelt("0")

	if _, err = NewUint256(low, high); err.Error() != "invalid low felt value" {
		t.Error("expected invalid low felt value")
	}
	low = BigToFelt(new(big.Int).Sub(max, big.NewInt(1)))
	if _, err = NewUint256(low, high); err != nil {
		t.Errorf("failed to convert to uint256: %v", err)
	}

	low = StrToFelt("0")
	high = BigToFelt(new(big.Int).Lsh(big.NewInt(1), 128))

	if _, err = NewUint256(low, high); err.Error() != "invalid high felt value" {
		t.Error("expected invalid high felt value")
	}
	high = BigToFelt(new(big.Int).Sub(max, big.NewInt(1)))
	if _, err = NewUint256(low, high); err != nil {
		t.Errorf("failed to convert to uint256: %v", err)
	}

	max = new(big.Int).Lsh(big.NewInt(1), 256)
	if _, err = Uint256FromBig(max); err.Error() != "invalid uint256 value" {
		t.Error("expected invalid uint256 value")
	}
	max = new(big.Int).Sub(max, big.NewInt(1))
	if _, err = Uint256FromBig(max); err != nil {
		t.Errorf("failed to convert to uint256: %v", err)
	}
}

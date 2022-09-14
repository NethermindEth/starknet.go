package types

import (
	"math/big"
	"testing"
)

func TestUint256_Big(t *testing.T) {
	low := StrToFelt("169642200937843940828817578484658339840")
	high := StrToFelt("5293408896007833759405984958868987187")

	uint256 := NewUint256(low, high)
	got := uint256.Big()
	want, _ := new(big.Int).SetString("1801253708213897481358155843874670753928272635927844366765878409841828954112", 10)

	if got.Cmp(want) != 0 {
		t.Fatalf("uint256 value mismatch, want: %s, got: %s", want.String(), got.String())
	}

	uint256 = Uint256FromBig(got)
	if uint256.Low.Cmp(low.Int) != 0 {
		t.Fatalf("uint256 low value mismatch, want: %s, got: %s", low.String(), uint256.Low.String())
	}

	uint256 = Uint256FromBig(got)
	if uint256.High.Cmp(high.Int) != 0 {
		t.Fatalf("uint256 high value mismatch, want: %s, got: %s", high.String(), uint256.High.String())
	}

	low = StrToFelt("0")
	high = StrToFelt("0")

	uint256 = NewUint256(low, high)
	got = uint256.Big()
	want, _ = new(big.Int).SetString("0", 10)

	if got.Cmp(want) != 0 {
		t.Fatalf("uint256 value mismatch, want: %s, got: %s", want.String(), got.String())
	}

	uint256 = Uint256FromBig(got)
	if uint256.Low.Cmp(low.Int) != 0 {
		t.Fatalf("uint256 low value mismatch, want: %s, got: %s", low.String(), uint256.Low.String())
	}

	uint256 = Uint256FromBig(got)
	if uint256.High.Cmp(high.Int) != 0 {
		t.Fatalf("uint256 high value mismatch, want: %s, got: %s", high.String(), uint256.High.String())
	}
}

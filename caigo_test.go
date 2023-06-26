package caigo

import (
	"crypto/elliptic"
	"fmt"
	"math/big"
	"testing"

	"github.com/smartcontractkit/caigo/types"
)

func BenchmarkSignatureVerify(b *testing.B) {
	private, _ := Curve.GetRandomPrivateKey()
	x, y, _ := Curve.PrivateToPoint(private)

	hash, _ := Curve.PedersenHash(
		[]*big.Int{
			types.HexToBN("0x7f15c38ea577a26f4f553282fcfe4f1feeb8ecfaad8f221ae41abf8224cbddd"),
			types.HexToBN("0x7f15c38ea577a26f4f553282fcfe4f1feeb8ecfaad8f221ae41abf8224cbdde"),
		})

	r, s, _ := Curve.Sign(hash, private)

	b.Run(fmt.Sprintf("sign_input_size_%d", hash.BitLen()), func(b *testing.B) {
		Curve.Sign(hash, private)
	})
	b.Run(fmt.Sprintf("verify_input_size_%d", hash.BitLen()), func(b *testing.B) {
		Curve.Verify(hash, r, s, x, y)
	})
}

func TestGeneral_ComputeHashOnElements(t *testing.T) {
	hashEmptyArray, err := Curve.ComputeHashOnElements([]*big.Int{})
	expectedHashEmmptyArray := types.HexToBN("0x49ee3eba8c1600700ee1b87eb599f16716b0b1022947733551fde4050ca6804")
	if err != nil {
		t.Errorf("Could no hash an empty array %v\n", err)
	}
	if hashEmptyArray.Cmp(expectedHashEmmptyArray) != 0 {
		t.Errorf("Hash empty array wrong value. Expected %v got %v\n", expectedHashEmmptyArray, hashEmptyArray)
	}

	hashFilledArray, err := Curve.ComputeHashOnElements([]*big.Int{
		big.NewInt(123782376),
		big.NewInt(213984),
		big.NewInt(128763521321),
	})
	expectedHashFilledArray := types.HexToBN("0x7b422405da6571242dfc245a43de3b0fe695e7021c148b918cd9cdb462cac59")

	if err != nil {
		t.Errorf("Could no hash an array with values %v\n", err)
	}
	if hashFilledArray.Cmp(expectedHashFilledArray) != 0 {
		t.Errorf("Hash filled array wrong value. Expected %v got %v\n", expectedHashFilledArray, hashFilledArray)
	}
}

func TestGeneral_HashAndSign(t *testing.T) {
	hashy, err := Curve.HashElements([]*big.Int{
		big.NewInt(1953658213),
		big.NewInt(126947999705460),
		big.NewInt(1953658213),
	})
	if err != nil {
		t.Errorf("Hasing elements: %v\n", err)
	}

	priv, _ := Curve.GetRandomPrivateKey()
	x, y, err := Curve.PrivateToPoint(priv)
	if err != nil {
		t.Errorf("Could not convert random private key to point: %v\n", err)
	}

	r, s, err := Curve.Sign(hashy, priv)
	if err != nil {
		t.Errorf("Could not convert gen signature: %v\n", err)
	}

	if !Curve.Verify(hashy, r, s, x, y) {
		t.Errorf("Verified bad signature %v %v\n", r, s)
	}
}

func TestGeneral_ComputeFact(t *testing.T) {
	testFacts := []struct {
		programHash   *big.Int
		programOutput []*big.Int
		expected      *big.Int
	}{
		{
			programHash:   types.HexToBN("0x114952172aed91e59f870a314e75de0a437ff550e4618068cec2d832e48b0c7"),
			programOutput: []*big.Int{big.NewInt(289)},
			expected:      types.HexToBN("0xe6168c0a865aa80d724ad05627fa65fbcfe4b1d66a586e9f348f461b076072c4"),
		},
		{
			programHash:   types.HexToBN("0x79920d895101ad1fbdea9adf141d8f362fdea9ee35f33dfcd07f38e4a589bab"),
			programOutput: []*big.Int{types.StrToBig("2754806153357301156380357983574496185342034785016738734224771556919270737441")},
			expected:      types.HexToBN("0x1d174fa1443deea9aab54bbca8d9be308dd14a0323dd827556c173bd132098db"),
		},
	}

	for _, tt := range testFacts {
		hash := types.ComputeFact(tt.programHash, tt.programOutput)
		if hash.Cmp(tt.expected) != 0 {
			t.Errorf("Fact does not equal ex %v %v\n", hash, tt.expected)
		}
	}
}

func TestGeneral_BadSignature(t *testing.T) {
	hash, err := Curve.PedersenHash([]*big.Int{types.HexToBN("0x12773"), types.HexToBN("0x872362")})
	if err != nil {
		t.Errorf("Hashing err: %v\n", err)
	}

	priv, _ := Curve.GetRandomPrivateKey()
	x, y, err := Curve.PrivateToPoint(priv)
	if err != nil {
		t.Errorf("Could not convert random private key to point: %v\n", err)
	}

	r, s, err := Curve.Sign(hash, priv)
	if err != nil {
		t.Errorf("Could not convert gen signature: %v\n", err)
	}

	badR := new(big.Int).Add(r, big.NewInt(1))
	if Curve.Verify(hash, badR, s, x, y) {
		t.Errorf("Verified bad signature %v %v\n", r, s)
	}

	badS := new(big.Int).Add(s, big.NewInt(1))
	if Curve.Verify(hash, r, badS, x, y) {
		t.Errorf("Verified bad signature %v %v\n", r, s)
	}

	badHash := new(big.Int).Add(hash, big.NewInt(1))
	if Curve.Verify(badHash, r, s, x, y) {
		t.Errorf("Verified bad signature %v %v\n", r, s)
	}
}

func TestGeneral_Signature(t *testing.T) {
	testSignature := []struct {
		private *big.Int
		publicX *big.Int
		publicY *big.Int
		hash    *big.Int
		rIn     *big.Int
		sIn     *big.Int
		raw     string
	}{
		{
			private: types.StrToBig("104397037759416840641267745129360920341912682966983343798870479003077644689"),
			publicX: types.StrToBig("1913222325711601599563860015182907040361852177892954047964358042507353067365"),
			publicY: types.StrToBig("798905265292544287704154888908626830160713383708400542998012716235575472365"),
			hash:    types.StrToBig("2680576269831035412725132645807649347045997097070150916157159360688041452746"),
			rIn:     types.StrToBig("607684330780324271206686790958794501662789535258258105407533051445036595885"),
			sIn:     types.StrToBig("453590782387078613313238308551260565642934039343903827708036287031471258875"),
		},
		{
			hash: types.HexToBN("0x7f15c38ea577a26f4f553282fcfe4f1feeb8ecfaad8f221ae41abf8224cbddd"),
			rIn:  types.StrToBig("2458502865976494910213617956670505342647705497324144349552978333078363662855"),
			sIn:  types.StrToBig("3439514492576562277095748549117516048613512930236865921315982886313695689433"),
			raw:  "04033f45f07e1bd1a51b45fc24ec8c8c9908db9e42191be9e169bfcac0c0d997450319d0f53f6ca077c4fa5207819144a2a4165daef6ee47a7c1d06c0dcaa3e456",
		},
		{
			hash:    types.HexToBN("0x324df642fcc7d98b1d9941250840704f35b9ac2e3e2b58b6a034cc09adac54c"),
			publicX: types.HexToBN("0x4e52f2f40700e9cdd0f386c31a1f160d0f310504fc508a1051b747a26070d10"),
			rIn:     types.StrToBig("2849277527182985104629156126825776904262411756563556603659114084811678482647"),
			sIn:     types.StrToBig("3156340738553451171391693475354397094160428600037567299774561739201502791079"),
		},
	}

	var err error
	for _, tt := range testSignature {
		if tt.raw != "" {
			h, _ := types.HexToBytes(tt.raw)
			tt.publicX, tt.publicY = elliptic.Unmarshal(Curve, h)
		} else if tt.private != nil {
			tt.publicX, tt.publicY, err = Curve.PrivateToPoint(tt.private)
			if err != nil {
				t.Errorf("Could not convert random private key to point: %v\n", err)
			}
		} else if tt.publicX != nil {
			tt.publicY = Curve.GetYCoordinate(tt.publicX)
		}

		if tt.rIn == nil && tt.private != nil {
			tt.rIn, tt.sIn, err = Curve.Sign(tt.hash, tt.private)
			if err != nil {
				t.Errorf("Could not sign good hash: %v\n", err)
			}
		}

		if !Curve.Verify(tt.hash, tt.rIn, tt.sIn, tt.publicX, tt.publicY) {
			t.Errorf("successful signature did not verify\n")
		}
	}
}

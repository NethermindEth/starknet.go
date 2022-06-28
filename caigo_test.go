package caigo

import (
	"crypto/elliptic"
	"fmt"
	"math/big"
	"testing"
)

var caigoCurve StarkCurve

func init() {
	var err error
	caigoCurve, err = SC(WithConstants("./pedersen_params.json"))
	if err != nil {
		panic(err.Error())
	}
}

func BenchmarkSignatureVerify(b *testing.B) {
	private, _ := caigoCurve.GetRandomPrivateKey()
	x, y, _ := caigoCurve.PrivateToPoint(private)

	hash, _ := caigoCurve.PedersenHash(
		[]*big.Int{
			HexToBN("0x7f15c38ea577a26f4f553282fcfe4f1feeb8ecfaad8f221ae41abf8224cbddd"),
			HexToBN("0x7f15c38ea577a26f4f553282fcfe4f1feeb8ecfaad8f221ae41abf8224cbdde"),
		})

	r, s, _ := caigoCurve.Sign(hash, private)

	b.Run(fmt.Sprintf("sign_input_size_%d", hash.BitLen()), func(b *testing.B) {
		caigoCurve.Sign(hash, private)
	})
	b.Run(fmt.Sprintf("verify_input_size_%d", hash.BitLen()), func(b *testing.B) {
		caigoCurve.Verify(hash, r, s, x, y)
	})
}

func TestHashAndSign(t *testing.T) {
	hashy, err := caigoCurve.HashElements([]*big.Int{
		big.NewInt(1953658213),
		big.NewInt(126947999705460),
		big.NewInt(1953658213),
	})
	if err != nil {
		t.Errorf("Hasing elements: %v\n", err)
	}

	priv, _ := caigoCurve.GetRandomPrivateKey()
	x, y, err := caigoCurve.PrivateToPoint(priv)
	if err != nil {
		t.Errorf("Could not convert random private key to point: %v\n", err)
	}

	r, s, err := caigoCurve.Sign(hashy, priv)
	if err != nil {
		t.Errorf("Could not convert gen signature: %v\n", err)
	}

	if !caigoCurve.Verify(hashy, r, s, x, y) {
		t.Errorf("Verified bad signature %v %v\n", r, s)
	}
}

func TestComputeFact(t *testing.T) {
	testFacts := []struct {
		programHash			*big.Int
		programOutput		[]*big.Int
		expected			*big.Int
	}{
		{
			programHash: HexToBN("0x114952172aed91e59f870a314e75de0a437ff550e4618068cec2d832e48b0c7"),
			programOutput: []*big.Int{big.NewInt(289)},
			expected: HexToBN("0xe6168c0a865aa80d724ad05627fa65fbcfe4b1d66a586e9f348f461b076072c4"),
		},
		{
			programHash: HexToBN("0x79920d895101ad1fbdea9adf141d8f362fdea9ee35f33dfcd07f38e4a589bab"),
			programOutput: []*big.Int{StrToBig("2754806153357301156380357983574496185342034785016738734224771556919270737441")},
			expected: HexToBN("0x1d174fa1443deea9aab54bbca8d9be308dd14a0323dd827556c173bd132098db"),
		},
	}

	for _, tt := range testFacts {
		hash := ComputeFact(tt.programHash, tt.programOutput)
		if hash.Cmp(tt.expected) != 0 {
			t.Errorf("Fact does not equal ex %v %v\n", hash, tt.expected)
		}
	}
}

func TestBadSignature(t *testing.T) {
	hash, err := caigoCurve.PedersenHash([]*big.Int{HexToBN("0x12773"), HexToBN("0x872362")})
	if err != nil {
		t.Errorf("Hashing err: %v\n", err)
	}

	priv, _ := caigoCurve.GetRandomPrivateKey()
	x, y, err := caigoCurve.PrivateToPoint(priv)
	if err != nil {
		t.Errorf("Could not convert random private key to point: %v\n", err)
	}

	r, s, err := caigoCurve.Sign(hash, priv)
	if err != nil {
		t.Errorf("Could not convert gen signature: %v\n", err)
	}

	badR := new(big.Int).Add(r, big.NewInt(1))
	if caigoCurve.Verify(hash, badR, s, x, y) {
		t.Errorf("Verified bad signature %v %v\n", r, s)
	}

	badS := new(big.Int).Add(s, big.NewInt(1))
	if caigoCurve.Verify(hash, r, badS, x, y) {
		t.Errorf("Verified bad signature %v %v\n", r, s)
	}

	badHash := new(big.Int).Add(hash, big.NewInt(1))
	if caigoCurve.Verify(badHash, r, s, x, y) {
		t.Errorf("Verified bad signature %v %v\n", r, s)
	}
}

func TestSignature(t *testing.T) {
	testSignature := []struct {
		private			*big.Int
		publicX			*big.Int
		publicY			*big.Int
		hash			*big.Int
		rIn				*big.Int
		sIn				*big.Int
		raw				string
	}{
		{
			private: StrToBig("104397037759416840641267745129360920341912682966983343798870479003077644689"),
			publicX: StrToBig("1913222325711601599563860015182907040361852177892954047964358042507353067365"),
			publicY: StrToBig("798905265292544287704154888908626830160713383708400542998012716235575472365"),
			hash: StrToBig("2680576269831035412725132645807649347045997097070150916157159360688041452746"),
			rIn: StrToBig("607684330780324271206686790958794501662789535258258105407533051445036595885"),
			sIn: StrToBig("453590782387078613313238308551260565642934039343903827708036287031471258875"),
		},
		{
			hash: HexToBN("0x7f15c38ea577a26f4f553282fcfe4f1feeb8ecfaad8f221ae41abf8224cbddd"),
			rIn: StrToBig("2458502865976494910213617956670505342647705497324144349552978333078363662855"),
			sIn: StrToBig("3439514492576562277095748549117516048613512930236865921315982886313695689433"),
			raw: "04033f45f07e1bd1a51b45fc24ec8c8c9908db9e42191be9e169bfcac0c0d997450319d0f53f6ca077c4fa5207819144a2a4165daef6ee47a7c1d06c0dcaa3e456",
		},
		{
			hash: HexToBN("0x324df642fcc7d98b1d9941250840704f35b9ac2e3e2b58b6a034cc09adac54c"),
			publicX: HexToBN("0x4e52f2f40700e9cdd0f386c31a1f160d0f310504fc508a1051b747a26070d10"),
			rIn: StrToBig("2849277527182985104629156126825776904262411756563556603659114084811678482647"),
			sIn: StrToBig("3156340738553451171391693475354397094160428600037567299774561739201502791079"),
		},
	}

	var err error
	for _, tt := range testSignature {
		if tt.raw != "" {
			h, _ := HexToBytes(tt.raw)
			tt.publicX, tt.publicY = elliptic.Unmarshal(curve, h)
		} else if tt.private != nil {
			tt.publicX, tt.publicY, err = caigoCurve.PrivateToPoint(tt.private)
			if err != nil {
				t.Errorf("Could not convert random private key to point: %v\n", err)
			}
		} else if tt.publicX != nil {
			tt.publicY = caigoCurve.GetYCoordinate(tt.publicX)
		}

		if tt.rIn == nil  && tt.private != nil {
			tt.rIn, tt.sIn, err = caigoCurve.Sign(tt.hash, tt.private)
			if err != nil {
				t.Errorf("Could not sign good hash: %v\n", err)
			}
		}
		
		if !caigoCurve.Verify(tt.hash, tt.rIn, tt.sIn, tt.publicX, tt.publicY) {
			t.Errorf("successful signature did not verify\n")
		}
	}
}

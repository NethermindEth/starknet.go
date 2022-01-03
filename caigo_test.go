package caigo

import (
	"testing"
	"math/big"
	"crypto/ecdsa"
	"crypto/elliptic"
)

func TestInitCurveWithConstants(t *testing.T) {
	curve, err := SCWithConstants("./pedersen_params.json")
	if err != nil {
		t.Errorf("Could not init with constant points: %v\n", err)
	}

	hash := HexToBN("0x7f15c38ea577a26f4f553282fcfe4f1feeb8ecfaad8f221ae41abf8224cbddd")
	r, _ := new(big.Int).SetString("2458502865976494910213617956670505342647705497324144349552978333078363662855", 10)
	s, _ := new(big.Int).SetString("3439514492576562277095748549117516048613512930236865921315982886313695689433", 10)

	h, _ := HexToBytes("04033f45f07e1bd1a51b45fc24ec8c8c9908db9e42191be9e169bfcac0c0d997450319d0f53f6ca077c4fa5207819144a2a4165daef6ee47a7c1d06c0dcaa3e456")
	x, y := elliptic.Unmarshal(curve, h)
	pub := ecdsa.PublicKey{
		Curve: curve,
		X:     x,
		Y:     y,
	}
	if !Verify(hash, r, s, pub, curve) {
		t.Errorf("successful signature did not verify\n")
	}
}

func TestDivMod(t *testing.T) {
	curve := SC()
	inX, _ := new(big.Int).SetString("311379432064974854430469844112069886938521247361583891764940938105250923060", 10)
	inY, _ := new(big.Int).SetString("621253665351494585790174448601059271924288186997865022894315848222045687999", 10)
	DIVMODRES, _ := new(big.Int).SetString("2577265149861519081806762825827825639379641276854712526969977081060187505740", 10)

	divR := DivMod(inX, inY, curve.P)
	if divR.Cmp(DIVMODRES) != 0 {
		t.Errorf("DivMod Res %v does not == expected %v\n", divR, DIVMODRES)
	}
}

func TestAdd(t *testing.T) {
	curve := SC()
	pub0, _ := new(big.Int).SetString("1468732614996758835380505372879805860898778283940581072611506469031548393285", 10)
	pub1, _ := new(big.Int).SetString("1402551897475685522592936265087340527872184619899218186422141407423956771926", 10)
	EXPX, _ := new(big.Int).SetString("2573054162739002771275146649287762003525422629677678278801887452213127777391", 10)
	EXPY, _ := new(big.Int).SetString("3086444303034188041185211625370405120551769541291810669307042006593736192813", 10)

	resX, resY := curve.Add(curve.Gx, curve.Gy, pub0, pub1)
	if resX.Cmp(EXPX) != 0 {
		t.Errorf("ResX %v does not == expected %v\n", resX, EXPX)

	}
	if resY.Cmp(EXPY) != 0 {
		t.Errorf("ResY %v does not == expected %v\n", resY, EXPY)
	}
}

func TestMultAir(t *testing.T) {
	curve := SC()
	ry, _ := new(big.Int).SetString("2458502865976494910213617956670505342647705497324144349552978333078363662855", 10)
	pubx, _ := new(big.Int).SetString("1468732614996758835380505372879805860898778283940581072611506469031548393285", 10)
	puby, _ := new(big.Int).SetString("1402551897475685522592936265087340527872184619899218186422141407423956771926", 10)
	resX, _ := new(big.Int).SetString("182543067952221301675635959482860590467161609552169396182763685292434699999", 10)
	resY, _ := new(big.Int).SetString("3154881600662997558972388646773898448430820936643060392452233533274798056266", 10)

	x, y, err := curve.MimicEcMultAir(ry, pubx, puby, curve.Gx, curve.Gy)
	if err != nil {
		t.Errorf("MultAirERR %v\n", err)
	}

	if x.Cmp(resX) != 0 {
		t.Errorf("ResX %v does not == expected %v\n", x, resX)

	}
	if y.Cmp(resY) != 0 {
		t.Errorf("ResY %v does not == expected %v\n", y, resY)
	}
}

func TestGetY(t *testing.T) {
	curve := SC()
	h, _ := HexToBytes("04033f45f07e1bd1a51b45fc24ec8c8c9908db9e42191be9e169bfcac0c0d997450319d0f53f6ca077c4fa5207819144a2a4165daef6ee47a7c1d06c0dcaa3e456")
	x, y := elliptic.Unmarshal(curve, h)

	yout := curve.GetYCoordinate(x)

	if y.Cmp(yout) != 0 {
		t.Errorf("Derived Y %v does not == expected %v\n", yout, y)
	}
}

func TestVerifySignature(t *testing.T) {
	curve := SC()
	hash := HexToBN("0x7f15c38ea577a26f4f553282fcfe4f1feeb8ecfaad8f221ae41abf8224cbddd")
	r, _ := new(big.Int).SetString("2458502865976494910213617956670505342647705497324144349552978333078363662855", 10)
	s, _ := new(big.Int).SetString("3439514492576562277095748549117516048613512930236865921315982886313695689433", 10)

	h, _ := HexToBytes("04033f45f07e1bd1a51b45fc24ec8c8c9908db9e42191be9e169bfcac0c0d997450319d0f53f6ca077c4fa5207819144a2a4165daef6ee47a7c1d06c0dcaa3e456")
	x, y := elliptic.Unmarshal(curve, h)
	pub := ecdsa.PublicKey{
		Curve: curve,
		X:     x,
		Y:     y,
	}
	if !Verify(hash, r, s, pub, curve) {
		t.Errorf("successful signature did not verify\n")
	}
}

func TestUIVerifySignature(t *testing.T) {
	curve := SC()
	hash := HexToBN("0x324df642fcc7d98b1d9941250840704f35b9ac2e3e2b58b6a034cc09adac54c")
	r, _ := new(big.Int).SetString("2849277527182985104629156126825776904262411756563556603659114084811678482647", 10)
	s, _ := new(big.Int).SetString("3156340738553451171391693475354397094160428600037567299774561739201502791079", 10)

	pub := XToPubKey("0x4e52f2f40700e9cdd0f386c31a1f160d0f310504fc508a1051b747a26070d10")

	if !Verify(hash, r, s, pub, curve) {
		t.Errorf("successful signature did not verify\n")
	}
}
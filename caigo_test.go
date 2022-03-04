package caigo

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"math/big"
	"testing"
)

func TestHashAndSign(t *testing.T) {
	curve, err := SCWithConstants("./pedersen_params.json")
	if err != nil {
		t.Errorf("Could not init with constant points: %v\n", err)
	}

	keys := []string{"true", "submit_l1", "true"}
	var bigs []*big.Int
	for _, key := range keys {
		bigs = append(bigs, UTF8StrToBig(key))
	}

	hashy, err := curve.HashElements(bigs)

	priv, _ := curve.GetRandomPrivateKey()

	r, s, err := curve.Sign(hashy, priv)
	if err != nil {
		t.Errorf("Could not convert gen signature: %v\n", err)
	}

	x, y, err := curve.PrivateToPoint(priv)
	if err != nil {
		t.Errorf("Could not convert random private key to point: %v\n", err)
	}

	if !curve.Verify(hashy, r, s, x, y) {
		t.Errorf("Verified bad signature %v %v\n", r, s)
	}
}

func TestComputeFact(t *testing.T) {
	progHash := "0x114952172aed91e59f870a314e75de0a437ff550e4618068cec2d832e48b0c7"
	progOutput := []*big.Int{big.NewInt(289)}
	hash := ComputeFact(HexToBN(progHash), progOutput)

	if hash.Cmp(HexToBN("0xe6168c0a865aa80d724ad05627fa65fbcfe4b1d66a586e9f348f461b076072c4")) != 0 {
		t.Errorf("Fact does not equal ex %v\n", hash)
	}

	progHash = "0x79920d895101ad1fbdea9adf141d8f362fdea9ee35f33dfcd07f38e4a589bab"
	out, _ := new(big.Int).SetString("2754806153357301156380357983574496185342034785016738734224771556919270737441", 10)
	progOutput = []*big.Int{out}
	hash = ComputeFact(HexToBN(progHash), progOutput)

	if hash.Cmp(HexToBN("0x1d174fa1443deea9aab54bbca8d9be308dd14a0323dd827556c173bd132098db")) != 0 {
		t.Errorf("Fact does not equal ex %v %v\n", hash, HexToBN("0x1d174fa1443deea9aab54bbca8d9be308dd14a0323dd827556c173bd132098db"))
	}
}

func TestBadSignature(t *testing.T) {
	curve, err := SCWithConstants("./pedersen_params.json")
	if err != nil {
		t.Errorf("Could not init with constant points: %v\n", err)
	}

	hash, err := curve.PedersenHash([]*big.Int{HexToBN("0x12773"), HexToBN("0x872362")})
	if err != nil {
		t.Errorf("Hashing err: %v\n", err)
	}

	priv, _ := curve.GetRandomPrivateKey()

	x, y, err := curve.PrivateToPoint(priv)
	if err != nil {
		t.Errorf("Could not convert random private key to point: %v\n", err)
	}

	r, s, err := curve.Sign(hash, priv)
	if err != nil {
		t.Errorf("Could not convert gen signature: %v\n", err)
	}
	badR := new(big.Int)
	badR = badR.Add(r, big.NewInt(1))
	if curve.Verify(hash, badR, s, x, y) {
		t.Errorf("Verified bad signature %v %v\n", r, s)
	}

	badS := new(big.Int)
	badS = badS.Add(s, big.NewInt(1))
	if curve.Verify(hash, r, badS, x, y) {
		t.Errorf("Verified bad signature %v %v\n", r, s)
	}

	badHash := new(big.Int)
	badHash = badHash.Add(hash, big.NewInt(1))
	if curve.Verify(badHash, r, s, x, y) {
		t.Errorf("Verified bad signature %v %v\n", r, s)
	}
}

func BenchmarkSignatureVerify(b *testing.B) {
	curve, _ := SCWithConstants("./pedersen_params.json")

	pr, _ := curve.GetRandomPrivateKey()

	prin := new(big.Int)
	prin = prin.Set(pr)
	x, y, _ := curve.PrivateToPoint(prin)

	priv := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: curve,
			X:     x,
			Y:     y,
		},
		D: pr,
	}

	hash, _ := curve.PedersenHash(
		[]*big.Int{
			HexToBN("0x7f15c38ea577a26f4f553282fcfe4f1feeb8ecfaad8f221ae41abf8224cbddd"),
			HexToBN("0x7f15c38ea577a26f4f553282fcfe4f1feeb8ecfaad8f221ae41abf8224cbdde"),
		})

	r, s, _ := curve.Sign(hash, priv.D)

	b.Run(fmt.Sprintf("sign_input_size_%d", hash.BitLen()), func(b *testing.B) {
		curve.Sign(hash, priv.D)
	})
	b.Run(fmt.Sprintf("verify_input_size_%d", hash.BitLen()), func(b *testing.B) {
		curve.Verify(hash, r, s, x, y)
	})
}

func TestKnownSignature(t *testing.T) {
	// good signature
	priv, _ := new(big.Int).SetString("104397037759416840641267745129360920341912682966983343798870479003077644689", 10)
	pubX, _ := new(big.Int).SetString("1913222325711601599563860015182907040361852177892954047964358042507353067365", 10)
	pubY, _ := new(big.Int).SetString("798905265292544287704154888908626830160713383708400542998012716235575472365", 10)
	hash, _ := new(big.Int).SetString("2680576269831035412725132645807649347045997097070150916157159360688041452746", 10)
	rIn, _ := new(big.Int).SetString("607684330780324271206686790958794501662789535258258105407533051445036595885", 10)
	sIn, _ := new(big.Int).SetString("453590782387078613313238308551260565642934039343903827708036287031471258875", 10)

	curve, err := SCWithConstants("./pedersen_params.json")
	if err != nil {
		t.Errorf("Could not init with constant points: %v\n", err)
	}

	if !curve.Verify(hash, rIn, sIn, pubX, pubY) {
		t.Errorf("'known good sig' as actually bad: %v\n", err)
	}

	r, s, err := curve.Sign(hash, priv)
	if err != nil {
		t.Errorf("Could not sign good hash: %v\n", err)
	}

	x, y, err := curve.PrivateToPoint(priv)
	if err != nil {
		t.Errorf("Could not convert random private key to point: %v\n", err)
	}

	if !curve.Verify(hash, r, s, x, y) {
		t.Errorf("Could not verify good signature: %v\ngot: %v %v\n", err, r, s)
	}
}

func TestDerivedSignature(t *testing.T) {
	curve, err := SCWithConstants("./pedersen_params.json")
	if err != nil {
		t.Errorf("Could not init with constant points: %v\n", err)
	}

	pr, _ := curve.GetRandomPrivateKey()

	prin := new(big.Int)
	prin = prin.Set(pr)
	x, y, err := curve.PrivateToPoint(prin)
	if err != nil {
		t.Errorf("Could not convert random private key to point: %v\n", err)
	}

	priv := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: curve,
			X:     x,
			Y:     y,
		},
		D: pr,
	}

	hash, err := curve.PedersenHash([]*big.Int{HexToBN("0x12773"), HexToBN("0x872362")})
	if err != nil {
		t.Errorf("Hashing err: %v\n", err)
	}

	r, s, err := curve.Sign(hash, priv.D)
	if err != nil {
		t.Errorf("Could not convert gen signature: %v\n", err)
	}

	if !curve.Verify(hash, r, s, priv.PublicKey.X, priv.PublicKey.Y) {
		t.Errorf("Could not verify good signature: %v\ngot: %v %v\n", err, r, s)
	}
}

func TestTransactionHash(t *testing.T) {
	curve := SC()

	jtx := JSTransaction{
		Calldata:           []string{"2914367423676101327401096153024331591451054625738519726725779300741401683065", "1284328616562954354594453552152941613439836383012703358554726925609665244667", "3", "1242951120254381876598", "9", "22108152553797646456187940211", "14"},
		ContractAddress:    "0x6f8b21c8354e8ba21ead656932eaa21e728f8c81f001488c186a336d7038cf1",
		EntryPointSelector: "0x240060cdb34fcc260f41eac7474ee1d7c80b7e3607daff9ac67c7ea2ebb1c44",
		EntryPointType:     "EXTERNAL",
		JSSignature:        []string{"1941185432155203218742540925113146991052744726484097092312705586406341211736", "1060098570318028605648271956533461104484177708855341648099672514178101492604"},
		TransactionHash:    "0x14ac93b17d35cc984ff7f186172175cd4341520d32748a406627e48605b38df",
		Nonce:              "0xe",
	}

	tx := jtx.ConvertTx()
	hashFinal, err := curve.HashMsg(
		HexToBN("0x6f8b21c8354e8ba21ead656932eaa21e728f8c81f001488c186a336d7038cf1"),
		tx,
	)
	if err != nil {
		t.Errorf("Could not hash tx arguments: %v\n", err)
	}
	if hashFinal.Cmp(HexToBN("0x2c50e0db592d8149ef09c215846d629206b0d2d40509d313a0b1072f172f0ad")) != 0 {
		t.Errorf("Incorrect hash: got %v expected %v\n", hashFinal, HexToBN("0x2c50e0db592d8149ef09c215846d629206b0d2d40509d313a0b1072f172f0ad"))
	}
}

func TestVerifySignature(t *testing.T) {
	curve := SC()
	hash := HexToBN("0x7f15c38ea577a26f4f553282fcfe4f1feeb8ecfaad8f221ae41abf8224cbddd")
	r, _ := new(big.Int).SetString("2458502865976494910213617956670505342647705497324144349552978333078363662855", 10)
	s, _ := new(big.Int).SetString("3439514492576562277095748549117516048613512930236865921315982886313695689433", 10)

	h, _ := HexToBytes("04033f45f07e1bd1a51b45fc24ec8c8c9908db9e42191be9e169bfcac0c0d997450319d0f53f6ca077c4fa5207819144a2a4165daef6ee47a7c1d06c0dcaa3e456")
	x, y := elliptic.Unmarshal(curve, h)

	if !curve.Verify(hash, r, s, x, y) {
		t.Errorf("successful signature did not verify\n")
	}
}

func TestUIVerifySignature(t *testing.T) {
	curve := SC()
	hash := HexToBN("0x324df642fcc7d98b1d9941250840704f35b9ac2e3e2b58b6a034cc09adac54c")
	r, _ := new(big.Int).SetString("2849277527182985104629156126825776904262411756563556603659114084811678482647", 10)
	s, _ := new(big.Int).SetString("3156340738553451171391693475354397094160428600037567299774561739201502791079", 10)

	pubX, pubY := curve.XToPubKey("0x4e52f2f40700e9cdd0f386c31a1f160d0f310504fc508a1051b747a26070d10")

	if !curve.Verify(hash, r, s, pubX, pubY) {
		t.Errorf("successful signature did not verify\n")
	}
}

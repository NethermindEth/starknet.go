package caigo

import (
	"fmt"
	"math/big"
	"testing"
)

type Mail struct {
	From     Person
	To       Person
	Contents string
}

type Person struct {
	Name   string
	Wallet string
}

func (mail Mail) FmtDefinitionEncoding(field string) (fmtEnc []*big.Int) {
	if field == "from" {
		fmtEnc = append(fmtEnc, UTF8StrToBig(mail.From.Name))
		fmtEnc = append(fmtEnc, HexToBN(mail.From.Wallet))
	} else if field == "to" {
		fmtEnc = append(fmtEnc, UTF8StrToBig(mail.To.Name))
		fmtEnc = append(fmtEnc, HexToBN(mail.To.Wallet))
	} else if field == "contents" {
		fmtEnc = append(fmtEnc, UTF8StrToBig(mail.Contents))
	}
	return fmtEnc
}

func MockTypedData() (ttd TypedData) {
	exampleTypes := make(map[string]TypeDef)
	domDefs := []Definition{{"name", "felt"}, {"version", "felt"}, {"chainId", "felt"}}
	exampleTypes["StarkNetDomain"] = TypeDef{Definitions: domDefs}
	mailDefs := []Definition{{"from", "Person"}, {"to", "Person"}, {"contents", "felt"}}
	exampleTypes["Mail"] = TypeDef{Definitions: mailDefs}
	persDefs := []Definition{{"name", "felt"}, {"wallet", "felt"}}
	exampleTypes["Person"] = TypeDef{Definitions: persDefs}

	dm := Domain{
		Name:    "StarkNet Mail",
		Version: "1",
		ChainId: "1",
	}

	ttd, _ = NewTypedData(exampleTypes, "Mail", dm)
	return ttd
}

func TestGetMessageHash(t *testing.T) {
	ttd := MockTypedData()

	mail := Mail{
		From: Person{
			Name:   "Cow",
			Wallet: "0xCD2a3d9F938E13CD947Ec05AbC7FE734Df8DD826",
		},
		To: Person{
			Name:   "Bob",
			Wallet: "0xbBbBBBBbbBBBbbbBbbBbbbbBBbBbbbbBbBbbBBbB",
		},
		Contents: "Hello, Bob!",
	}

	hash, err := ttd.GetMessageHash(HexToBN("0xCD2a3d9F938E13CD947Ec05AbC7FE734Df8DD826"), mail, Curve)
	if err != nil {
		t.Errorf("Could not hash message: %v\n", err)
	}

	exp := "0x6fcff244f63e38b9d88b9e3378d44757710d1b244282b435cb472053c8d78d0"
	if BigToHex(hash) != exp {
		t.Errorf("type hash: %v does not match expected %v\n", BigToHex(hash), exp)
	}
}

func BenchmarkGetMessageHash(b *testing.B) {
	ttd := MockTypedData()

	mail := Mail{
		From: Person{
			Name:   "Cow",
			Wallet: "0xCD2a3d9F938E13CD947Ec05AbC7FE734Df8DD826",
		},
		To: Person{
			Name:   "Bob",
			Wallet: "0xbBbBBBBbbBBBbbbBbbBbbbbBBbBbbbbBbBbbBBbB",
		},
		Contents: "Hello, Bob!",
	}
	addr := HexToBN("0xCD2a3d9F938E13CD947Ec05AbC7FE734Df8DD826")
	b.Run(fmt.Sprintf("input_size_%d", addr.BitLen()), func(b *testing.B) {
		ttd.GetMessageHash(addr, mail, Curve)
	})
}

func TestGetDomainHash(t *testing.T) {
	ttd := MockTypedData()

	hash, err := ttd.GetTypedMessageHash("StarkNetDomain", ttd.Domain, Curve)
	if err != nil {
		t.Errorf("Could not hash message: %v\n", err)
	}

	exp := "0x54833b121883a3e3aebff48ec08a962f5742e5f7b973469c1f8f4f55d470b07"
	if BigToHex(hash) != exp {
		t.Errorf("type hash: %v does not match expected %v\n", BigToHex(hash), exp)
	}
}

// equivalent of get struct hash
func TestGetTypedMessageHash(t *testing.T) {
	ttd := MockTypedData()

	mail := Mail{
		From: Person{
			Name:   "Cow",
			Wallet: "0xCD2a3d9F938E13CD947Ec05AbC7FE734Df8DD826",
		},
		To: Person{
			Name:   "Bob",
			Wallet: "0xbBbBBBBbbBBBbbbBbbBbbbbBBbBbbbbBbBbbBBbB",
		},
		Contents: "Hello, Bob!",
	}

	hash, err := ttd.GetTypedMessageHash("Mail", mail, Curve)
	if err != nil {
		t.Errorf("Could get typed message hash: %v\n", err)
	}

	exp := "0x4758f1ed5e7503120c228cbcaba626f61514559e9ef5ed653b0b885e0f38aec"
	if BigToHex(hash) != exp {
		t.Errorf("type hash: %v does not match expected %v\n", BigToHex(hash), exp)
	}
}

func TestGetTypeHash(t *testing.T) {
	tdd := MockTypedData()

	hash, err := tdd.GetTypeHash("StarkNetDomain")
	if err != nil {
		t.Errorf("error enccoding type %v\n", err)
	}

	exp := "0x1bfc207425a47a5dfa1a50a4f5241203f50624ca5fdf5e18755765416b8e288"
	if BigToHex(hash) != exp {
		t.Errorf("type hash: %v does not match expected %v\n", BigToHex(hash), exp)
	}

	enc := tdd.Types["StarkNetDomain"]
	if BigToHex(enc.Encoding) != exp {
		t.Errorf("type hash: %v does not match expected %v\n", BigToHex(hash), exp)
	}

	pHash, err := tdd.GetTypeHash("Person")
	if err != nil {
		t.Errorf("error enccoding type %v\n", err)
	}

	exp = "0x2896dbe4b96a67110f454c01e5336edc5bbc3635537efd690f122f4809cc855"
	if BigToHex(pHash) != exp {
		t.Errorf("type hash: %v does not match expected %v\n", BigToHex(pHash), exp)
	}

	enc = tdd.Types["Person"]
	if BigToHex(enc.Encoding) != exp {
		t.Errorf("type hash: %v does not match expected %v\n", BigToHex(hash), exp)
	}
}

func TestGetSelectorFromName(t *testing.T) {
	sel1 := BigToHex(GetSelectorFromName("initialize"))
	sel2 := BigToHex(GetSelectorFromName("mint"))
	sel3 := BigToHex(GetSelectorFromName("test"))

	exp1 := "0x79dc0da7c54b95f10aa182ad0a46400db63156920adb65eca2654c0945a463"
	exp2 := "0x2f0b3c5710379609eb5495f1ecd348cb28167711b73609fe565a72734550354"
	exp3 := "0x22ff5f21f0b81b113e63f7db6da94fedef11b2119b4088b89664fb9a3cb658"

	if sel1 != exp1 || sel2 != exp2 || sel3 != exp3 {
		t.Errorf("invalid Keccak256 encoding: %v %v %v\n", sel1, sel2, sel3)
	}
}

func TestEncodeType(t *testing.T) {
	tdd := MockTypedData()

	enc, err := tdd.EncodeType("Mail")
	if err != nil {
		t.Errorf("error enccoding type %v\n", err)
	}

	exp := "Mail(from:Person,to:Person,contents:felt)Person(name:felt,wallet:felt)"
	if enc != exp {
		t.Errorf("type encoding: %v does not match expected %v\n", enc, exp)
	}
}

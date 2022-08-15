package felt

import (
	"fmt"
	"math/big"
	"testing"
)

func BenchmarkPedersenHash(b *testing.B) {
	suite := [][]*big.Int{
		{strToBig("0x12773"), strToBig("0x872362")},
		{strToBig("0x1277312773"), strToBig("0x872362872362")},
		{strToBig("0x1277312773"), strToBig("0xCD2a3d9F938E13CD947Ec05AbC7FE734Df8DD826")},
		{strToBig("0xbBbBBBBbbBBBbbbBbbBbbbbBBbBbbbbBbBbbBBbB"), strToBig("0x872362872362")},
		{strToBig("0xCD2a3d9F938E13CD947Ec05AbC7FE734Df8DD826"), strToBig("0xbBbBBBBbbBBBbbbBbbBbbbbBBbBbbbbBbBbbBBbB")},
		{strToBig("0x7f15c38ea577a26f4f553282fcfe4f1feeb8ecfaad8f221ae41abf8224cbddd"), strToBig("0x13d41f388b8ea4db56c5aa6562f13359fab192b3db57651af916790f9debee9")},
		{strToBig("0x7f15c38ea577a26f4f553282fcfe4f1feeb8ecfaad8f221ae41abf8224cbddd"), strToBig("0x7f15c38ea577a26f4f553282fcfe4f1feeb8ecfaad8f221ae41abf8224cbdde")},
	}

	for _, test := range suite {
		b.Run(fmt.Sprintf("input_size_%d_%d", test[0].BitLen(), test[1].BitLen()), func(b *testing.B) {
			Curve.PedersenHash(test)
		})
	}
}

func TestPedersenHash(t *testing.T) {
	testPedersen := []struct {
		elements []*big.Int
		expected *big.Int
	}{
		{
			elements: []*big.Int{strToBig("0x12773"), strToBig("0x872362")},
			expected: strToBig("0x5ed2703dfdb505c587700ce2ebfcab5b3515cd7e6114817e6026ec9d4b364ca"),
		},
		{
			elements: []*big.Int{strToBig("0x13d41f388b8ea4db56c5aa6562f13359fab192b3db57651af916790f9debee9"), strToBig("0x537461726b4e6574204d61696c")},
			expected: strToBig("0x180c0a3d13c1adfaa5cbc251f4fc93cc0e26cec30ca4c247305a7ce50ac807c"),
		},
		{
			elements: []*big.Int{big.NewInt(100), big.NewInt(1000)},
			expected: strToBig("0x45a62091df6da02dce4250cb67597444d1f465319908486b836f48d0f8bf6e7"),
		},
	}

	for _, tt := range testPedersen {
		hash, err := Curve.PedersenHash(tt.elements)
		if err != nil {
			t.Errorf("Hashing err: %v\n", err)
		}
		if hash.Cmp(tt.expected) != 0 {
			t.Errorf("incorrect hash: got %v expected %v\n", hash, tt.expected)
		}
	}
}

func TestDivMod(t *testing.T) {
	testDivmod := []struct {
		x        *big.Int
		y        *big.Int
		expected *big.Int
	}{
		{
			x:        strToBig("311379432064974854430469844112069886938521247361583891764940938105250923060"),
			y:        strToBig("621253665351494585790174448601059271924288186997865022894315848222045687999"),
			expected: strToBig("2577265149861519081806762825827825639379641276854712526969977081060187505740"),
		},
		{
			x:        big.NewInt(1),
			y:        big.NewInt(2),
			expected: strToBig("0x0400000000000008800000000000000000000000000000000000000000000001"),
		},
	}

	for _, tt := range testDivmod {
		divR := DivMod(tt.x, tt.y, Curve.P)

		if divR.Cmp(tt.expected) != 0 {
			t.Errorf("DivMod Res %v does not == expected %v\n", divR, tt.expected)
		}
	}
}

func TestAdd(t *testing.T) {
	testAdd := []struct {
		x         *big.Int
		y         *big.Int
		expectedX *big.Int
		expectedY *big.Int
	}{
		{
			x:         strToBig("1468732614996758835380505372879805860898778283940581072611506469031548393285"),
			y:         strToBig("1402551897475685522592936265087340527872184619899218186422141407423956771926"),
			expectedX: strToBig("2573054162739002771275146649287762003525422629677678278801887452213127777391"),
			expectedY: strToBig("3086444303034188041185211625370405120551769541291810669307042006593736192813"),
		},
		{
			x:         big.NewInt(1),
			y:         big.NewInt(2),
			expectedX: strToBig("225199957243206662471193729647752088571005624230831233470296838210993906468"),
			expectedY: strToBig("190092378222341939862849656213289777723812734888226565973306202593691957981"),
		},
	}

	for _, tt := range testAdd {
		resX, resY := Curve.Add(Curve.Gx, Curve.Gy, tt.x, tt.y)
		if resX.Cmp(tt.expectedX) != 0 {
			t.Errorf("ResX %v does not == expected %v\n", resX, tt.expectedX)

		}
		if resY.Cmp(tt.expectedY) != 0 {
			t.Errorf("ResY %v does not == expected %v\n", resY, tt.expectedY)
		}
	}
}

func TestMultAir(t *testing.T) {
	testMult := []struct {
		r         *big.Int
		x         *big.Int
		y         *big.Int
		expectedX *big.Int
		expectedY *big.Int
	}{
		{
			r:         strToBig("2458502865976494910213617956670505342647705497324144349552978333078363662855"),
			x:         strToBig("1468732614996758835380505372879805860898778283940581072611506469031548393285"),
			y:         strToBig("1402551897475685522592936265087340527872184619899218186422141407423956771926"),
			expectedX: strToBig("182543067952221301675635959482860590467161609552169396182763685292434699999"),
			expectedY: strToBig("3154881600662997558972388646773898448430820936643060392452233533274798056266"),
		},
	}

	for _, tt := range testMult {
		x, y, err := Curve.MimicEcMultAir(tt.r, tt.x, tt.y, Curve.Gx, Curve.Gy)
		if err != nil {
			t.Errorf("MultAirERR %v\n", err)
		}

		if x.Cmp(tt.expectedX) != 0 {
			t.Errorf("ResX %v does not == expected %v\n", x, tt.expectedX)

		}
		if y.Cmp(tt.expectedY) != 0 {
			t.Errorf("ResY %v does not == expected %v\n", y, tt.expectedY)
		}
	}
}

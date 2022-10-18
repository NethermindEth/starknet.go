package caigo

import (
	"crypto/subtle"
	"fmt"
	"math"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/dontpanicdao/caigo/types"
	"gonum.org/v1/gonum/stat"
)

func BenchmarkPedersenHash(b *testing.B) {
	suite := [][]*big.Int{
		{types.HexToBN("0x12773"), types.HexToBN("0x872362")},
		{types.HexToBN("0x1277312773"), types.HexToBN("0x872362872362")},
		{types.HexToBN("0x1277312773"), types.HexToBN("0xCD2a3d9F938E13CD947Ec05AbC7FE734Df8DD826")},
		{types.HexToBN("0xbBbBBBBbbBBBbbbBbbBbbbbBBbBbbbbBbBbbBBbB"), types.HexToBN("0x872362872362")},
		{types.HexToBN("0xCD2a3d9F938E13CD947Ec05AbC7FE734Df8DD826"), types.HexToBN("0xbBbBBBBbbBBBbbbBbbBbbbbBBbBbbbbBbBbbBBbB")},
		{types.HexToBN("0x7f15c38ea577a26f4f553282fcfe4f1feeb8ecfaad8f221ae41abf8224cbddd"), types.HexToBN("0x13d41f388b8ea4db56c5aa6562f13359fab192b3db57651af916790f9debee9")},
		{types.HexToBN("0x7f15c38ea577a26f4f553282fcfe4f1feeb8ecfaad8f221ae41abf8224cbddd"), types.HexToBN("0x7f15c38ea577a26f4f553282fcfe4f1feeb8ecfaad8f221ae41abf8224cbdde")},
	}

	for _, test := range suite {
		b.Run(fmt.Sprintf("input_size_%d_%d", test[0].BitLen(), test[1].BitLen()), func(b *testing.B) {
			Curve.PedersenHash(test)
		})
	}
}

func BenchmarkCurveSign(b *testing.B) {
	type data struct {
		MessageHash *big.Int
		PrivateKey  *big.Int
		Seed        *big.Int
	}

	dataSet := []data{}
	MessageHash := big.NewInt(0).Exp(big.NewInt(2), big.NewInt(250), nil)
	PrivateKey := big.NewInt(0).Add(MessageHash, big.NewInt(1))
	Seed := big.NewInt(0)
	for i := int64(0); i < 20; i++ {
		dataSet = append(dataSet, data{
			MessageHash: big.NewInt(0).Add(MessageHash, big.NewInt(i)),
			PrivateKey:  big.NewInt(0).Add(PrivateKey, big.NewInt(i)),
			Seed:        big.NewInt(0).Add(Seed, big.NewInt(i)),
		})

		for _, test := range dataSet {
			Curve.Sign(test.MessageHash, test.PrivateKey, test.Seed)
		}
	}
}

func TestPedersenHash(t *testing.T) {
	testPedersen := []struct {
		elements []*big.Int
		expected *big.Int
	}{
		{
			elements: []*big.Int{types.HexToBN("0x12773"), types.HexToBN("0x872362")},
			expected: types.HexToBN("0x5ed2703dfdb505c587700ce2ebfcab5b3515cd7e6114817e6026ec9d4b364ca"),
		},
		{
			elements: []*big.Int{types.HexToBN("0x13d41f388b8ea4db56c5aa6562f13359fab192b3db57651af916790f9debee9"), types.HexToBN("0x537461726b4e6574204d61696c")},
			expected: types.HexToBN("0x180c0a3d13c1adfaa5cbc251f4fc93cc0e26cec30ca4c247305a7ce50ac807c"),
		},
		{
			elements: []*big.Int{big.NewInt(100), big.NewInt(1000)},
			expected: types.HexToBN("0x45a62091df6da02dce4250cb67597444d1f465319908486b836f48d0f8bf6e7"),
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
			x:        types.StrToBig("311379432064974854430469844112069886938521247361583891764940938105250923060"),
			y:        types.StrToBig("621253665351494585790174448601059271924288186997865022894315848222045687999"),
			expected: types.StrToBig("2577265149861519081806762825827825639379641276854712526969977081060187505740"),
		},
		{
			x:        big.NewInt(1),
			y:        big.NewInt(2),
			expected: types.HexToBN("0x0400000000000008800000000000000000000000000000000000000000000001"),
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
			x:         types.StrToBig("1468732614996758835380505372879805860898778283940581072611506469031548393285"),
			y:         types.StrToBig("1402551897475685522592936265087340527872184619899218186422141407423956771926"),
			expectedX: types.StrToBig("2573054162739002771275146649287762003525422629677678278801887452213127777391"),
			expectedY: types.StrToBig("3086444303034188041185211625370405120551769541291810669307042006593736192813"),
		},
		{
			x:         big.NewInt(1),
			y:         big.NewInt(2),
			expectedX: types.StrToBig("225199957243206662471193729647752088571005624230831233470296838210993906468"),
			expectedY: types.StrToBig("190092378222341939862849656213289777723812734888226565973306202593691957981"),
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
	tests := []struct {
		r         *big.Int
		x         *big.Int
		y         *big.Int
		expectedX *big.Int
		expectedY *big.Int
	}{
		{
			r:         types.StrToBig("2458502865976494910213617956670505342647705497324144349552978333078363662855"),
			x:         types.StrToBig("1468732614996758835380505372879805860898778283940581072611506469031548393285"),
			y:         types.StrToBig("1402551897475685522592936265087340527872184619899218186422141407423956771926"),
			expectedX: types.StrToBig("182543067952221301675635959482860590467161609552169396182763685292434699999"),
			expectedY: types.StrToBig("3154881600662997558972388646773898448430820936643060392452233533274798056266"),
		},
	}

	for _, tt := range tests {
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

// swappable ec multiplication fn
type ecMultiFn func(m, x1, y1 *big.Int) (x, y *big.Int)
type ecMultOption struct {
	algo   string
	fn     ecMultiFn
	stddev float64
}

// Get multiple ec multiplication algo options to test and benchmark
func (sc StarkCurve) ecMultOptions() []ecMultOption {
	return []ecMultOption{
		{
			algo: "Double-And-Add",
			fn:   sc.ecMult_DoubleAndAdd, // original algo
		},
		{
			algo: "Double-And-Always-Add",
			fn:   sc.EcMult, // best algo (currently used)
		},
		{
			algo: "Montgomery-Ladder",
			fn:   sc.ecMult_Montgomery,
		},
		{
			algo: "Montgomery-Ladder-Lsh",
			fn:   sc.ecMult_MontgomeryLsh,
		},
	}
}

func FuzzEcMult(f *testing.F) {
	// Generate the scalar value k, where 0 < k < order(point)
	var _genScalar = func(a int, b int) (k *big.Int) {
		k = new(big.Int).Mul(big.NewInt(int64(a)), big.NewInt(int64(b)))
		k = k.Mul(k, k).Mul(k, k) // generate moar big number
		k = k.Abs(k)
		k = k.Add(k, big.NewInt(1)) // edge case: avoid zero
		k = k.Mod(k, Curve.N)
		return
	}

	// Seed the fuzzer (examples)
	f.Add(-12121501143923232, 142312310232324552) // negative numbers used as seeds but the resulting
	f.Add(41289371293219038, -179566705053432322) // scalar is normalized to 0 < k < order(point)
	f.Add(927302501143912223, 220390912389202149)
	f.Add(874739451078007766, 868575557812948233)
	f.Add(302150520188025637, 670505342647705232)
	f.Add(778320444456588442, 932884823101831273)
	f.Add(658844239552133924, 933442778319932884)
	f.Add(494910213617956623, 976290247577832044)

	f.Fuzz(func(t *testing.T, a int, b int) {
		k := _genScalar(a, b)

		var x0, y0 *big.Int
		for _, tt := range Curve.ecMultOptions() {
			x, y, err := Curve.privateToPoint(k, tt.fn)
			if err != nil {
				t.Errorf("EcMult err: %v, algo=%v\n", err, tt.algo)
			}

			// Store the initial result from the first algo and test against it
			if x0 == nil {
				x0 = x
				y0 = y
			} else if x0.Cmp(x) != 0 {
				t.Errorf("EcMult x mismatch: %v != %v, algo=%v\n", x, x0, tt.algo)
			} else if y0.Cmp(y) != 0 {
				t.Errorf("EcMult y mismatch: %v != %v, algo=%v\n", y, y0, tt.algo)
			}
		}
	})
}

func BenchmarkEcMultAll(b *testing.B) {
	// Generate the scalar value k, where n number of bits are set, no trailing zeros
	var _genScalarBits = func(n int) (k *big.Int) {
		k = big.NewInt(1)
		for i := 1; i < n; i++ {
			k = k.Lsh(k, 1).Add(k, big.NewInt(1))
		}
		return
	}

	ecMultiBest := ecMultOption{
		algo:   "",
		stddev: math.MaxFloat64,
	}

	var out strings.Builder
	for _, tt := range Curve.ecMultOptions() {
		// test (+ time) injected ec multi fn performance via Curve.privateToPoint
		var _test = func(k *big.Int) int64 {
			start := time.Now()
			Curve.privateToPoint(k, tt.fn)
			return time.Since(start).Nanoseconds()
		}

		xs := []float64{}
		// generate numbers with 1 to 251 bits set
		for i := 1; i < Curve.N.BitLen(); i++ {
			k := _genScalarBits(i)
			b.Run(fmt.Sprintf("%s/input_bits_len/%d", tt.algo, k.BitLen()), func(b *testing.B) {
				ns := _test(k)
				xs = append(xs, float64(ns))
			})
		}

		// generate numbers with 1 to 250 trailing zero bits set
		k := _genScalarBits(Curve.N.BitLen() - 1)
		for i := 1; i < Curve.N.BitLen()-1; i++ {
			k.Rsh(k, uint(i)).Lsh(k, uint(i))
			b.Run(fmt.Sprintf("%s/input_bits_len/%d#%d", tt.algo, k.BitLen(), k.TrailingZeroBits()), func(b *testing.B) {
				ns := _test(k)
				xs = append(xs, float64(ns))
			})
		}

		// computes the weighted mean of the dataset.
		// we don't have any weights (ie: all weights are 1) so we pass a nil slice.
		mean := stat.Mean(xs, nil)
		variance := stat.Variance(xs, nil)
		stddev := math.Sqrt(variance)
		// Keep track of the best one (min stddev)
		if stddev < ecMultiBest.stddev {
			ecMultiBest.stddev = stddev
			ecMultiBest.algo = tt.algo
		}

		out.WriteString("-----------------------------\n")
		out.WriteString(fmt.Sprintf("algo=       %v\n", tt.algo))
		out.WriteString(fmt.Sprintf("stats(ns)\n"))
		out.WriteString(fmt.Sprintf("  mean=     %v\n", mean))
		out.WriteString(fmt.Sprintf("  variance= %v\n", variance))
		out.WriteString(fmt.Sprintf("  std-dev=  %v\n", stddev))
		out.WriteString("\n")
	}

	// final stats output
	fmt.Println(out.String())
	// assert benchmark result is as expected
	expectedBest := "Double-And-Always-Add"
	if ecMultiBest.algo != expectedBest {
		b.Errorf("ecMultiBest.algo %v does not == expected %v\n", ecMultiBest.algo, expectedBest)
	}
}

// Multiplies by m a point on the elliptic curve with equation y^2 = x^3 + alpha*x + beta mod p.
// Assumes affine form (x, y) is spread (x1 *big.Int, y1 *big.Int) and that 0 < m < order(point).
//
// (ref: https://github.com/starkware-libs/cairo-lang/blob/master/src/starkware/crypto/starkware/crypto/signature/math_utils.py)
func (sc StarkCurve) ecMult_DoubleAndAdd(m, x1, y1 *big.Int) (x, y *big.Int) {
	var _ecMult func(m, x1, y1 *big.Int) (x, y *big.Int)
	_ecMult = func(m, x1, y1 *big.Int) (x, y *big.Int) {
		if m.BitLen() == 1 {
			return x1, y1
		}
		mk := new(big.Int).Mod(m, big.NewInt(2))
		if mk.Cmp(big.NewInt(0)) == 0 {
			h := new(big.Int).Div(m, big.NewInt(2))
			c, d := sc.Double(x1, y1)
			return _ecMult(h, c, d)
		}
		n := new(big.Int).Sub(m, big.NewInt(1))
		e, f := _ecMult(n, x1, y1)

		return sc.Add(e, f, x1, y1)
	}

	// Notice: no need for scalar rewrite trick via `StarkCurve.rewriteScalar`
	//   This algorithm is not affected, as it doesn't do a fixed number of operations,
	//   nor directly depends on the binary representation of the scalar.
	return _ecMult(m, x1, y1)
}

// Multiplies by m a point on the elliptic curve with equation y^2 = x^3 + alpha*x + beta mod p.
// Assumes affine form (x, y) is spread (x1 *big.Int, y1 *big.Int) and that 0 < m < order(point).
//
// (ref: https://en.wikipedia.org/wiki/Elliptic_curve_point_multiplication#Montgomery_ladder)
func (sc StarkCurve) ecMult_Montgomery(m, x1, y1 *big.Int) (x, y *big.Int) {
	var _ecMultMontgomery = func(m, x0, y0, x1, y1 *big.Int) (x, y *big.Int) {
		// Do constant number of operations
		for i := sc.N.BitLen() - 1; i >= 0; i-- {
			// Check if next bit set
			if m.Bit(i) == 0 {
				x1, y1 = sc.Add(x0, y0, x1, y1)
				x0, y0 = sc.Double(x0, y0)
			} else {
				x0, y0 = sc.Add(x0, y0, x1, y1)
				x1, y1 = sc.Double(x1, y1)
			}
		}
		return x0, y0
	}

	return _ecMultMontgomery(sc.rewriteScalar(m), big.NewInt(0), big.NewInt(0), x1, y1)
}

// Multiplies by m a point on the elliptic curve with equation y^2 = x^3 + alpha*x + beta mod p.
// Assumes affine form (x, y) is spread (x1 *big.Int, y1 *big.Int) and that 0 < m < order(point).
//
// (ref: https://en.wikipedia.org/wiki/Elliptic_curve_point_multiplication#Montgomery_ladder)
func (sc StarkCurve) ecMult_MontgomeryLsh(m, x1, y1 *big.Int) (x, y *big.Int) {
	var _ecMultMontgomery = func(m, x0, y0, x1, y1 *big.Int) (x, y *big.Int) {
		// Fill a fixed 32 byte buffer (2 ** 251)
		// NOTICE: this will take an absolute value first
		buf := m.FillBytes(make([]byte, 32))

		for i, byte := range buf {
			for bitNum := 0; bitNum < 8; bitNum++ {
				// Skip first 4 bits, do constant 252 operations
				if i == 0 && bitNum < 4 {
					byte <<= 1
					continue
				}

				// Check if next bit set
				if subtle.ConstantTimeByteEq(byte&0x80, 0x80) == 0 {
					x1, y1 = sc.Add(x0, y0, x1, y1)
					x0, y0 = sc.Double(x0, y0)
				} else {
					x0, y0 = sc.Add(x0, y0, x1, y1)
					x1, y1 = sc.Double(x1, y1)
				}
				byte <<= 1
			}
		}
		return x0, y0
	}

	return _ecMultMontgomery(sc.rewriteScalar(m), big.NewInt(0), big.NewInt(0), x1, y1)
}

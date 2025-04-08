package curve

import (
	"crypto/elliptic"
	"crypto/subtle"
	"fmt"
	"math"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/NethermindEth/juno/core/felt"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/stretchr/testify/require"
	"gonum.org/v1/gonum/stat"
)

// package level variable to be used by the benchmarking code
// to prevent the compiler from optimizing the code away
var result any

// BenchmarkPedersenHash benchmarks the performance of the PedersenHash function.
//
// The function takes a 2D slice of big.Int values as input and measures the time
// it takes to execute the PedersenHash function for each test case.
//
// Parameters:
// - b: a *testing.B value representing the testing context
// Returns:
//
//	none
func BenchmarkPedersenHash(b *testing.B) {
	suite := [][]string{
		{"0x12773", "0x872362"},
		{"0x1277312773", "0x872362872362"},
		{"0x1277312773", "0xCD2a3d9F938E13CD947Ec05AbC7FE734Df8DD826"},
		{"0xbBbBBBBbbBBBbbbBbbBbbbbBBbBbbbbBbBbbBBbB", "0x872362872362"},
		{"0xCD2a3d9F938E13CD947Ec05AbC7FE734Df8DD826", "0xbBbBBBBbbBBBbbbBbbBbbbbBBbBbbbbBbBbbBBbB"},
		{"0x7f15c38ea577a26f4f553282fcfe4f1feeb8ecfaad8f221ae41abf8224cbddd", "0x13d41f388b8ea4db56c5aa6562f13359fab192b3db57651af916790f9debee9"},
		{"0x7f15c38ea577a26f4f553282fcfe4f1feeb8ecfaad8f221ae41abf8224cbddd", "0x7f15c38ea577a26f4f553282fcfe4f1feeb8ecfaad8f221ae41abf8224cbdde"},
	}

	for _, test := range suite {
		b.Run(fmt.Sprintf("input_size_%d_%d", len(test[0]), len(test[1])), func(b *testing.B) {
			hexArr, err := internalUtils.HexArrToFelt(test)
			require.NoError(b, err)
			result = Pedersen(hexArr[0], hexArr[1])
		})
	}
}

// BenchmarkCurveSign benchmarks the Curve.Sign function.
//
// Parameters:
// - b: a *testing.B value representing the testing context
// Returns:
//
//	none
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
			result, _, err := Curve.Sign(test.MessageHash, test.PrivateKey, test.Seed)
			require.NoError(b, err)
			require.NotEmpty(b, result)
		}
	}
}

// BenchmarkSignatureVerify benchmarks the SignatureVerify function.
//
// The function takes a testing.B object as a parameter and performs a series
// of operations to benchmark the SignatureVerify function. It generates a
// random private key, computes the corresponding public key, computes a hash
// using the PedersenHash function, signs the hash using the private key, and
// finally verifies the signature. The function runs two benchmarks: one for
// signing and one for verification. Each benchmark measures the time taken to
// perform the respective operation.
//
// Parameters:
// - b: a *testing.B value representing the testing context
// Returns:
//
//	none
func BenchmarkSignatureVerify(b *testing.B) {
	private, err := Curve.GetRandomPrivateKey()
	require.NoError(b, err)
	x, y, err := Curve.PrivateToPoint(private)
	require.NoError(b, err)

	hash := Pedersen(
		internalUtils.TestHexToFelt(b, "0x7f15c38ea577a26f4f553282fcfe4f1feeb8ecfaad8f221ae41abf8224cbddd"),
		internalUtils.TestHexToFelt(b, "0x7f15c38ea577a26f4f553282fcfe4f1feeb8ecfaad8f221ae41abf8224cbdde"),
	)
	hashBigInt := internalUtils.FeltToBigInt(hash)

	r, s, err := Curve.Sign(hashBigInt, private)
	require.NoError(b, err)

	b.Run(fmt.Sprintf("sign_input_size_%d", hashBigInt.BitLen()), func(b *testing.B) {
		result, _, err = Curve.Sign(hashBigInt, private)
		require.NoError(b, err)
		require.NotEmpty(b, result)
	})
	b.Run(fmt.Sprintf("verify_input_size_%d", hashBigInt.BitLen()), func(b *testing.B) {
		result = Curve.Verify(hashBigInt, r, s, x, y)
		require.NotEmpty(b, result)
	})
}

// TestGeneral_PrivateToPoint tests the PrivateToPoint function.
//
// Parameters:
// - t: a *testing.T value representing the testing context
// Returns:
//
//	none
func TestGeneral_PrivateToPoint(t *testing.T) {
	x, _, err := Curve.PrivateToPoint(big.NewInt(2))
	require.NoError(t, err)
	expectedX, ok := new(big.Int).SetString("3324833730090626974525872402899302150520188025637965566623476530814354734325", 10)
	require.True(t, ok)
	require.Equal(t, expectedX, x)
}

// TestGeneral_PedersenHash is a test function for the PedersenHash method in the General struct.
//
// The function tests the PedersenHash method by providing different test cases and comparing the computed hash with the expected hash.
// It uses the testing.T type from the testing package to report any errors encountered during the tests.
//
// Parameters:
// - t: a *testing.T value representing the testing context
// Returns:
//
//	none
func TestGeneral_PedersenHash(t *testing.T) {
	testPedersen := []struct {
		elements []string
		expected string
	}{
		{
			elements: []string{"0x12773", "0x872362"},
			expected: "0x5ed2703dfdb505c587700ce2ebfcab5b3515cd7e6114817e6026ec9d4b364ca",
		},
		{
			elements: []string{"0x13d41f388b8ea4db56c5aa6562f13359fab192b3db57651af916790f9debee9", "0x537461726b4e6574204d61696c"},
			expected: "0x180c0a3d13c1adfaa5cbc251f4fc93cc0e26cec30ca4c247305a7ce50ac807c",
		},
		{
			elements: []string{"100", "1000"},
			expected: "0x45a62091df6da02dce4250cb67597444d1f465319908486b836f48d0f8bf6e7",
		},
	}

	for _, test := range testPedersen {
		elementsFelt, err := internalUtils.HexArrToFelt(test.elements)
		require.NoError(t, err)
		expected := internalUtils.TestHexToFelt(t, test.expected)

		result := Pedersen(elementsFelt[0], elementsFelt[1])
		require.Equal(t, expected, result)
	}
}

// TestGeneral_DivMod tests the DivMod function.
//
// The function takes in a list of test cases which consist of inputs x, y, and the expected output.
// The inputs x and y are of type *big.Int and the expected output is also of type *big.Int.
// The function iterates through each test case and calls the DivMod function with the inputs x, y, and the prime number Curve.P.
// It then compares the output of DivMod with the expected output and throws an error if they are not equal.
// The function is used to test the correctness of the DivMod function.
//
// Parameters:
// - t: a *testing.T value representing the testing context
// Returns:
//
//	none
func TestGeneral_DivMod(t *testing.T) {
	testDivmod := []struct {
		x        *big.Int
		y        *big.Int
		expected *big.Int
	}{
		{
			x:        internalUtils.StrToBig("311379432064974854430469844112069886938521247361583891764940938105250923060"),
			y:        internalUtils.StrToBig("621253665351494585790174448601059271924288186997865022894315848222045687999"),
			expected: internalUtils.StrToBig("2577265149861519081806762825827825639379641276854712526969977081060187505740"),
		},
		{
			x:        big.NewInt(1),
			y:        big.NewInt(2),
			expected: internalUtils.HexToBN("0x0400000000000008800000000000000000000000000000000000000000000001"),
		},
	}

	for _, tt := range testDivmod {
		divR := DivMod(tt.x, tt.y, Curve.P)

		require.Equal(t, tt.expected, divR)
	}
}

// TestGeneral_Add tests the Add function in the General package.
//
// It tests the addition of two big integers and compares the result with the expected values.
// The function takes a slice of test cases, each containing two big integers and their expected sum.
// It iterates over the test cases, computes the sum using the Add function, and checks if it matches the expected sum.
// If the computed sum does not match the expected sum, an error is reported using the require.Equal function.
//
// Parameters:
// - t: a *testing.T value representing the testing context
// Returns:
//
//	none
func TestGeneral_Add(t *testing.T) {
	testAdd := []struct {
		x         *big.Int
		y         *big.Int
		expectedX *big.Int
		expectedY *big.Int
	}{
		{
			x:         internalUtils.StrToBig("1468732614996758835380505372879805860898778283940581072611506469031548393285"),
			y:         internalUtils.StrToBig("1402551897475685522592936265087340527872184619899218186422141407423956771926"),
			expectedX: internalUtils.StrToBig("2573054162739002771275146649287762003525422629677678278801887452213127777391"),
			expectedY: internalUtils.StrToBig("3086444303034188041185211625370405120551769541291810669307042006593736192813"),
		},
		{
			x:         big.NewInt(1),
			y:         big.NewInt(2),
			expectedX: internalUtils.StrToBig("225199957243206662471193729647752088571005624230831233470296838210993906468"),
			expectedY: internalUtils.StrToBig("190092378222341939862849656213289777723812734888226565973306202593691957981"),
		},
	}

	for _, tt := range testAdd {
		resX, resY := Curve.Add(Curve.Gx, Curve.Gy, tt.x, tt.y)
		require.Equal(t, tt.expectedX, resX)
		require.Equal(t, tt.expectedY, resY)
	}
}

// TestGeneral_MultAir is a test function that tests the MultAir function in the General package.
//
// It tests the behavior of the MultAir function with different input values and checks if the
// returned values match the expected values.
// The MultAir function takes a big.Int value 'r', and two other big.Int values 'x' and 'y', and
// performs a mathematical operation on them to calculate new big.Int values 'x' and 'y'.
// The function then compares the calculated values with the expected values and reports any
// mismatches as test failures.
//
// Parameters:
// - t: a *testing.T value representing the testing context
// Returns:
//
//	none
func TestGeneral_MultAir(t *testing.T) {
	testMult := []struct {
		r         *big.Int
		x         *big.Int
		y         *big.Int
		expectedX *big.Int
		expectedY *big.Int
	}{
		{
			r:         internalUtils.StrToBig("2458502865976494910213617956670505342647705497324144349552978333078363662855"),
			x:         internalUtils.StrToBig("1468732614996758835380505372879805860898778283940581072611506469031548393285"),
			y:         internalUtils.StrToBig("1402551897475685522592936265087340527872184619899218186422141407423956771926"),
			expectedX: internalUtils.StrToBig("182543067952221301675635959482860590467161609552169396182763685292434699999"),
			expectedY: internalUtils.StrToBig("3154881600662997558972388646773898448430820936643060392452233533274798056266"),
		},
	}

	for _, tt := range testMult {
		x, y, err := Curve.MimicEcMultAir(tt.r, tt.x, tt.y, Curve.Gx, Curve.Gy)
		require.NoError(t, err)
		require.Equal(t, tt.expectedX, x)
		require.Equal(t, tt.expectedY, y)
	}
}

// TestGeneral_ComputeHashOnElements is a test function that verifies the correctness of the ComputeHashOnElements and PedersenArray functions in the General package.
//
// This function tests both functions by passing in different arrays of big.Int elements and comparing the computed hash with the expected hash.
// It checks the behavior of the functions when an empty array is passed as input, as well as when an array with multiple elements is passed.
//
// Parameters:
// - t: a *testing.T value representing the testing context
// Returns:
//
//	none
func TestGeneral_ComputeHashOnElements(t *testing.T) {
	hashEmptyArray := ComputeHashOnElements([]*big.Int{})
	hashEmptyArrayFelt := PedersenArray([]*felt.Felt{}...)

	expectedHashEmmptyArray := internalUtils.HexToBN("0x49ee3eba8c1600700ee1b87eb599f16716b0b1022947733551fde4050ca6804")
	require.Equal(t, hashEmptyArray, expectedHashEmmptyArray, "Hash empty array wrong value.")
	require.Equal(t, internalUtils.FeltToBigInt(hashEmptyArrayFelt), expectedHashEmmptyArray, "Hash empty array wrong value.")

	filledArray := []*big.Int{
		big.NewInt(123782376),
		big.NewInt(213984),
		big.NewInt(128763521321),
	}

	hashFilledArray := ComputeHashOnElements(filledArray)
	hashFilledArrayFelt := PedersenArray(internalUtils.BigIntArrToFeltArr(filledArray)...)

	expectedHashFilledArray := internalUtils.HexToBN("0x7b422405da6571242dfc245a43de3b0fe695e7021c148b918cd9cdb462cac59")
	require.Equal(t, hashFilledArray, expectedHashFilledArray, "Hash filled array wrong value.")
	require.Equal(t, internalUtils.FeltToBigInt(hashFilledArrayFelt), expectedHashFilledArray, "Hash filled array wrong value.")
}

// TestGeneral_HashAndSign is a test function that verifies the hashing and signing process.
//
// Parameters:
// - t: The testing.T object for running the test.
// Returns:
//
//	none
func TestGeneral_HashAndSign(t *testing.T) {
	hashy := HashPedersenElements([]*big.Int{
		big.NewInt(1953658213),
		big.NewInt(126947999705460),
		big.NewInt(1953658213),
	})

	priv, err := Curve.GetRandomPrivateKey()
	require.NoError(t, err)
	x, y, err := Curve.PrivateToPoint(priv)
	require.NoError(t, err)

	r, s, err := Curve.Sign(hashy, priv)
	require.NoError(t, err)

	require.True(t, Curve.Verify(hashy, r, s, x, y))
}

// TestGeneral_ComputeFact tests the ComputeFact function.
//
// It tests the ComputeFact function by providing a set of test cases
// and comparing the computed hash with the expected hash.
// The test cases consist of program hashes, program outputs,
// and expected hash values.
//
// Parameters:
// - t: The testing.T object for running the test
// Returns:
//
//	none
func TestGeneral_ComputeFact(t *testing.T) {
	testFacts := []struct {
		programHash   *big.Int
		programOutput []*big.Int
		expected      *big.Int
	}{
		{
			programHash:   internalUtils.HexToBN("0x114952172aed91e59f870a314e75de0a437ff550e4618068cec2d832e48b0c7"),
			programOutput: []*big.Int{big.NewInt(289)},
			expected:      internalUtils.HexToBN("0xe6168c0a865aa80d724ad05627fa65fbcfe4b1d66a586e9f348f461b076072c4"),
		},
		{
			programHash:   internalUtils.HexToBN("0x79920d895101ad1fbdea9adf141d8f362fdea9ee35f33dfcd07f38e4a589bab"),
			programOutput: []*big.Int{internalUtils.StrToBig("2754806153357301156380357983574496185342034785016738734224771556919270737441")},
			expected:      internalUtils.HexToBN("0x1d174fa1443deea9aab54bbca8d9be308dd14a0323dd827556c173bd132098db"),
		},
	}

	for _, tt := range testFacts {
		hash := internalUtils.ComputeFact(tt.programHash, tt.programOutput)
		require.Equal(t, tt.expected, hash)
	}
}

// TestGeneral_BadSignature tests the behavior of the function that checks for bad signatures.
//
// Parameters:
// - t: The testing.T object for running the test
// Returns:
//
//	none
func TestGeneral_BadSignature(t *testing.T) {
	hash := Pedersen(internalUtils.TestHexToFelt(t, "0x12773"), internalUtils.TestHexToFelt(t, "0x872362"))
	hashBigInt := internalUtils.FeltToBigInt(hash)

	priv, err := Curve.GetRandomPrivateKey()
	require.NoError(t, err)
	x, y, err := Curve.PrivateToPoint(priv)
	require.NoError(t, err)

	r, s, err := Curve.Sign(hashBigInt, priv)
	require.NoError(t, err)

	badR := new(big.Int).Add(r, big.NewInt(1))
	require.False(t, Curve.Verify(hashBigInt, badR, s, x, y))

	badS := new(big.Int).Add(s, big.NewInt(1))
	require.False(t, Curve.Verify(hashBigInt, r, badS, x, y))

	badHash := new(big.Int).Add(hashBigInt, big.NewInt(1))
	require.False(t, Curve.Verify(badHash, r, s, x, y))
}

// TestGeneral_Signature tests the Signature function.
//
// testSignature is a test struct containing private, publicX, publicY,
// hash, rIn, sIn, and raw fields.
// The function iterates over the testSignature slice and performs various
// operations including signing, verifying, converting, and initializing
// variables.
//
// Parameters:
// - t: The testing.T object for running the test
// Returns:
//
//	none
func TestGeneral_Signature(t *testing.T) {
	testSignature := []struct {
		private *big.Int
		publicX *big.Int
		publicY *big.Int
		hash    *big.Int
		rIn     *big.Int
		sIn     *big.Int
		rOut    *big.Int
		sOut    *big.Int
		raw     string
	}{
		{
			private: internalUtils.StrToBig("104397037759416840641267745129360920341912682966983343798870479003077644689"),
			publicX: internalUtils.StrToBig("1913222325711601599563860015182907040361852177892954047964358042507353067365"),
			publicY: internalUtils.StrToBig("798905265292544287704154888908626830160713383708400542998012716235575472365"),
			hash:    internalUtils.StrToBig("2680576269831035412725132645807649347045997097070150916157159360688041452746"),
			rIn:     internalUtils.StrToBig("607684330780324271206686790958794501662789535258258105407533051445036595885"),
			sIn:     internalUtils.StrToBig("453590782387078613313238308551260565642934039343903827708036287031471258875"),
		},
		{
			hash: internalUtils.HexToBN("0x7f15c38ea577a26f4f553282fcfe4f1feeb8ecfaad8f221ae41abf8224cbddd"),
			rIn:  internalUtils.StrToBig("2458502865976494910213617956670505342647705497324144349552978333078363662855"),
			sIn:  internalUtils.StrToBig("3439514492576562277095748549117516048613512930236865921315982886313695689433"),
			raw:  "04033f45f07e1bd1a51b45fc24ec8c8c9908db9e42191be9e169bfcac0c0d997450319d0f53f6ca077c4fa5207819144a2a4165daef6ee47a7c1d06c0dcaa3e456",
		},
		{
			hash:    internalUtils.HexToBN("0x324df642fcc7d98b1d9941250840704f35b9ac2e3e2b58b6a034cc09adac54c"),
			publicX: internalUtils.HexToBN("0x4e52f2f40700e9cdd0f386c31a1f160d0f310504fc508a1051b747a26070d10"),
			rIn:     internalUtils.StrToBig("2849277527182985104629156126825776904262411756563556603659114084811678482647"),
			sIn:     internalUtils.StrToBig("3156340738553451171391693475354397094160428600037567299774561739201502791079"),
		},
		// Example ref: https://github.com/starkware-libs/crypto-cpp/blob/master/src/starkware/crypto/ecdsa_test.cc
		// NOTICE: s component of the {r, s} signature is not available at source, but was manually computed/confirmed as
		//   `s := sc.InvModCurveSize(w)` for w: 0x1f2c44a7798f55192f153b4c48ea5c1241fbb69e6132cc8a0da9c5b62a4286e
		{
			private: internalUtils.HexToBN("0x3c1e9550e66958296d11b60f8e8e7a7ad990d07fa65d5f7652c4a6c87d4e3cc"),
			hash:    internalUtils.HexToBN("0x397e76d1667c4454bfb83514e120583af836f8e32a516765497823eabe16a3f"),
			rIn:     internalUtils.HexToBN("0x173fd03d8b008ee7432977ac27d1e9d1a1f6c98b1a2f05fa84a21c84c44e882"),
			sIn:     internalUtils.HexToBN("4b6d75385aed025aa222f28a0adc6d58db78ff17e51c3f59e259b131cd5a1cc"),
		},
		{
			publicX: internalUtils.HexToBN("0x77a3b314db07c45076d11f62b6f9e748a39790441823307743cf00d6597ea43"),
			hash:    internalUtils.HexToBN("0x397e76d1667c4454bfb83514e120583af836f8e32a516765497823eabe16a3f"),
			rIn:     internalUtils.HexToBN("0x173fd03d8b008ee7432977ac27d1e9d1a1f6c98b1a2f05fa84a21c84c44e882"),
			sIn:     internalUtils.HexToBN("4b6d75385aed025aa222f28a0adc6d58db78ff17e51c3f59e259b131cd5a1cc"),
		},
		{
			private: internalUtils.HexToBN("0x3c1e9550e66958296d11b60f8e8e7a7ad990d07fa65d5f7652c4a6c87d4e3cc"),
			hash:    internalUtils.HexToBN("0x397e76d1667c4454bfb83514e120583af836f8e32a516765497823eabe16a3f"),
			rOut:    internalUtils.HexToBN("0x173fd03d8b008ee7432977ac27d1e9d1a1f6c98b1a2f05fa84a21c84c44e882"),
			sOut:    internalUtils.HexToBN("4b6d75385aed025aa222f28a0adc6d58db78ff17e51c3f59e259b131cd5a1cc"),
		},
	}

	var err error
	for _, tt := range testSignature {
		require := require.New(t)
		if tt.raw != "" {
			h, err := internalUtils.HexToBytes(tt.raw)
			require.NoError(err)
			tt.publicX, tt.publicY = elliptic.Unmarshal(Curve, h) //nolint:all
		} else if tt.private != nil {
			tt.publicX, tt.publicY, err = Curve.PrivateToPoint(tt.private)
			require.NoError(err)
		} else if tt.publicX != nil {
			tt.publicY = Curve.GetYCoordinate(tt.publicX)
		}

		if tt.rIn == nil && tt.private != nil {
			tt.rIn, tt.sIn, err = Curve.Sign(tt.hash, tt.private)
			require.NoError(err)
			if tt.rOut != nil && tt.rOut.Cmp(tt.rIn) != 0 {
				t.Errorf("Signature {r!, s} mismatch: %x != %x\n", tt.rIn, tt.rOut)
			}
			if tt.sOut != nil && tt.sOut.Cmp(tt.sIn) != 0 {
				t.Errorf("Signature {r, s!} mismatch: %x != %x\n", tt.sIn, tt.sOut)
			}
		}

		require.True(Curve.Verify(tt.hash, tt.rIn, tt.sIn, tt.publicX, tt.publicY))
	}
}

type ecMultOption struct {
	algo   string
	fn     EcMultiFn
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

// TestGeneral_SplitFactStr is a test function that tests the SplitFactStr function.
//
// It verifies the behavior of the SplitFactStr function by providing different inputs and checking the output.
// The function takes no parameters and returns no values.
//
// Parameters:
// - t: The testing.T object for running the test
// Returns:
//
//	none
func TestGeneral_SplitFactStr(t *testing.T) {
	data := []map[string]string{
		{"input": "0x3", "h": "0x0", "l": "0x3"},
		{"input": "0x300000000000000000000000000000000", "h": "0x3", "l": "0x0"},
	}
	for _, d := range data {
		l, h := internalUtils.SplitFactStr(d["input"]) // 0x3
		require.Equal(t, d["l"], l)
		require.Equal(t, d["h"], h)
	}
}

// TestGeneral_VerifySignature is a test function that verifies the correctness of the VerifySignature function.
//
// It checks if the signature of a given message hash is valid using the provided r, s values and the public key.
// The function takes no parameters and returns no values.
//
// Parameters:
// - t: The testing.T object for running the test
// Returns:
//
//	none
func TestGeneral_VerifySignature(t *testing.T) {
	// values verified with starknet.js

	msgHash := "0x2789daed76c8b750d5a609a706481034db9dc8b63ae01f505d21e75a8fc2336"
	r := "0x13e4e383af407f7ccc1f13195ff31a58cad97bbc6cf1d532798b8af616999d4"
	s := "0x44dd06cf67b2ba7ea4af346d80b0b439e02a0b5893c6e4dfda9ee204211c879"
	fullPubKey := "0x6c7c4408e178b2999cef9a5b3fa2a3dffc876892ad6a6bd19d1451a2256906c"

	require.True(t, VerifySignature(msgHash, r, s, fullPubKey))

	// Change the last digit of the message hash to test invalid signature
	wrongMsgHash := "0x2789daed76c8b750d5a609a706481034db9dc8b63ae01f505d21e75a8fc2337"
	require.False(t, VerifySignature(wrongMsgHash, r, s, fullPubKey))
}

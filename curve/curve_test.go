package curve

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
//   - b: a *testing.B value representing the testing context
//
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
		{
			"0x7f15c38ea577a26f4f553282fcfe4f1feeb8ecfaad8f221ae41abf8224cbddd",
			"0x13d41f388b8ea4db56c5aa6562f13359fab192b3db57651af916790f9debee9",
		},
		{
			"0x7f15c38ea577a26f4f553282fcfe4f1feeb8ecfaad8f221ae41abf8224cbddd",
			"0x7f15c38ea577a26f4f553282fcfe4f1feeb8ecfaad8f221ae41abf8224cbdde",
		},
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
//   - b: a *testing.B value representing the testing context
//
// Returns:
//
//	none
func BenchmarkCurveSign(b *testing.B) {
	MessageHash := big.NewInt(0).Exp(big.NewInt(2), big.NewInt(250), nil)
	PrivateKey := big.NewInt(0).Add(MessageHash, big.NewInt(1))
	b.ResetTimer()

	for i := int64(0); i < int64(b.N); i++ {
		b.StopTimer()
		MessageHash = big.NewInt(0).Add(MessageHash, big.NewInt(i))
		PrivateKey = big.NewInt(0).Add(PrivateKey, big.NewInt(i))
		b.StartTimer()

		result, _, _ = Sign(MessageHash, PrivateKey)
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
//   - b: a *testing.B value representing the testing context
//
// Returns:
//
//	none
func BenchmarkSignatureVerify(b *testing.B) {
	private, err := GetRandomPrivateKey()
	require.NoError(b, err)
	x, y, err := PrivateToPoint(private)
	require.NoError(b, err)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// setup
		b.StopTimer()
		randFelt, err := new(felt.Felt).SetRandom()
		require.NoError(b, err)
		hash := Pedersen(
			internalUtils.RANDOM_FELT,
			randFelt,
		)
		hashBigInt := internalUtils.FeltToBigInt(hash)
		r, s, err := Sign(hashBigInt, private)
		require.NoError(b, err)

		b.StartTimer()
		result, _ = Verify(hashBigInt, r, s, x, y)
		b.StopTimer()

		resp := result.(bool)
		require.True(b, resp)
	}
}

// TestGeneral_PrivateToPoint tests the PrivateToPoint function.
//
// Parameters:
//   - t: a *testing.T value representing the testing context
//
// Returns:
//
//	none
func TestGeneral_PrivateToPoint(t *testing.T) {
	x, _, err := PrivateToPoint(big.NewInt(2))
	require.NoError(t, err)
	expectedX, ok := new(big.Int).SetString("3324833730090626974525872402899302150520188025637965566623476530814354734325", 10)
	require.True(t, ok)

	assert.Equal(t, expectedX, x)
}

// TestGeneral_PedersenHash is a test function for the PedersenHash method in the General struct.
//
// The function tests the PedersenHash method by providing different test cases and comparing the computed hash with the expected hash.
// It uses the testing.T type from the testing package to report any errors encountered during the tests.
//
// Parameters:
//   - t: a *testing.T value representing the testing context
//
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

// TestGeneral_ComputeHashOnElements is a test function that verifies the correctness of the ComputeHashOnElements and PedersenArray functions in the General package.
//
// This function tests both functions by passing in different arrays of big.Int elements and comparing the computed hash with the expected hash.
// It checks the behavior of the functions when an empty array is passed as input, as well as when an array with multiple elements is passed.
//
// Parameters:
//   - t: a *testing.T value representing the testing context
//
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
//   - t: The testing.T object for running the test.
//
// Returns:
//
//	none
func TestGeneral_HashAndSign(t *testing.T) {
	hashy := HashPedersenElements([]*big.Int{
		big.NewInt(1953658213),
		big.NewInt(126947999705460),
		big.NewInt(1953658213),
	})

	priv, err := GetRandomPrivateKey()
	require.NoError(t, err)
	x, y, err := PrivateToPoint(priv)
	require.NoError(t, err)

	r, s, err := Sign(hashy, priv)
	require.NoError(t, err)

	resp, err := Verify(hashy, r, s, x, y)
	require.NoError(t, err)
	require.True(t, resp)
}

// TestGeneral_ComputeFact tests the ComputeFact function.
//
// It tests the ComputeFact function by providing a set of test cases
// and comparing the computed hash with the expected hash.
// The test cases consist of program hashes, program outputs,
// and expected hash values.
//
// Parameters:
//   - t: The testing.T object for running the test
//
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
			programHash: internalUtils.HexToBN("0x79920d895101ad1fbdea9adf141d8f362fdea9ee35f33dfcd07f38e4a589bab"),
			programOutput: []*big.Int{
				internalUtils.StrToBig("2754806153357301156380357983574496185342034785016738734224771556919270737441"),
			},
			expected: internalUtils.HexToBN("0x1d174fa1443deea9aab54bbca8d9be308dd14a0323dd827556c173bd132098db"),
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
//   - t: The testing.T object for running the test
//
// Returns:
//
//	none
func TestGeneral_BadSignature(t *testing.T) {
	hash := Pedersen(internalUtils.TestHexToFelt(t, "0x12773"), internalUtils.TestHexToFelt(t, "0x872362"))
	hashBigInt := internalUtils.FeltToBigInt(hash)

	priv, err := GetRandomPrivateKey()
	require.NoError(t, err)
	x, y, err := PrivateToPoint(priv)
	require.NoError(t, err)

	r, s, err := Sign(hashBigInt, priv)
	require.NoError(t, err)

	// validating the correct signatures
	result, err := Verify(hashBigInt, r, s, x, y)
	require.NoError(t, err)
	assert.True(t, result)

	// testing bad R signature
	badR := new(big.Int).Add(r, big.NewInt(1))
	result, err = Verify(hashBigInt, badR, s, x, y)
	require.NoError(t, err)
	assert.False(t, result)

	// testing bad S signature
	badS := new(big.Int).Add(s, big.NewInt(1))
	result, err = Verify(hashBigInt, r, badS, x, y)
	require.NoError(t, err)
	assert.False(t, result)

	// testing bad hash
	badHash := new(big.Int).Add(hashBigInt, big.NewInt(1))
	result, err = Verify(badHash, r, s, x, y)
	require.NoError(t, err)
	assert.False(t, result)
}

// TestGeneral_SplitFactStr is a test function that tests the SplitFactStr function.
//
// It verifies the behavior of the SplitFactStr function by providing different inputs and checking the output.
// The function takes no parameters and returns no values.
//
// Parameters:
//   - t: The testing.T object for running the test
//
// Returns:
//
//	none
func TestGeneral_SplitFactStr(t *testing.T) {
	type tescase struct {
		input string
		h     string
		l     string
		err   bool
	}
	data := []tescase{
		{
			input: "0x3",
			h:     "0x0",
			l:     "0x3",
		},
		{
			input: "0x300000000000000000000000000000000",
			h:     "0x3",
			l:     "0x0",
		},
		{
			input: "11111111111111111111111111111111111111111111111111111111111111010",
			err:   true,
		},
		{
			input: "X",
			err:   true,
		},
	}
	for _, d := range data {
		l, h, err := internalUtils.SplitFactStr(d.input)
		if d.err {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
			assert.Equal(t, d.l, l)
			assert.Equal(t, d.h, h)
		}
	}
}

// TestGeneral_VerifySignature is a test function that verifies the correctness of the VerifySignature function.
//
// It checks if the signature of a given message hash is valid using the provided r, s values and the public key.
// The function takes no parameters and returns no values.
//
// Parameters:
//   - t: The testing.T object for running the test
//
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

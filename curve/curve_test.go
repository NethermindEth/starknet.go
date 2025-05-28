package curve

import (
	"math/big"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// package level variable to be used by the benchmarking code
// to prevent the compiler from optimising the code away
var result any

// BenchmarkCurveSign benchmarks the curve.Sign function.
func BenchmarkCurveSign(b *testing.B) {
	MessageHash := big.NewInt(0).Exp(big.NewInt(2), big.NewInt(250), nil)
	PrivateKey := big.NewInt(0).Add(MessageHash, big.NewInt(1))
	b.ResetTimer()

	for range b.N {
		result, _, _ = Sign(MessageHash, PrivateKey)
	}
}

// BenchmarkSignatureVerify benchmarks the SignatureVerify function.
//
// This benchmark generates a random private key and public key pair once,
// then for each iteration creates a random hash, signs it, and measures
// only the time taken to verify the signature.
func BenchmarkSignatureVerify(b *testing.B) {
	private, x, _, err := GetRandomKeys()
	require.NoError(b, err)
	b.ResetTimer()

	for range b.N {
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
		result, _ = Verify(hashBigInt, r, s, x)
		b.StopTimer()

		resp := result.(bool)
		require.True(b, resp)
	}
}

// TestPrivateToPoint tests the PrivateToPoint function.
func TestPrivateToPoint(t *testing.T) {
	x, _ := PrivateKeyToPoint(big.NewInt(2))
	expectedX, ok := new(big.Int).SetString("3324833730090626974525872402899302150520188025637965566623476530814354734325", 10)
	require.True(t, ok)

	assert.Equal(t, expectedX, x)
}

// TestComputeHashOnElements is a test function that verifies the correctness of the
// ComputeHashOnElements and PedersenArray functions in the General package.
//
// This function tests both functions by passing in different arrays of big.Int
// elements and comparing the computed hash with the expected hash.
// It checks the behaviour of the functions when an empty array is passed as input,
// as well as when an array with multiple elements is passed.
func TestComputeHashOnElements(t *testing.T) {
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

// TestSignature tests the behaviour of the Sign and Verify functions against
// the expected values.
func TestSignature(t *testing.T) {
	hash := Pedersen(internalUtils.TestHexToFelt(t, "0x12773"), internalUtils.TestHexToFelt(t, "0x872362"))
	hashBigInt := internalUtils.FeltToBigInt(hash)

	priv, x, _, err := GetRandomKeys()
	require.NoError(t, err)

	r, s, err := Sign(hashBigInt, priv)
	require.NoError(t, err)

	// validating the correct signatures
	result, err := Verify(hashBigInt, r, s, x)
	require.NoError(t, err)
	assert.True(t, result)

	// testing bad R signature
	badR := new(big.Int).Add(r, big.NewInt(1))
	result, err = Verify(hashBigInt, badR, s, x)
	require.NoError(t, err)
	assert.False(t, result)

	// testing bad S signature
	badS := new(big.Int).Add(s, big.NewInt(1))
	result, err = Verify(hashBigInt, r, badS, x)
	require.NoError(t, err)
	assert.False(t, result)

	// testing bad hash
	badHash := new(big.Int).Add(hashBigInt, big.NewInt(1))
	result, err = Verify(badHash, r, s, x)
	require.NoError(t, err)
	assert.False(t, result)
}

// TestVerifySignature is a test function that verifies the correctness of the VerifySignature function.
//
// It checks if the signature of a given message hash is valid using the provided r, s values and the public key.
func TestVerifySignature(t *testing.T) {
	// values verified with starknet.js

	msgHash := internalUtils.TestHexToFelt(t, "0x2789daed76c8b750d5a609a706481034db9dc8b63ae01f505d21e75a8fc2336")
	r := internalUtils.TestHexToFelt(t, "0x13e4e383af407f7ccc1f13195ff31a58cad97bbc6cf1d532798b8af616999d4")
	s := internalUtils.TestHexToFelt(t, "0x44dd06cf67b2ba7ea4af346d80b0b439e02a0b5893c6e4dfda9ee204211c879")
	pubKey := internalUtils.TestHexToFelt(t, "0x6c7c4408e178b2999cef9a5b3fa2a3dffc876892ad6a6bd19d1451a2256906c")

	resp, err := VerifyFelts(msgHash, r, s, pubKey)
	require.NoError(t, err)
	require.True(t, resp)

	// Change the last digit of the message hash to test invalid signature
	wrongMsgHash := internalUtils.TestHexToFelt(t, "0x2789daed76c8b750d5a609a706481034db9dc8b63ae01f505d21e75a8fc2337")
	resp, err = VerifyFelts(wrongMsgHash, r, s, pubKey)
	require.NoError(t, err)
	require.False(t, resp)
}

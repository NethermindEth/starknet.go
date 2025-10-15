package curve

import (
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	starkcurve "github.com/consensys/gnark-crypto/ecc/stark-curve"
	"github.com/consensys/gnark-crypto/ecc/stark-curve/ecdsa"
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
			internalUtils.DeadBeef,
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
	t.Parallel()
	x, _ := PrivateKeyToPoint(big.NewInt(2))
	expectedX, ok := new(
		big.Int,
	).SetString("3324833730090626974525872402899302150520188025637965566623476530814354734325", 10)
	require.True(t, ok)

	assert.Equal(t, expectedX, x)
}

func TestPrivateKeyEndToEnd(t *testing.T) {
	t.Parallel()

	testSet := []struct {
		name  string
		privK string
		pubK  string
	}{ // taken from devnet accounts
		{
			name:  "case1",
			privK: "0x0000000000000000000000000000000085b0ed141c12d4297a9f6fa3032b9757",
			pubK:  "0x043135f5e8e5e73d9750659bb5cccc803bc63318584933d584f9b5372ee8ffa6",
		},
		{
			name:  "case2",
			privK: "0x00000000000000000000000000000000b3de4d1a7a54cb19e2fbf0897cdaa555",
			pubK:  "0x06591082275c7da568b1542044eb08c2fcf3e0c75121a5275dce4960367f2bb8",
		},
		{
			name:  "case3",
			privK: "0x00000000000000000000000000000000522e4cd212156cf8ef4052615570ad8f",
			pubK:  "0x006a78b5ad5abdb109d4d362c14895efbd45a111d5f80157f669fd127ad0c0fd",
		},
		{
			name:  "case4",
			privK: "0x000000000000000000000000000000000baa0de5814f3b01f797c26b8e4e15c5",
			pubK:  "0x021a2016d43180337d76210eb85a016d2e28e315330c70fc66151c60981b0a18",
		},
		{
			name:  "case5",
			privK: "0x0000000000000000000000000000000085dee1deeb9c5212f92ccaaae0891bc4",
			pubK:  "0x0590047e22670dd8338582556e69e2874113ecfeb96e948f38e9d8d49258bb3c",
		},
		{
			name:  "case6",
			privK: "0x000000000000000000000000000000005528b45100c856799d326bc1340a68d7",
			pubK:  "0x04c8606899ef4fa13bd87683710832ebdab5ebb8550f2781594fc819190ec478",
		},
		{
			name:  "case7",
			privK: "0x000000000000000000000000000000001d7ca805b693b571b95b6858a2d7f55b",
			pubK:  "0x04f8c79272766f492c7a753efce366e277c6da76b01f7c380c555004614f9403",
		},
		{
			name:  "case8",
			privK: "0x000000000000000000000000000000007e594f2a0862cfc474eff190c4d8a53b",
			pubK:  "0x01bc0e0a1589a5364334c91b1dbccd5e3fb723cedc842625472ba0e9ffb5a16e",
		},
		{
			name:  "case9",
			privK: "0x00000000000000000000000000000000c028449dab59500f6abef03ae61d9306",
			pubK:  "0x0559d382842465d1add1e04e1e873d34e4102dfc9a49bb64f414f3c037006e6a",
		},
		{
			name:  "case10",
			privK: "0x00000000000000000000000000000000e2f6c88bd587ea90e65db1cda8a90918",
			pubK:  "0x03796f3fbe494243b1afeb7a0921500082f6c3c4cfb2af7f963700eb4e0a6c89",
		},
		{
			name:  "without some zeroes",
			privK: "0x0000000000e2f6c88bd587ea90e65db1cda8a90918",
			pubK:  "0x3796f3fbe494243b1afeb7a0921500082f6c3c4cfb2af7f963700eb4e0a6c89",
		},
		{
			name:  "without leading zeroes",
			privK: "0xe2f6c88bd587ea90e65db1cda8a90918",
			pubK:  "0x3796f3fbe494243b1afeb7a0921500082f6c3c4cfb2af7f963700eb4e0a6c89",
		},
		{
			name:  "without '0x' prefix + without leading zeroes",
			privK: "e2f6c88bd587ea90e65db1cda8a90918",
			pubK:  "3796f3fbe494243b1afeb7a0921500082f6c3c4cfb2af7f963700eb4e0a6c89",
		},
	}

	t.Run("deriving keys and comparing", func(t *testing.T) {
		t.Parallel()
		for _, test := range testSet {
			t.Run(test.name, func(t *testing.T) {
				t.Parallel()

				privK := internalUtils.HexToBN(test.privK)

				g1Affline := starkcurve.G1Affine{}
				g1a := g1Affline.ScalarMultiplicationBase(privK)

				// ****** asserts whether a public key returned by the 'ecdsa.PublicKey' struct is
				// the same as the original public key
				var pubKeyStruct ecdsa.PublicKey
				pubKeyBytes := g1a.Bytes()
				_, err := pubKeyStruct.SetBytes(pubKeyBytes[:])
				require.NoError(t, err)

				assert.Contains(t, test.pubK, pubKeyStruct.A.X.Text(16))
				assert.Equal(t,
					internalUtils.FillHexWithZeroes(test.pubK),
					internalUtils.FillHexWithZeroes(pubKeyStruct.A.X.Text(16)))

				// ****** asserts whether a private key returned by the 'ecdsa.PrivateKey' struct is
				// the same as the original private key.

				// Assigning the private key
				var privKeyStruct ecdsa.PrivateKey
				privKeyBytes, err := fmtPrivKey(privK)
				require.NoError(t, err)
				privKeyInput := append(pubKeyStruct.Bytes(), privKeyBytes...)
				_, err = privKeyStruct.SetBytes(privKeyInput)
				require.NoError(t, err)

				// Getting the private key
				// A 64 bytes array containing both public (compressed) and private keys.
				fullPrivKBytes := privKeyStruct.Bytes()
				// The remaining 32 bytes are the private key.
				privKBytes := fullPrivKBytes[32:]
				privKey := new(big.Int).SetBytes(privKBytes)

				assert.Contains(t, test.privK, privKey.Text(16))
				assert.Equal(t,
					internalUtils.FillHexWithZeroes(test.privK),
					internalUtils.FillHexWithZeroes(privKey.Text(16)))
			})
		}
	})

	t.Run("sign and verify", func(t *testing.T) {
		t.Parallel()
		for _, test := range testSet {
			t.Run(test.name, func(t *testing.T) {
				t.Parallel()

				privK := internalUtils.HexToBN(test.privK)
				pubK := internalUtils.HexToBN(test.pubK)
				randMsh, err := rand.Int(rand.Reader, privK)
				require.NoError(t, err)

				r, s, err := Sign(randMsh, privK)
				require.NoError(t, err)

				// validating the correct signature
				result, err := Verify(randMsh, r, s, pubK)
				require.NoError(t, err)
				assert.True(t, result)

				// testing bad signature
				result, err = Verify(randMsh, r.Add(r, big.NewInt(1)), s, pubK)
				require.NoError(t, err)
				assert.False(t, result)
			})
		}
	})
}

// TestComputeHashOnElements is a test function that verifies the correctness of the
// ComputeHashOnElements and PedersenArray functions in the General package.
//
// This function tests both functions by passing in different arrays of big.Int
// elements and comparing the computed hash with the expected hash.
// It checks the behaviour of the functions when an empty array is passed as input,
// as well as when an array with multiple elements is passed.
func TestComputeHashOnElements(t *testing.T) {
	t.Parallel()
	hashEmptyArray := ComputeHashOnElements([]*big.Int{})
	hashEmptyArrayFelt := PedersenArray([]*felt.Felt{}...)

	expectedHashEmmptyArray := internalUtils.HexToBN(
		"0x49ee3eba8c1600700ee1b87eb599f16716b0b1022947733551fde4050ca6804",
	)
	require.Equal(t, hashEmptyArray, expectedHashEmmptyArray, "Hash empty array wrong value.")
	require.Equal(
		t,
		internalUtils.FeltToBigInt(hashEmptyArrayFelt),
		expectedHashEmmptyArray,
		"Hash empty array wrong value.",
	)

	filledArray := []*big.Int{
		big.NewInt(123782376),
		big.NewInt(213984),
		big.NewInt(128763521321),
	}

	hashFilledArray := ComputeHashOnElements(filledArray)
	hashFilledArrayFelt := PedersenArray(internalUtils.BigIntArrToFeltArr(filledArray)...)

	expectedHashFilledArray := internalUtils.HexToBN(
		"0x7b422405da6571242dfc245a43de3b0fe695e7021c148b918cd9cdb462cac59",
	)
	require.Equal(t, hashFilledArray, expectedHashFilledArray, "Hash filled array wrong value.")
	require.Equal(
		t,
		internalUtils.FeltToBigInt(hashFilledArrayFelt),
		expectedHashFilledArray,
		"Hash filled array wrong value.",
	)
}

// TestSignature tests the behaviour of the Sign and Verify functions against
// the expected values.
func TestSignature(t *testing.T) {
	t.Parallel()
	hash := Pedersen(
		internalUtils.TestHexToFelt(t, "0x12773"),
		internalUtils.TestHexToFelt(t, "0x872362"),
	)
	hashBigInt := internalUtils.FeltToBigInt(hash)

	priv, x, _, err := GetRandomKeys()
	require.NoError(t, err)
	privFelt := internalUtils.BigIntToFelt(priv)
	xFelt := internalUtils.BigIntToFelt(x)

	r, s, err := Sign(hashBigInt, priv)
	require.NoError(t, err)

	// testing SignFelts
	rFelt, sFelt, err := SignFelts(hash, privFelt)
	require.NoError(t, err)
	result, err := VerifyFelts(hash, rFelt, sFelt, xFelt)
	require.NoError(t, err)
	assert.True(t, result)

	// validating the correct signatures
	result, err = Verify(hashBigInt, r, s, x)
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
	t.Parallel()
	// values verified with starknet.js

	msgHash := internalUtils.TestHexToFelt(
		t,
		"0x2789daed76c8b750d5a609a706481034db9dc8b63ae01f505d21e75a8fc2336",
	)
	r := internalUtils.TestHexToFelt(
		t,
		"0x13e4e383af407f7ccc1f13195ff31a58cad97bbc6cf1d532798b8af616999d4",
	)
	s := internalUtils.TestHexToFelt(
		t,
		"0x44dd06cf67b2ba7ea4af346d80b0b439e02a0b5893c6e4dfda9ee204211c879",
	)
	pubKey := internalUtils.TestHexToFelt(
		t,
		"0x6c7c4408e178b2999cef9a5b3fa2a3dffc876892ad6a6bd19d1451a2256906c",
	)

	resp, err := VerifyFelts(msgHash, r, s, pubKey)
	require.NoError(t, err)
	require.True(t, resp)

	// Change the last digit of the message hash to test invalid signature
	wrongMsgHash := internalUtils.TestHexToFelt(
		t,
		"0x2789daed76c8b750d5a609a706481034db9dc8b63ae01f505d21e75a8fc2337",
	)
	resp, err = VerifyFelts(wrongMsgHash, r, s, pubKey)
	require.NoError(t, err)
	require.False(t, resp)
}

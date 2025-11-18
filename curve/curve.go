package curve

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"

	junoCrypto "github.com/NethermindEth/juno/core/crypto"
	"github.com/NethermindEth/juno/core/crypto/blake2s"
	"github.com/NethermindEth/juno/core/felt"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	starkcurve "github.com/consensys/gnark-crypto/ecc/stark-curve"
	"github.com/consensys/gnark-crypto/ecc/stark-curve/ecdsa"
	"github.com/consensys/gnark-crypto/ecc/stark-curve/fp"
)

// Verify verifies the validity of the signature for a given message hash using
// the StarkCurve.
//
// Parameters:
//   - msgHash: The message hash to be verified
//   - r: The r component of the signature
//   - s: The s component of the signature
//   - pubX: The x-coordinate of the public key used for verification,
//     usually used as the account public key.
//
// Returns:
//   - bool: true if the signature is valid, false otherwise
//   - error: An error if any occurred during the verification process
func Verify(msgHash, r, s, pubX *big.Int) (bool, error) {
	pubKey := junoCrypto.NewPublicKey(new(felt.Felt).SetBigInt(pubX))
	msgHashFelt := new(felt.Felt).SetBigInt(msgHash)
	rFelt := new(felt.Felt).SetBigInt(r)
	sFelt := new(felt.Felt).SetBigInt(s)

	return pubKey.Verify(&junoCrypto.Signature{R: *rFelt, S: *sFelt}, msgHashFelt)
}

// VerifyFelts verifies the validity of the signature for a given message hash
// using the StarkCurve.
// It does the same as Verify, but with felt.Felt parameters.
//
// Parameters:
//   - msgHash: The message hash to be verified
//   - r: The r component of the signature
//   - s: The s component of the signature
//   - pubX: The x-coordinate of the public key used for verification,
//     usually used as the account public key.
//
// Returns:
//   - bool: true if the signature is valid, false otherwise
//   - error: An error if any occurred during the verification process
func VerifyFelts(msgHash, r, s, pubX *felt.Felt) (bool, error) {
	pubKey := junoCrypto.NewPublicKey(pubX)

	return pubKey.Verify(&junoCrypto.Signature{R: *r, S: *s}, msgHash)
}

// Sign calculates the signature of a message using the StarkCurve algorithm.
//
// Parameters:
//   - msgHash: The message hash to be signed
//   - privKey: The private key used for signing
//
// Returns:
//   - r: The r component of the signature
//   - s: The s component of the signature
//   - error: An error if any occurred during the signing process
func Sign(msgHash, privKey *big.Int) (r, s *big.Int, err error) {
	var g1Affline starkcurve.G1Affine
	// generating pub and priv key types from the 'privKey' parameter
	g1a := g1Affline.ScalarMultiplicationBase(privKey)

	// pub
	var pubKeyStruct ecdsa.PublicKey
	pubKeyBytes := g1a.Bytes()
	_, err = pubKeyStruct.SetBytes(pubKeyBytes[:])
	if err != nil {
		return nil, nil, err
	}

	// priv
	var privKeyStruct ecdsa.PrivateKey
	privKeyBytes, err := fmtPrivKey(privKey)
	if err != nil {
		return nil, nil, err
	}
	privKeyInput := append(pubKeyStruct.Bytes(), privKeyBytes...)
	_, err = privKeyStruct.SetBytes(privKeyInput)
	if err != nil {
		return nil, nil, err
	}

	// signing
	_, r, s, err = privKeyStruct.SignForRecover(msgHash.Bytes(), nil)

	return r, s, err
}

// SignFelts calculates the signature of a message using the StarkCurve algorithm.
// It does the same as Sign, but with felt.Felt parameters.
//
// Parameters:
//   - msgHash: The message hash to be signed
//   - privKey: The private key used for signing
//
// Returns:
//   - r: The r component of the signature
//   - s: The s component of the signature
//   - error: An error if any occurred during the signing process
func SignFelts(msgHash, privKey *felt.Felt) (r, s *felt.Felt, err error) {
	msgHashBig := msgHash.BigInt(new(big.Int))
	privKeyBig := privKey.BigInt(new(big.Int))

	rBig, sBig, err := Sign(msgHashBig, privKeyBig)
	if err != nil {
		return nil, nil, err
	}

	r = new(felt.Felt).SetBigInt(rBig)
	s = new(felt.Felt).SetBigInt(sBig)

	return r, s, nil
}

// GetRandomKeys generates a random private key and its corresponding public key.
//
// Returns:
//   - privKey: The private key
//   - x: The x-coordinate of the public key (the Starknet public key)
//   - y: The y-coordinate of the public key
//   - err: An error if any occurred during the key generation process
func GetRandomKeys() (privKey, x, y *big.Int, err error) {
	fullPrivK, err := ecdsa.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, nil, err
	}

	// A 64 bytes array containing both public (compressed) and private keys.
	fullPrivKBytes := fullPrivK.Bytes()
	// The remaining 32 bytes are the private key.
	privKBytes := fullPrivKBytes[32:]

	privKey = new(big.Int).SetBytes(privKBytes)

	return privKey,
		fullPrivK.PublicKey.A.X.BigInt(new(big.Int)),
		fullPrivK.PublicKey.A.Y.BigInt(new(big.Int)),
		nil
}

// PrivateKeyToPoint generates a point on the StarkCurve from a private key.
//
// It takes a private key as a parameter and returns the x and y coordinates of
// the generated point on the curve.
//
// Parameters:
//   - privKey: The private key
//
// Returns:
//   - x: The x-coordinate of the point on the curve
//   - y: The y-coordinate of the point on the curve
func PrivateKeyToPoint(privKey *big.Int) (x, y *big.Int) {
	var g1Affline starkcurve.G1Affine
	res := g1Affline.ScalarMultiplicationBase(privKey)

	return res.X.BigInt(new(big.Int)), res.Y.BigInt(new(big.Int))
}

// GetYCoordinate returns the y-coordinate of a point on the curve given the x-coordinate.
//
// Parameters:
//   - starkX: The x-coordinate of the point
//
// Returns:
//   - *big.Int: The y-coordinate of the point
func GetYCoordinate(starkX *felt.Felt) *felt.Felt {
	//nolint:lll // The link would be unclickable if we break the line.
	// ref: https://github.com/NethermindEth/juno/blob/7d64642de90b6957c40a3b3ea75e6ad548a37f39/core/crypto/ecdsa.go#L26
	xEl := starkX.Impl()

	var ySquared fp.Element
	ySquared.Mul(xEl, xEl).Mul(&ySquared, xEl) // x^3
	ySquared.Add(&ySquared, xEl)               // + x

	_, b := starkcurve.CurveCoefficients()
	ySquared.Add(&ySquared, &b) // ySquared equals to (x^3 + x + b)

	starkY := ySquared.Sqrt(&ySquared)
	yFelt := felt.Felt(*starkY)

	return &yFelt
}

// HashPedersenElements calculates the hash of a list of elements using a
// golang Pedersen Hash.
//
// Parameters:
//   - elems: slice of big.Int pointers to be hashed
//
// Returns:
//   - hash: The hash of the list of elements
func HashPedersenElements(elems []*big.Int) (hash *big.Int) {
	feltArr := internalUtils.BigIntArrToFeltArr(elems)
	if len(elems) == 0 {
		feltArr = append(feltArr, new(felt.Felt))
	}

	feltHash := new(felt.Felt)
	for _, felt := range feltArr {
		feltHash = Pedersen(feltHash, felt)
	}

	hash = internalUtils.FeltToBigInt(feltHash)

	return hash
}

// ComputeHashOnElements computes the hash on the given elements using a
// golang Pedersen Hash implementation.
//
// The function appends the length of `elems` to the slice and then calls the
// `HashPedersenElements` method passing in `elems` as an argument. The
// resulting hash is returned.
//
// Parameters:
//   - elems: slice of big.Int pointers to be hashed
//
// Returns:
//   - hash: The hash of the list of elements
func ComputeHashOnElements(elems []*big.Int) (hash *big.Int) {
	elems = append(elems, big.NewInt(int64(len(elems))))

	return HashPedersenElements(elems)
}

// Pedersen is a function that implements the Pedersen hash.
//
// Parameters:
//   - a: a pointers to felt.Felt to be hashed.
//   - b: a pointers to felt.Felt to be hashed.
//
// Returns:
//   - *felt.Felt: a pointer to a felt.Felt storing the resulting hash.
func Pedersen(a, b *felt.Felt) *felt.Felt {
	hash := junoCrypto.Pedersen(a, b)

	return &hash
}

// Poseidon is a function that implements the Poseidon hash.
//
// Parameters:
//   - a: a pointers to felt.Felt to be hashed.
//   - b: a pointers to felt.Felt to be hashed.
//
// Returns:
//   - *felt.Felt: a pointer to a felt.Felt storing the resulting hash.
func Poseidon(a, b *felt.Felt) *felt.Felt {
	hash := junoCrypto.Poseidon(a, b)

	return &hash
}

// Blake2s is a function that implements the Blake2s hash.
//
// Parameters:
//   - a: a pointers to felt.Felt to be hashed.
//   - b: a pointers to felt.Felt to be hashed.
//
// Returns:
//   - *felt.Felt: a pointer to a felt.Felt storing the resulting hash.
func Blake2s(a, b *felt.Felt) *felt.Felt {
	hash := blake2s.Blake2s(a, b)
	hashFelt := felt.Felt(hash)

	return &hashFelt
}

// PedersenArray is a function that takes a variadic number of felt.Felt
// pointers as parameters and calls the PedersenArray function from the
// junoCrypto package with the provided parameters.
//
// Parameters:
//   - felts: A variadic number of pointers to felt.Felt
//
// Returns:
//   - *felt.Felt: pointer to a felt.Felt
func PedersenArray(felts ...*felt.Felt) *felt.Felt {
	hash := junoCrypto.PedersenArray(felts...)

	return &hash
}

// PoseidonArray is a function that takes a variadic number of felt.Felt
// pointers as parameters and calls the PoseidonArray function from the
// junoCrypto package with the provided parameters.
//
// Parameters:
//   - felts: A variadic number of pointers to felt.Felt
//
// Returns:
//   - *felt.Felt: pointer to a felt.Felt
func PoseidonArray(felts ...*felt.Felt) *felt.Felt {
	hash := junoCrypto.PoseidonArray(felts...)

	return &hash
}

// Blake2sArray is a function that takes a variadic number of felt.Felt
// pointers as parameters and calls the Blake2sArray function from the
// junoCrypto package with the provided parameters.
//
// Parameters:
//   - felts: A variadic number of pointers to felt.Felt
//
// Returns:
//   - *felt.Felt: pointer to a felt.Felt
func Blake2sArray(felts ...*felt.Felt) *felt.Felt {
	hash := blake2s.Blake2sArray(felts...)
	hashFelt := felt.Felt(hash)

	return &hashFelt
}

// StarknetKeccak computes the Starknet Keccak hash of the given byte slice.
//
// Parameters:
//   - b: The byte slice to hash
//
// Returns:
//   - *felt.Felt: pointer to a felt.Felt
//   - error: An error if any
func StarknetKeccak(b []byte) *felt.Felt {
	hash := junoCrypto.StarknetKeccak(b)

	return &hash
}

// fmtPrivKey formats a private key to a 32 bytes array by padding it
// with leading zeroes if necessary, which is required by the ecdsa.PrivateKey type.
func fmtPrivKey(privKey *big.Int) ([]byte, error) {
	return hex.DecodeString(fmt.Sprintf("%064s", privKey.Text(16))) //nolint:mnd // hex base
}

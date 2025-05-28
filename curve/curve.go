package curve

/*
	Although the library adheres to the 'elliptic/curve' interface.
	All testing has been done against library function explicity.
	It is recommended to use in the same way(i.e. `curve.Sign` and not `ecdsa.Sign`).
*/
import (
	"crypto/rand"
	_ "embed"
	"math/big"

	"github.com/NethermindEth/juno/core/crypto"
	junoCrypto "github.com/NethermindEth/juno/core/crypto"
	"github.com/NethermindEth/juno/core/felt"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	starkcurve "github.com/consensys/gnark-crypto/ecc/stark-curve"
	"github.com/consensys/gnark-crypto/ecc/stark-curve/ecdsa"
	"github.com/consensys/gnark-crypto/ecc/stark-curve/fp"
)

// HashPedersenElements calculates the hash of a list of elements using a golang Pedersen Hash.
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
	return
}

// ComputeHashOnElements computes the hash on the given elements using a golang Pedersen Hash implementation.
//
// The function appends the length of `elems` to the slice and then calls the `HashPedersenElements` method
// passing in `elems` as an argument. The resulting hash is returned.
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
// NOTE: This function just wraps the Juno implementation
// (ref: https://github.com/NethermindEth/juno/blob/32fd743c774ec11a1bb2ce3dceecb57515f4873e/core/crypto/pedersen_hash.go#L20)
//
// Parameters:
//   - a: a pointers to felt.Felt to be hashed.
//   - b: a pointers to felt.Felt to be hashed.
//
// Returns:
//   - *felt.Felt: a pointer to a felt.Felt storing the resulting hash.
func Pedersen(a, b *felt.Felt) *felt.Felt {
	return junoCrypto.Pedersen(a, b)
}

// Poseidon is a function that implements the Poseidon hash.
// NOTE: This function just wraps the Juno implementation
// (ref: https://github.com/NethermindEth/juno/blob/32fd743c774ec11a1bb2ce3dceecb57515f4873e/core/crypto/poseidon_hash.go#L59)
//
// Parameters:
//   - a: a pointers to felt.Felt to be hashed.
//   - b: a pointers to felt.Felt to be hashed.
//
// Returns:
//   - *felt.Felt: a pointer to a felt.Felt storing the resulting hash.
func Poseidon(a, b *felt.Felt) *felt.Felt {
	return junoCrypto.Poseidon(a, b)
}

// PedersenArray is a function that takes a variadic number of felt.Felt pointers as parameters and
// calls the PedersenArray function from the junoCrypto package with the provided parameters.
// NOTE: This function just wraps the Juno implementation
// (ref: https://github.com/NethermindEth/juno/blob/32fd743c774ec11a1bb2ce3dceecb57515f4873e/core/crypto/pedersen_hash.go#L12)
//
// Parameters:
//   - felts: A variadic number of pointers to felt.Felt
//
// Returns:
//   - *felt.Felt: pointer to a felt.Felt
func PedersenArray(felts ...*felt.Felt) *felt.Felt {
	return junoCrypto.PedersenArray(felts...)
}

// PoseidonArray is a function that takes a variadic number of felt.Felt pointers as parameters and
// calls the PoseidonArray function from the junoCrypto package with the provided parameters.
// NOTE: This function just wraps the Juno implementation
// (ref: https://github.com/NethermindEth/juno/blob/main/core/crypto/poseidon_hash.go#L74)
//
// Parameters:
//   - felts: A variadic number of pointers to felt.Felt
//
// Returns:
//   - *felt.Felt: pointer to a felt.Felt
func PoseidonArray(felts ...*felt.Felt) *felt.Felt {
	return junoCrypto.PoseidonArray(felts...)
}

// StarknetKeccak computes the Starknet Keccak hash of the given byte slice.
// NOTE: This function just wraps the Juno implementation
// (ref: https://github.com/NethermindEth/juno/blob/main/core/crypto/keccak.go#L11)
//
// Parameters:
//   - b: The byte slice to hash
//
// Returns:
//   - *felt.Felt: pointer to a felt.Felt
//   - error: An error if any
func StarknetKeccak(b []byte) *felt.Felt {
	return junoCrypto.StarknetKeccak(b)
}

// VerifySignature verifies the ECDSA signature of a given message hash using the provided public key.
//
// It takes the message hash, the r and s values of the signature, and the public key as strings and
// verifies the signature using the public key.
//
// Parameters:
//   - msgHash: The hash of the message to be verified as a string
//   - r: The r value (the first part) of the signature as a string
//   - s: The s value (the second part) of the signature as a string
//   - pubKey: The public key (only the x coordinate) as a string
//
// Return values:
//   - bool: A boolean indicating whether the signature is valid
//   - error: An error if any occurred during the verification process
func VerifySignature(msgHash, r, s, pubKey string) bool {
	feltMsgHash, err := new(felt.Felt).SetString(msgHash)
	if err != nil {
		return false
	}
	feltR, err := new(felt.Felt).SetString(r)
	if err != nil {
		return false
	}
	feltS, err := new(felt.Felt).SetString(s)
	if err != nil {
		return false
	}
	pubKeyFelt, err := new(felt.Felt).SetString(pubKey)
	if err != nil {
		return false
	}

	signature := junoCrypto.Signature{
		R: *feltR,
		S: *feltS,
	}

	pubKeyStruct := junoCrypto.NewPublicKey(pubKeyFelt)
	resp, err := pubKeyStruct.Verify(&signature, feltMsgHash)
	if err != nil {
		return false
	}

	return resp
}

func GetYCoordinate(starkX *big.Int) *big.Int {
	// ref: https://github.com/NethermindEth/juno/blob/7d64642de90b6957c40a3b3ea75e6ad548a37f39/core/crypto/ecdsa.go#L26
	xEl := new(fp.Element).SetBigInt(starkX)

	var ySquared fp.Element
	ySquared.Mul(xEl, xEl).Mul(&ySquared, xEl) // x^3
	ySquared.Add(&ySquared, xEl)               // + x

	_, b := starkcurve.CurveCoefficients()
	ySquared.Add(&ySquared, &b) // ySquared equals to (x^3 + x + b)
	return ySquared.Sqrt(&ySquared).BigInt(new(big.Int))
}

func Verify(msgHash, r, s, pubX, pubY *big.Int) (bool, error) {
	pubKey := crypto.NewPublicKey(new(felt.Felt).SetBigInt(pubX))
	msgHashFelt := new(felt.Felt).SetBigInt(msgHash)
	rFelt := new(felt.Felt).SetBigInt(r)
	sFelt := new(felt.Felt).SetBigInt(s)

	return pubKey.Verify(&crypto.Signature{R: *rFelt, S: *sFelt}, msgHashFelt)
}

func Sign(msgHash, privKey *big.Int, seed ...*big.Int) (r, s *big.Int, err error) {
	// generating pub and priv keys
	g1a := new(starkcurve.G1Affine).ScalarMultiplicationBase(privKey)

	var pubKeyStruct ecdsa.PublicKey
	pubKeyBytes := g1a.Bytes()
	_, err = pubKeyStruct.SetBytes(pubKeyBytes[:])
	if err != nil {
		return nil, nil, err
	}

	var privKeyStruct ecdsa.PrivateKey
	privKeyBytes := privKey.Bytes()
	privKeyInput := append(pubKeyStruct.Bytes(), privKeyBytes...)
	_, err = privKeyStruct.SetBytes(privKeyInput[:])
	if err != nil {
		return nil, nil, err
	}

	// signing
	_, r, s, err = privKeyStruct.SignForRecover(msgHash.Bytes(), nil)

	return r, s, err
}

func GetRandomPrivateKey() (*big.Int, error) {
	priv, err := ecdsa.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	privBytes := priv.Bytes()           // a 64 bytes array containing both public (compressed) and private keys
	finalPrivKeyBytes := privBytes[32:] // the remaining 32 bytes are the private key

	finalPrivKey := new(big.Int).SetBytes(finalPrivKeyBytes)

	return finalPrivKey, nil
}

func PrivateToPoint(privKey *big.Int) (x, y *big.Int, err error) {
	g1a := new(starkcurve.G1Affine)
	res := g1a.ScalarMultiplicationBase(privKey)
	return res.X.BigInt(new(big.Int)), res.Y.BigInt(new(big.Int)), nil
}

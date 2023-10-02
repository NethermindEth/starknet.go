package main

import (
	"fmt"
	"math/big"

	starknetgo "github.com/NethermindEth/starknet.go"
	"github.com/NethermindEth/starknet.go/types"
)

// main is the entry point of the curve example.
//
// This function performs various cryptographic operations using the StarknetGo library.
// It calculates the hash using the PedersenHash function, generates a random private key,
// converts the private key to a point on the curve, signs the hash using the private key,
// and verifies the signature. The function prints whether the signature is valid or invalid.
//
// No parameters are passed to this function.
// No return types are specified.
func main() {
	/*
		Although the library adheres to the 'elliptic/curve' interface.
		All testing has been done against library function explicity.
		It is recommended to use in the same way(i.e. `curve.Sign` and not `ecdsa.Sign`).
		NOTE: when not given local file path this pulls the curve data from Starkware github repo
	*/
	hash, err := starknetgo.Curve.PedersenHash([]*big.Int{types.HexToBN("0x12773"), types.HexToBN("0x872362")})
	if err != nil {
		panic(err.Error())
	}

	priv, _ := starknetgo.Curve.GetRandomPrivateKey()

	x, y, err := starknetgo.Curve.PrivateToPoint(priv)
	if err != nil {
		panic(err.Error())
	}

	r, s, err := starknetgo.Curve.Sign(hash, priv)
	if err != nil {
		panic(err.Error())
	}

	if starknetgo.Curve.Verify(hash, r, s, x, y) {
		fmt.Println("signature is valid")
	} else {
		fmt.Println("signature is invalid")
	}
}

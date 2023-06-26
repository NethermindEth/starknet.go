package main

import (
	"fmt"
	"math/big"

	"github.com/smartcontractkit/caigo"
	"github.com/smartcontractkit/caigo/types"
)

func main() {
	/*
		Although the library adheres to the 'elliptic/curve' interface.
		All testing has been done against library function explicity.
		It is recommended to use in the same way(i.e. `curve.Sign` and not `ecdsa.Sign`).
		NOTE: when not given local file path this pulls the curve data from Starkware github repo
	*/
	hash, err := caigo.Curve.PedersenHash([]*big.Int{types.HexToBN("0x12773"), types.HexToBN("0x872362")})
	if err != nil {
		panic(err.Error())
	}

	priv, _ := caigo.Curve.GetRandomPrivateKey()

	x, y, err := caigo.Curve.PrivateToPoint(priv)
	if err != nil {
		panic(err.Error())
	}

	r, s, err := caigo.Curve.Sign(hash, priv)
	if err != nil {
		panic(err.Error())
	}

	if caigo.Curve.Verify(hash, r, s, x, y) {
		fmt.Println("signature is valid")
	} else {
		fmt.Println("signature is invalid")
	}
}

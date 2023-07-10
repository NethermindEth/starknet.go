package main

import (
	"fmt"
	"math/big"

	"github.com/NethermindEth/starknet.go"
	"github.com/NethermindEth/starknet.go/types"
)

func main() {
	/*
		Although the library adheres to the 'elliptic/curve' interface.
		All testing has been done against library function explicity.
		It is recommended to use in the same way(i.e. `curve.Sign` and not `ecdsa.Sign`).
		NOTE: when not given local file path this pulls the curve data from Starkware github repo
	*/
	hash, err := starknet.go.Curve.PedersenHash([]*big.Int{types.HexToBN("0x12773"), types.HexToBN("0x872362")})
	if err != nil {
		panic(err.Error())
	}

	priv, _ := starknet.go.Curve.GetRandomPrivateKey()

	x, y, err := starknet.go.Curve.PrivateToPoint(priv)
	if err != nil {
		panic(err.Error())
	}

	r, s, err := starknet.go.Curve.Sign(hash, priv)
	if err != nil {
		panic(err.Error())
	}

	if starknet.go.Curve.Verify(hash, r, s, x, y) {
		fmt.Println("signature is valid")
	} else {
		fmt.Println("signature is invalid")
	}
}

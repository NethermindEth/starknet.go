package main

import (
	"fmt"
	"math/big"

	"github.com/dontpanicdao/caigo"
)

func main() {
	/*
		Although the library adheres to the 'elliptic/curve' interface.
		All testing has been done against library function explicity.
		It is recommended to use in the same way(i.e. `curve.Sign` and not `ecdsa.Sign`).
		NOTE: when not given local file path this pulls the curve data from Starkware github repo
	*/
	curve, err := caigo.SC(caigo.WithConstants())
	if err != nil {
		panic(err.Error())
	}
	v1, _ := big.NewInt(0).SetString("0x12773", 0)
	v2, _ := big.NewInt(0).SetString("0x872362", 0)
	hash, err := curve.PedersenHash([]*big.Int{v1, v2})
	if err != nil {
		panic(err.Error())
	}

	priv, _ := curve.GetRandomPrivateKey()

	x, y, err := curve.PrivateToPoint(priv)
	if err != nil {
		panic(err.Error())
	}

	r, s, err := curve.Sign(hash, priv)
	if err != nil {
		panic(err.Error())
	}

	if curve.Verify(hash, r, s, x, y) {
		fmt.Println("signature is valid")
	} else {
		fmt.Println("signature is invalid")
	}
}

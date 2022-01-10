<h1 align="center">Golang Library for StarkNet/Cairo</h1>

<p align="center">
    <a href="https://pkg.go.dev/github.com/dontpanicdao/caigo">
        <img src="https://pkg.go.dev/badge/github.com/dontpanicdao/caigo.svg" alt="Go Reference">
    </a>
    <a href="https://github.com/dontpanicdao/caigo/blob/main/LICENSE">
        <img src="https://img.shields.io/badge/license-MIT-black">
    </a>
    <a href="https://starkware.co/">
        <img src="https://img.shields.io/badge/powered_by-StarkWare-navy">
    </a>
</p>

### Caigo is predominately a transcription of the following libraries:
- https://github.com/starkware-libs/cairo-lang/blob/master/src/starkware/crypto/starkware/crypto/signature
- https://github.com/seanjameshan/starknet.js
- https://github.com/software-mansion/starknet.py
- https://github.com/codahale/rfc6979/blob/master/rfc6979.go

### !!! THIS LIBRARY HAS YET BEEN AUDITED BY THE STARKWARE TEAM !!!

### Usage
Although the library adheres to the 'elliptic/curve' interface. All testing has been done against library function explicity. It is recommended to use in the same way(i.e. `curve.Sign` and not `ecdsa.Sign`).

####
```go
package main

import (
	"fmt"
	"math/big"

	"github.com/dontpanicdao/caigo"
)

func main() {
	// NOTE: when not given local file path this pulls the curve data from Starkware github repo
	curve, err := caigo.SCWithConstants("")
	if err != nil {
		panic(err.Error())
	}

	hash, err := curve.PedersenHash([]*big.Int{caigo.HexToBN("0x12773"), caigo.HexToBN("0x872362")})
	if err != nil {
		panic(err.Error())
	}

	priv := curve.GetRandomPrivateKey()

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
```

#### Test
```go
go test -v
=== RUN   TestBadSignature
--- PASS: TestBadSignature (0.06s)
=== RUN   TestKnownSignature
--- PASS: TestKnownSignature (0.02s)
=== RUN   TestDerivedSignature
--- PASS: TestDerivedSignature (0.01s)
=== RUN   TestTransactionHash
--- PASS: TestTransactionHash (0.02s)
=== RUN   TestVerifySignature
--- PASS: TestVerifySignature (0.01s)
=== RUN   TestUIVerifySignature
--- PASS: TestUIVerifySignature (0.02s)
=== RUN   TestPedersenHash
--- PASS: TestPedersenHash (0.00s)
=== RUN   TestInitCurveWithConstants
--- PASS: TestInitCurveWithConstants (0.01s)
=== RUN   TestDivMod
--- PASS: TestDivMod (0.00s)
=== RUN   TestEcMult
--- PASS: TestEcMult (0.00s)
=== RUN   TestAdd
--- PASS: TestAdd (0.00s)
=== RUN   TestMultAir
--- PASS: TestMultAir (0.00s)
=== RUN   TestGetY
--- PASS: TestGetY (0.00s)
PASS
ok      github.com/dontpanicdao/caigo   0.454s
```

## Issues

If you find an issue/bug or have a feature request please submit an issue here
[Issues](https://github.com/dontpanicdao/caigo/issues)

## Contributing

If you are looking to contribute, please head to the
[Contributing](https://github.com/dontpanicdao/caigo/blob/main/CONTRIBUTING.md) section.
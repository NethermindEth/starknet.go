<h1 align="center">Golang Library for StarkNet</h1>

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

## Examples

### deploy/call/invoke
Deploy a compiled contract to testnet, query the initial state, and invoke a state transition.

```go
package main

import (
	"fmt"

	"github.com/dontpanicdao/caigo"
)

func main() {
	// init the stark curve with constants
	// will pull the 'pedersen_params.json' file if you don't have it locally
	curve, err := caigo.SCWithConstants("")
	if err != nil {
		panic(err.Error())
	}
	
	// init starknet gateway client
	gw := caigo.NewGateway() //defaults to goerli

	// get random value for salt
	priv, _ := curve.GetRandomPrivateKey()

	// deploy StarkNet contract with random salt
	deployRequest := caigo.DeployRequest{
		ContractAddressSalt: caigo.BigToHex(priv),
		ConstructorCalldata: []string{},
	}

	// example: https://github.com/starknet-edu/ultimate-env/blob/main/counter.cairo
	deployResponse, err := gw.Deploy("counter_compiled.json", deployRequest)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Deployment Response: \n\t%+v\n\n", deployResponse)

	// poll until the desired transaction status
	pollInterval := 5
	n, status, err := gw.PollTx(deployResponse.TransactionHash, caigo.ACCEPTED_ON_L2, pollInterval, 150)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Poll %dsec %dx \n\ttransaction(%s) status: %s\n\n", n * pollInterval, n, deployResponse.TransactionHash, status)

	// fetch transaction details
	tx, err := gw.GetTransaction(deployResponse.TransactionHash)
	if err != nil {
		panic(err.Error())
	}

	// call StarkNet contract
	callResp, err := gw.Call(caigo.StarknetRequest{
		ContractAddress:    tx.Transaction.ContractAddress,
		EntryPointSelector: caigo.BigToHex(caigo.GetSelectorFromName("get_count")),
	})
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Counter is currently at: ", callResp[0])
	
	// invoke StarkNet contract external function
	invResp, err := gw.Invoke(caigo.StarknetRequest{
		ContractAddress:    tx.Transaction.ContractAddress,
		EntryPointSelector: caigo.BigToHex(caigo.GetSelectorFromName("increment")),
	})
	if err != nil {
		panic(err.Error())
	}

	n, status, err = gw.PollTx(invResp.TransactionHash, caigo.ACCEPTED_ON_L2, 5, 150)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Poll %dsec %dx \n\ttransaction(%s) status: %s\n\n", n * pollInterval, n, deployResponse.TransactionHash, status)

	callResp, err = gw.Call(caigo.StarknetRequest{
		ContractAddress:    tx.Transaction.ContractAddress,
		EntryPointSelector: caigo.BigToHex(caigo.GetSelectorFromName("get_count")),
	})
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Counter is currently at: ", callResp[0])
}
```

### sign/verify
Although the library adheres to the 'elliptic/curve' interface. All testing has been done against library function explicity. It is recommended to use in the same way(i.e. `curve.Sign` and not `ecdsa.Sign`).

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
```

### Test && Benchmark
```go
// run tests
go test -v

// run benchmarks
go test -bench=.
```

## Issues

If you find an issue/bug or have a feature request please submit an issue here
[Issues](https://github.com/dontpanicdao/caigo/issues)

## Contributing

If you are looking to contribute, please head to the
[Contributing](https://github.com/dontpanicdao/caigo/blob/main/CONTRIBUTING.md) section.
---
sidebar_position: 4
---

# Cryptography API Reference

StarkNet.go provides cryptographic utilities for working with StarkNet's cryptographic primitives. This reference documents the available cryptographic methods and how to use them.

## Stark Curve

StarkNet uses the Stark curve for cryptographic operations. The `curve` package provides methods for working with the Stark curve:

```go
package main

import (
    "fmt"
    "math/big"
    
    "github.com/NethermindEth/starknet.go/curve"
    "github.com/NethermindEth/starknet.go/utils"
)

func main() {
    // Create a new Stark curve instance
    starkCurve := curve.NewStarkCurve()
    
    // Generate a private key
    privateKey, err := utils.GeneratePrivateKey()
    if err != nil {
        panic(err)
    }
    
    // Get the public key from the private key
    publicKey, err := starkCurve.GetPublicKey(privateKey)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Private key: 0x%s\n", privateKey.Text(16))
    fmt.Printf("Public key: 0x%s\n", publicKey.Text(16))
}
```

### Stark Curve Parameters

The Stark curve is defined by the following parameters:

```go
// Stark curve parameters
var (
    StarkCurveOrder, _ = new(big.Int).SetString("3618502788666131213697322783095070105623107215331596699973092056135872020481", 10)
    StarkCurveP, _     = new(big.Int).SetString("3618502788666131213697322783095070105623107215331596699973092056135872020481", 10)
    StarkCurveAlpha, _ = new(big.Int).SetString("1", 10)
    StarkCurveBeta, _  = new(big.Int).SetString("3141592653589793238462643383279502884197169399375105820974944592307816406665", 10)
)
```

### Signing Messages

To sign a message using the Stark curve:

```go
// Message to sign
message := big.NewInt(123)

// Sign the message
r, s, err := starkCurve.Sign(privateKey, message)
if err != nil {
    panic(err)
}

fmt.Printf("Signature (r): 0x%s\n", r.Text(16))
fmt.Printf("Signature (s): 0x%s\n", s.Text(16))
```

### Verifying Signatures

To verify a signature using the Stark curve:

```go
// Verify the signature
isValid, err := starkCurve.Verify(publicKey, message, r, s)
if err != nil {
    panic(err)
}

if isValid {
    fmt.Println("Signature is valid")
} else {
    fmt.Println("Signature is invalid")
}
```

## Pedersen Hash

StarkNet uses the Pedersen hash function for various cryptographic operations. The `hash` package provides methods for computing Pedersen hashes:

```go
package main

import (
    "fmt"
    "math/big"
    
    "github.com/NethermindEth/starknet.go/hash"
)

func main() {
    // Values to hash
    a := big.NewInt(123)
    b := big.NewInt(456)
    
    // Compute the Pedersen hash
    result, err := hash.PedersenHash([]*big.Int{a, b})
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Pedersen hash: 0x%s\n", result.Text(16))
}
```

### Computing Pedersen Hashes

To compute a Pedersen hash of multiple values:

```go
// Values to hash
values := []*big.Int{
    big.NewInt(123),
    big.NewInt(456),
    big.NewInt(789),
}

// Compute the Pedersen hash
result, err := hash.PedersenHash(values)
if err != nil {
    panic(err)
}

fmt.Printf("Pedersen hash: 0x%s\n", result.Text(16))
```

### Computing Array Pedersen Hashes

To compute a Pedersen hash of an array:

```go
// Array to hash
array := []*big.Int{
    big.NewInt(123),
    big.NewInt(456),
    big.NewInt(789),
}

// Compute the array Pedersen hash
result, err := hash.ComputeArrayPedersenHash(array)
if err != nil {
    panic(err)
}

fmt.Printf("Array Pedersen hash: 0x%s\n", result.Text(16))
```

## Typed Data

StarkNet supports typed data signing, similar to Ethereum's EIP-712. The `typedData` package provides methods for working with typed data:

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/NethermindEth/starknet.go/account"
    "github.com/NethermindEth/starknet.go/typedData"
)

func main() {
    // Initialize provider and account (see Account Management section)
    // ... (provider and account setup code)
    
    // Define typed data
    typedDataJSON := `{
        "types": {
            "StarkNetDomain": {
                "name": "string",
                "version": "felt",
                "chainId": "felt"
            },
            "Person": {
                "name": "string",
                "wallet": "felt"
            },
            "Mail": {
                "from": "Person",
                "to": "Person",
                "contents": "string"
            }
        },
        "primaryType": "Mail",
        "domain": {
            "name": "StarkNet Mail",
            "version": "1",
            "chainId": "1"
        },
        "message": {
            "from": {
                "name": "Alice",
                "wallet": "0x1234567890123456789012345678901234567890"
            },
            "to": {
                "name": "Bob",
                "wallet": "0x0987654321098765432109876543210987654321"
            },
            "contents": "Hello, Bob!"
        }
    }`
    
    // Parse typed data
    td, err := typedData.TypedDataFromJSON([]byte(typedDataJSON))
    if err != nil {
        panic(err)
    }
    
    // Get the message hash
    messageHash, err := td.GetMessageHash(context.Background())
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Message hash: 0x%s\n", messageHash.Text(16))
    
    // Sign the typed data
    signature, err := acc.SignMessage(context.Background(), []string{messageHash.Text(16)})
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Signature: %v\n", signature)
    
    // Verify the signature
    isValid, err := acc.VerifyMessageSignature(context.Background(), []string{messageHash.Text(16)}, signature)
    if err != nil {
        panic(err)
    }
    
    if isValid {
        fmt.Println("Signature is valid")
    } else {
        fmt.Println("Signature is invalid")
    }
}
```

### Typed Data Structure

The `TypedData` struct represents typed data:

```go
type TypedData struct {
    Types       map[string][]TypedDataField
    PrimaryType string
    Domain      map[string]interface{}
    Message     map[string]interface{}
}
```

### Getting the Message Hash

To get the message hash of typed data:

```go
// Get the message hash
messageHash, err := td.GetMessageHash(context.Background())
if err != nil {
    panic(err)
}

fmt.Printf("Message hash: 0x%s\n", messageHash.Text(16))
```

### Getting the Struct Hash

To get the struct hash of typed data:

```go
// Get the struct hash
structHash, err := td.GetStructHash(context.Background(), "Mail", td.Message)
if err != nil {
    panic(err)
}

fmt.Printf("Struct hash: 0x%s\n", structHash.Text(16))
```

### Getting the Domain Hash

To get the domain hash of typed data:

```go
// Get the domain hash
domainHash, err := td.GetDomainHash(context.Background())
if err != nil {
    panic(err)
}

fmt.Printf("Domain hash: 0x%s\n", domainHash.Text(16))
```

## Utilities

The `utils` package provides various cryptographic utilities:

```go
package main

import (
    "fmt"
    "math/big"
    
    "github.com/NethermindEth/starknet.go/utils"
)

func main() {
    // Generate a private key
    privateKey, err := utils.GeneratePrivateKey()
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Private key: 0x%s\n", privateKey.Text(16))
    
    // Get a selector from a function name
    selector := utils.GetSelectorFromName("transfer")
    fmt.Printf("Selector: 0x%s\n", selector)
    
    // Convert a hex string to a big.Int
    hexString := "0x1234567890abcdef"
    bigInt, err := utils.HexToBN(hexString)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Big.Int: %s\n", bigInt.Text(10))
}
```

### Generating Private Keys

To generate a private key:

```go
// Generate a private key
privateKey, err := utils.GeneratePrivateKey()
if err != nil {
    panic(err)
}

fmt.Printf("Private key: 0x%s\n", privateKey.Text(16))
```

### Getting Selectors

To get a selector from a function name:

```go
// Get a selector from a function name
selector := utils.GetSelectorFromName("transfer")
fmt.Printf("Selector: 0x%s\n", selector)
```

### Converting Hex Strings to Big.Int

To convert a hex string to a `big.Int`:

```go
// Convert a hex string to a big.Int
hexString := "0x1234567890abcdef"
bigInt, err := utils.HexToBN(hexString)
if err != nil {
    panic(err)
}

fmt.Printf("Big.Int: %s\n", bigInt.Text(10))
```

### Converting Big.Int to Hex Strings

To convert a `big.Int` to a hex string:

```go
// Convert a big.Int to a hex string
bigInt := big.NewInt(123456789)
hexString := utils.GetHexString(bigInt)

fmt.Printf("Hex string: %s\n", hexString)
```

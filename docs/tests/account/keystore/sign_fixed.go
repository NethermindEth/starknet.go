package main

import (
	"context"
	"fmt"
	"math/big"

	"github.com/NethermindEth/starknet.go/account"
)

func main() {
	// Create a keystore with keys
	ks := account.NewMemKeystore()

	// Use simple test keys
	pubKey1 := "0xabc123def456789"
	privKey1 := new(big.Int).SetUint64(123456)
	ks.Put(pubKey1, privKey1)

	fmt.Printf("Stored key pair:\n")
	fmt.Printf("Public Key:  %s\n", pubKey1)
	fmt.Printf("Private Key: %s\n", privKey1)

	// Create a message to sign (as big.Int)
	fmt.Println("\nSigning a simple message:")
	msgHash := new(big.Int).SetUint64(42)
	fmt.Printf("Message hash: %d\n", msgHash)

	// Sign the message
	ctx := context.Background()
	r, s, err := ks.Sign(ctx, pubKey1, msgHash)
	if err != nil {
		fmt.Printf("Error signing: %v\n", err)
		return
	}

	fmt.Printf("Signature R: %s\n", r)
	fmt.Printf("Signature S: %s\n", s)

	// Note: In Starknet, signatures may include a random nonce component,
	// so consecutive signatures might differ even for the same message.
	// This is a security feature, not a bug.
	fmt.Println("\nNote: Starknet signatures may include randomness for security.")
	fmt.Println("Multiple signatures of the same message may differ - this is expected behavior.")

	// Sign a different message
	fmt.Println("\nSigning a different message:")
	msgHash2 := new(big.Int).SetUint64(999)
	fmt.Printf("Message hash 2: %d\n", msgHash2)

	r3, s3, err := ks.Sign(ctx, pubKey1, msgHash2)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Signature R: %s\n", r3)
	fmt.Printf("Signature S: %s\n", s3)

	diffSig := r.String() != r3.String() || s.String() != s3.String()
	fmt.Printf("Different signature for different message: %v\n", diffSig)

	// Try to sign with a non-existent key
	fmt.Println("\nAttempting to sign with non-existent key:")
	_, _, err = ks.Sign(ctx, "0xnonexistent", msgHash)
	if err != nil {
		fmt.Printf("Error (expected): %v\n", err)
	}

	// Sign with another key
	fmt.Println("\nAdding and signing with another key:")
	pubKey2 := "0xfedcba9876543210"
	privKey2 := new(big.Int).SetUint64(654321)
	ks.Put(pubKey2, privKey2)

	fmt.Printf("Added second key with public: %s\n", pubKey2)

	r4, s4, err := ks.Sign(ctx, pubKey2, msgHash)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Signature from key 2 - R: %s\n", r4)
	fmt.Printf("Signature from key 2 - S: %s\n", s4)

	// Verify different keys produce different signatures for same message
	diffKeys := r.String() != r4.String() || s.String() != s4.String()
	fmt.Printf("Different keys produce different signatures: %v\n", diffKeys)

	// Sign a larger message hash
	fmt.Println("\nSigning with large message hash:")
	largeMsgHash := new(big.Int)
	largeMsgHash.SetString("123456789012345678901234567890123456789012345678901234567890", 10)
	fmt.Printf("Large message hash: %s\n", largeMsgHash)

	rLarge, sLarge, err := ks.Sign(ctx, pubKey1, largeMsgHash)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Signature R: %s\n", rLarge)
	fmt.Printf("Signature S: %s\n", sLarge)
	fmt.Println("Successfully signed large message hash")

	// Test with zero message
	fmt.Println("\nSigning zero message:")
	zeroMsg := new(big.Int).SetUint64(0)
	rZero, sZero, err := ks.Sign(ctx, pubKey1, zeroMsg)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Zero message signature R: %s\n", rZero)
	fmt.Printf("Zero message signature S: %s\n", sZero)
	fmt.Println("Successfully signed zero message")
}
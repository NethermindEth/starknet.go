package account

import (
	"context"
	"fmt"
	"math/big"

	"github.com/NethermindEth/starknet.go/curve"
	"github.com/NethermindEth/starknet.go/utils"
)

type Signer struct {
	keystore  *MemKeystore
	publicKey string
}

// NewSigner creates a new Signer instance with the provided private key.
//
// Parameters:
// - privateKey: a string representing the private key (in any format).
//
// Returns:
// - *Signer: a pointer to the created Signer instance.
// - error: an error if the private key is invalid or if the public key cannot be derived.
func NewSigner(privateKey string) (*Signer, error) {
	pkBigInt, ok := new(big.Int).SetString(privateKey, 0)
	if !ok {
		return nil, fmt.Errorf("invalid private key value")
	}

	pubKey, err := getPublicKey(pkBigInt)
	if err != nil {
		return nil, err
	}

	keyStore := SetNewMemKeystore(pubKey, pkBigInt)

	return &Signer{
		keystore:  keyStore,
		publicKey: pubKey,
	}, nil
}

// PublicKey returns the public key associated with the Signer.
//
// Returns:
// - string: the public key.
func (s *Signer) PublicKey() string {
	return s.publicKey
}

// MemKeyStore returns the keystore used by the Signer.
//
// Returns:
// - *MemKeystore: a pointer to the keystore.
func (s *Signer) MemKeyStore() *MemKeystore {
	return s.keystore
}

// Put stores a new private key in the keystore.
//
// Parameters:
// - priv: a string representing the private key in hexadecimal format.
// - existingKeystore: a pointer to an existing MemKeystore to use. If nil, the Signer's keystore is used.
//
// Returns:
// - error: an error if the private key is invalid.
func (s *Signer) Put(priv string, existingKeystore *MemKeystore) error {
	privateKey, ok := new(big.Int).SetString(priv, 16)
	if !ok {
		return fmt.Errorf("invalid private key value")
	}

	var keystoreToUse *MemKeystore
	if existingKeystore != nil {
		keystoreToUse = existingKeystore
	} else {
		keystoreToUse = s.keystore
	}

	keystoreToUse.Put(s.publicKey, privateKey)
	return nil
}

// Sign signs a message hash using the private key stored in the keystore.
//
// Parameters:
// - ctx: the context.Context object for the signing operation.
// - msgHash: the hash of the message to sign.
//
// Returns:
// - *big.Int: the r part of the signature.
// - *big.Int: the s part of the signature.
// - error: an error if the signing process fails.
func (s *Signer) Sign(ctx context.Context, msgHash *big.Int) (*big.Int, *big.Int, error) {
	return s.keystore.Sign(ctx, s.publicKey, msgHash)
}

// getPublicKey derives the public key from the given private key.
//
// Parameters:
// - priv: a big.Int representing the private key.
//
// Returns:
// - string: the derived public key.
// - error: an error if the public key cannot be derived.
func getPublicKey(priv *big.Int) (string, error) {
	pubX, _, err := curve.Curve.PrivateToPoint(priv)
	if err != nil {
		return "", err
	}

	pub := utils.BigIntToFelt(pubX).String()
	return pub, nil
}

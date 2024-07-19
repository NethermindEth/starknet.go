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

func NewSigner(privateKeyHex string, existingKeystore *MemKeystore) (*Signer, error) {
	privateKey, ok := new(big.Int).SetString(privateKeyHex, 16)
	if !ok {
		return nil, fmt.Errorf("invalid private key format")
	}

	pubKey, err := getPublicKey(privateKey)
	if err != nil {
		return nil, err
	}

	var keyStore *MemKeystore
	if existingKeystore != nil {
		keyStore = existingKeystore
	} else {
		keyStore = SetNewMemKeystore(pubKey, privateKey)
	}

	return &Signer{
		keystore:  keyStore,
		publicKey: pubKey,
	}, nil
}

func (s *Signer) PublicKey() string {
	return s.publicKey
}

func (s *Signer) MemKeyStore() *MemKeystore {
	return s.keystore
}

func (s *Signer) Put(priv string) error {
	privateKey, ok := new(big.Int).SetString(priv, 16)
	if !ok {
		return fmt.Errorf("invalid private key format")
	}
	s.keystore.Put(s.publicKey, privateKey)
	return nil
}

func (s *Signer) Sign(ctx context.Context, msgHash *big.Int) (*big.Int, *big.Int, error) {
	return s.keystore.Sign(ctx, s.publicKey, msgHash)
}

func getPublicKey(priv *big.Int) (string, error) {
	pubX, _, err := curve.Curve.PrivateToPoint(priv)
	if err != nil {
		return "", err
	}

	pub := utils.BigIntToFelt(pubX).String()
	return pub, nil
}

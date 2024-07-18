package account

import (
	"context"
	"math/big"

	"github.com/NethermindEth/starknet.go/curve"
	"github.com/NethermindEth/starknet.go/utils"
)

type Signer struct {
	keystore  *MemKeystore
	publicKey string
}

func NewSigner(privateKey *big.Int) (*Signer, error) {
	pubKey, err := getPublicKey(privateKey)
	if err != nil {
		return nil, err
	}

	keyStore := SetNewMemKeystore(pubKey, privateKey)

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

func (s *Signer) Put(priv *big.Int) {
	s.keystore.Put(s.publicKey, priv)
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

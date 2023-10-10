package account

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"sync"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/curve"
	"github.com/NethermindEth/starknet.go/utils"
)

type Keystore interface {
	Sign(ctx context.Context, id string, msgHash *big.Int) (x *big.Int, y *big.Int, err error)
}

// MemKeystore implements the Keystore interface and is intended for example and test code.
type MemKeystore struct {
	mu   sync.Mutex
	keys map[string]*big.Int
}

func NewMemKeystore() *MemKeystore {
	return &MemKeystore{
		keys: make(map[string]*big.Int),
	}
}

func SetNewMemKeystore(pub string, priv *big.Int) *MemKeystore {
	ks := NewMemKeystore()
	ks.Put(pub, priv)
	return ks
}

func (ks *MemKeystore) Put(senderAddress string, k *big.Int) {
	ks.mu.Lock()
	defer ks.mu.Unlock()
	ks.keys[senderAddress] = k
}

var ErrSenderNoExist = errors.New("sender does not exist")

func (ks *MemKeystore) Get(senderAddress string) (*big.Int, error) {
	ks.mu.Lock()
	defer ks.mu.Unlock()
	k, exists := ks.keys[senderAddress]
	if !exists {
		return nil, fmt.Errorf("error getting key for sender %s: %w", senderAddress, ErrSenderNoExist)
	}
	return k, nil
}

func (ks *MemKeystore) Sign(ctx context.Context, id string, msgHash *big.Int) (*big.Int, *big.Int, error) {

	k, err := ks.Get(id)
	if err != nil {
		return nil, nil, err
	}

	return sign(ctx, msgHash, k)
}

// sign illustrates one way to handle context cancellation
func sign(ctx context.Context, msgHash *big.Int, key *big.Int) (x *big.Int, y *big.Int, err error) {

	select {
	case <-ctx.Done():
		x = nil
		y = nil
		err = ctx.Err()

	default:
		x, y, err = curve.Curve.Sign(msgHash, key)
	}
	return x, y, err
}

// GetRandomKeys gets a random set of pub-priv keys. Note: This should be used for testing purposes only, do NOT send real funds to these addresses.
func GetRandomKeys() (*MemKeystore, *felt.Felt, *felt.Felt) {
	// Get random keys
	privateKey, err := curve.Curve.GetRandomPrivateKey()
	if err != nil {
		fmt.Println("can't get random private key:", err)
		os.Exit(1)
	}
	pubX, _, err := curve.Curve.PrivateToPoint(privateKey)
	if err != nil {
		fmt.Println("can't generate public key:", err)
		os.Exit(1)
	}
	privFelt := utils.BigIntToFelt(privateKey)
	pubFelt := utils.BigIntToFelt(pubX)

	// set up keystore
	ks := NewMemKeystore()
	fakePrivKeyBI, ok := new(big.Int).SetString(privFelt.String(), 0)
	if !ok {
		panic("Error setting up account key store")
	}
	ks.Put(pubFelt.String(), fakePrivKeyBI)

	return ks, pubFelt, privFelt
}

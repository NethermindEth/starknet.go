package caigo

import (
	"errors"
	"fmt"
	"math/big"
	"sync"
)

type Keystore interface {
	Sign(senderAddress string, msgHash *big.Int) (x *big.Int, y *big.Int, err error)
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

func (ks *MemKeystore) Sign(id string, msgHash *big.Int) (*big.Int, *big.Int, error) {

	k, err := ks.Get(id)
	if err != nil {
		return nil, nil, err
	}

	return Curve.Sign(msgHash, k)
}

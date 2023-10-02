package starknetgo

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
)

type Keystore interface {
	Sign(ctx context.Context, id string, msgHash *big.Int) (x *big.Int, y *big.Int, err error)
}

// MemKeystore implements the Keystore interface and is intended for example and test code.
type MemKeystore struct {
	mu   sync.Mutex
	keys map[string]*big.Int
}

// NewMemKeystore creates a new instance of the MemKeystore struct.
//
// It initializes the keys field as an empty map[string]*big.Int object.
// It returns a pointer to the newly created MemKeystore.
func NewMemKeystore() *MemKeystore {
	return &MemKeystore{
		keys: make(map[string]*big.Int),
	}
}

// Put adds the given key to the keystore for the specified sender address.
//
// Parameters:
// - senderAddress: the address of the sender.
// - k: the key to be added.
func (ks *MemKeystore) Put(senderAddress string, k *big.Int) {
	ks.mu.Lock()
	defer ks.mu.Unlock()
	ks.keys[senderAddress] = k
}

var ErrSenderNoExist = errors.New("sender does not exist")

// Get retrieves the value associated with a sender address from the MemKeystore.
//
// It takes a parameter:
// - senderAddress: a string representing the sender address.
//
// It returns:
// - *big.Int: the value associated with the sender address.
// - error: an error if the sender address does not exist in the keystore.
func (ks *MemKeystore) Get(senderAddress string) (*big.Int, error) {
	ks.mu.Lock()
	defer ks.mu.Unlock()
	k, exists := ks.keys[senderAddress]
	if !exists {
		return nil, fmt.Errorf("error getting key for sender %s: %w", senderAddress, ErrSenderNoExist)
	}
	return k, nil
}

// Sign signs a message hash using the given identifier.
//
// Parameters:
// - ctx: The context.Context for the function execution.
// - id: The identifier used to retrieve the key.
// - msgHash: The message hash to be signed.
//
// Returns:
// - The r and s values of the signature.
// - An error if there was a problem signing the message hash.
func (ks *MemKeystore) Sign(ctx context.Context, id string, msgHash *big.Int) (*big.Int, *big.Int, error) {

	k, err := ks.Get(id)
	if err != nil {
		return nil, nil, err
	}

	return sign(ctx, msgHash, k)
}


// sign generates a digital signature for a given message hash using a given private key, illustrates one way to handle context cancellation
//
// Parameters:
// - ctx: The context.Context object for cancellation and timeout functionality.
// - msgHash: The message hash represented as a *big.Int.
// - key: The private key represented as a *big.Int.
//
// Returns:
// - x: The x-coordinate of the generated signature point represented as a *big.Int.
// - y: The y-coordinate of the generated signature point represented as a *big.Int.
// - err: An error object indicating any error that occurred during the signature generation process.
func sign(ctx context.Context, msgHash *big.Int, key *big.Int) (x *big.Int, y *big.Int, err error) {

	select {
	case <-ctx.Done():
		x = nil
		y = nil
		err = ctx.Err()

	default:
		x, y, err = Curve.Sign(msgHash, key)
	}
	return x, y, err
}

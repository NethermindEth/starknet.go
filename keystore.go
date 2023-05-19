package caigo

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"sync"
)

type Keystore interface {
	Sign(senderAddress string, hash []byte) ([]byte, error)
}

const (
	maxPointByteLen = 32 // stark curve max is 252 bits
	signatureLen    = 2 * maxPointByteLen
)

// note: we aren't extending the types.Signature because it looks super brittle. defined as a slice,
// but all uses seems to have an expected order, and JSON handling explicitly restricts to two element.
// would need to add tests before touching it
type signature struct {
	x, y *big.Int
}

func (s *signature) bytes() ([]byte, error) {
	buf := new(bytes.Buffer)
	n, err := buf.Write(padBytes(s.x.Bytes(), maxPointByteLen))
	if err != nil {
		return nil, fmt.Errorf("error writing 'x' component of signature: %w", err)
	}
	if n != maxPointByteLen {
		return nil, fmt.Errorf("unexpected write length of 'x' component of signature: wrote %d expected %d", n, maxPointByteLen)
	}

	n, err = buf.Write(padBytes(s.y.Bytes(), maxPointByteLen))
	if err != nil {
		return nil, fmt.Errorf("error writing 'y' component of signature: %w", err)
	}
	if n != maxPointByteLen {
		return nil, fmt.Errorf("unexpected write length of 'y' component of signature: wrote %d expected %d", n, maxPointByteLen)
	}

	if buf.Len() != signatureLen {
		return nil, fmt.Errorf("error in signature length")
	}
	return buf.Bytes(), nil
}

// b is expected to encode x,y components in accordance with [signature.bytes]
func signatureFromBytes(b []byte) (*signature, error) {
	if len(b) != signatureLen {
		return nil, fmt.Errorf("expected slice len %d got %d", signatureLen, len(b))
	}
	x := b[:maxPointByteLen]
	y := b[maxPointByteLen:]

	return &signature{
		x: new(big.Int).SetBytes(x),
		y: new(big.Int).SetBytes(y),
	}, nil
}

// pad bytes  to specific length
func padBytes(a []byte, length int) []byte {
	if len(a) < length {
		pad := make([]byte, length-len(a))
		return append(pad, a...)
	}

	// return original if length is >= to specified length
	return a
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

// Sign implements the Keystore interface
// this implementation wraps starknet specific curve and expects
// hash: byte representation (big-endian) of *big.Int
func (ks *MemKeystore) Sign(id string, hash []byte) ([]byte, error) {
	k, err := ks.Get(id)
	if err != nil {
		return nil, err
	}

	starkHash := new(big.Int).SetBytes(hash)

	x, y, err := Curve.Sign(starkHash, k)
	if err != nil {
		return nil, fmt.Errorf("error signing data with curve: %w", err)
	}

	s := &signature{
		x: x,
		y: y,
	}
	return s.bytes()
}

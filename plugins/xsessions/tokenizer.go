package xsessions

import (
	"fmt"
	"math/big"
	"time"

	starknetgo "github.com/NethermindEth/starknet.go"
	ctypes "github.com/NethermindEth/starknet.go/types"
)

type Session struct {
	Key      string   `json:"key"`
	Expires  *big.Int `json:"expires"`
	Policies []Policy `json:"policies"`
}

type SignedSession struct {
	ChainID        string     `json:"cahinId"`
	AccountAddress string     `json:"account"`
	Root           string     `json:"root"`
	Signature      []*big.Int `json:"Signature"`
}

type SessionKeyToken struct {
	session       Session
	signedSession SignedSession
}

// computeSessionHash computes the session hash for a given session key, expiration date, root, chain ID, and account address.
//
// sessionKey: The session key.
// expires: The expiration date.
// root: The root.
// chainId: The chain ID.
// accountAddress: The account address.
// Returns the session hash as a *big.Int and an error if there was an error computing the hash.
// TODO remove use of `HexToBN`
func computeSessionHash(sessionKey, expires, root, chainId, accountAddress string) (*big.Int, error) {
	hashDomain, err := starknetgo.Curve.ComputeHashOnElements([]*big.Int{
		STARKNET_DOMAIN_TYPE_HASH,
		ctypes.HexToBN(chainId),
	})
	if err != nil {
		return nil, err
	}
	hashMessage, err := starknetgo.Curve.ComputeHashOnElements([]*big.Int{
		SESSION_TYPE_HASH,
		ctypes.HexToBN(sessionKey),
		ctypes.HexToBN(expires),
		ctypes.HexToBN(root),
	})
	if err != nil {
		return nil, err
	}
	return starknetgo.Curve.ComputeHashOnElements([]*big.Int{
		STARKNET_MESSAGE,
		hashDomain,
		ctypes.HexToBN(accountAddress),
		hashMessage,
	})
}

// getMerkleRoot calculates the Merkle root for the given policies.
//
// policies: A slice of Policy objects representing the policies.
// Returns the Merkle root as a string and an error if any.
func getMerkleRoot(policies []Policy) (string, error) {
	leaves := []*big.Int{}
	for _, policy := range policies {
		leave, err := starknetgo.Curve.ComputeHashOnElements([]*big.Int{
			POLICY_TYPE_HASH,
			ctypes.HexToBN(policy.ContractAddress),
			ctypes.GetSelectorFromName(policy.Selector),
		})
		if err != nil {
			return "", err
		}
		leaves = append(leaves, leave)
	}
	tree, err := starknetgo.NewFixedSizeMerkleTree(leaves...)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("0x%s", tree.Root.Text(16)), nil
}

// SignToken signs a token using the provided parameters.
//
// Parameters:
// - privateKey: the private key used for signing the token.
// - chainId: the chain ID of the token.
// - sessionPublicKey: the public key of the session.
// - accountAddress: the address of the account.
// - duration: the duration for which the token is valid.
// - policies: the policies associated with the token.
//
// Returns:
// - (*SessionKeyToken): the signed session key token.
// - (error): an error if the signing process fails.
func SignToken(privateKey, chainId, sessionPublicKey, accountAddress string, duration time.Duration, policies []Policy) (*SessionKeyToken, error) {
	root, err := getMerkleRoot(policies)
	if err != nil {
		return nil, err
	}
	expires := big.NewInt(time.Now().Add(duration).Unix())
	res, err := computeSessionHash(
		sessionPublicKey,
		fmt.Sprintf("0x%s", expires.Text(16)),
		root,
		chainId,
		accountAddress,
	)
	if err != nil {
		return nil, err
	}
	x, y, err := starknetgo.Curve.Sign(res, ctypes.HexToBN(privateKey))
	if err != nil {
		return nil, err
	}
	return &SessionKeyToken{
		session: Session{
			Key:      sessionPublicKey,
			Expires:  expires,
			Policies: policies,
		},
		signedSession: SignedSession{
			ChainID:        chainId,
			AccountAddress: accountAddress,
			Root:           root,
			Signature:      []*big.Int{x, y},
		},
	}, nil
}

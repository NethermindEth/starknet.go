package xsessions

import (
	"fmt"
	"math/big"
	"time"

	"github.com/dontpanicdao/caigo"
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

// TODO remove use of `HexToBN`
func computeSessionHash(sessionKey, expires, root, chainId, accountAddress string) (*big.Int, error) {
	hashDomain, err := caigo.Curve.ComputeHashOnElements([]*big.Int{
		STARKNET_DOMAIN_TYPE_HASH,
		caigo.HexToBN(chainId),
	})
	if err != nil {
		return nil, err
	}
	hashMessage, err := caigo.Curve.ComputeHashOnElements([]*big.Int{
		SESSION_TYPE_HASH,
		caigo.HexToBN(sessionKey),
		caigo.HexToBN(expires),
		caigo.HexToBN(root),
	})
	if err != nil {
		return nil, err
	}
	return caigo.Curve.ComputeHashOnElements([]*big.Int{
		STARKNET_MESSAGE,
		hashDomain,
		caigo.HexToBN(accountAddress),
		hashMessage,
	})
}

func getMerkleRoot(policies []Policy) (string, error) {
	leaves := []*big.Int{}
	for _, policy := range policies {
		leave, err := caigo.Curve.ComputeHashOnElements([]*big.Int{
			POLICY_TYPE_HASH,
			caigo.HexToBN(policy.ContractAddress),
			caigo.GetSelectorFromName(policy.Selector),
		})
		if err != nil {
			return "", err
		}
		leaves = append(leaves, leave)
	}
	tree, err := caigo.NewFixedSizeMerkleTree(leaves...)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("0x%s", tree.Root.Text(16)), nil
}

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
	x, y, err := caigo.Curve.Sign(res, caigo.HexToBN(privateKey))
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

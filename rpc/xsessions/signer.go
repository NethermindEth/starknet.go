package xsessions

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/dontpanicdao/caigo"
	"github.com/dontpanicdao/caigo/rpc"
	"github.com/dontpanicdao/caigo/rpc/types"
)

var (
	_ rpc.AccountPlugin = &SessionKeyPlugin{}
)

type SessionKeyPlugin struct {
	accountAddress types.Hash
	classHash      *big.Int
	private        *big.Int
	token          *SessionKeyToken
}

func WithSessionKeyPlugin(pluginClassHash string, token *SessionKeyToken) rpc.AccountOptionFunc {
	return func(private, address string) (rpc.AccountOption, error) {
		plugin, ok := big.NewInt(0).SetString(pluginClassHash, 0)
		if !ok {
			return rpc.AccountOption{}, errors.New("could not convert plugin class hash")
		}
		pk, ok := big.NewInt(0).SetString(private, 0)
		if !ok {
			return rpc.AccountOption{}, errors.New("could not convert plugin class hash")
		}
		return rpc.AccountOption{
			AccountPlugin: &SessionKeyPlugin{
				accountAddress: types.HexToHash(address),
				classHash:      plugin,
				private:        pk,
				token:          token,
			},
		}, nil
	}
}

// TODO: write get merkle proof
func getMerkleProof(policies []Policy, call types.FunctionCall) ([]string, error) {
	leaves := []*big.Int{}
	for _, policy := range policies {
		leave, err := caigo.Curve.ComputeHashOnElements([]*big.Int{
			POLICY_TYPE_HASH,
			caigo.HexToBN(policy.ContractAddress),
			caigo.GetSelectorFromName(policy.Selector),
		})
		if err != nil {
			return nil, err
		}
		leaves = append(leaves, leave)
	}
	tree, err := caigo.NewFixedSizeMerkleTree(leaves...)
	if err != nil {
		return nil, err
	}
	callkey, err := caigo.Curve.ComputeHashOnElements([]*big.Int{
		POLICY_TYPE_HASH,
		call.ContractAddress.Big(),
		caigo.GetSelectorFromName(call.EntryPointSelector),
	})
	if err != nil {
		return nil, err
	}
	res, err := tree.Proof(callkey)
	if err != nil {
		return nil, err
	}
	output := []string{}
	for _, r := range res {
		output = append(output, fmt.Sprintf("0x%v", r.Text(16)))
	}
	return output, nil
}

func (plugin *SessionKeyPlugin) PluginCall(calls []types.FunctionCall) (types.FunctionCall, error) {
	data := []string{
		fmt.Sprintf("0x%s", plugin.classHash.Text(16)),
		plugin.token.session.Key,
		fmt.Sprintf("0x%s", plugin.token.session.Expires.Text(16)),
		plugin.token.signedSession.Root,
	}

	firstIteration := true
	for _, call := range calls {
		proof, err := getMerkleProof(plugin.token.session.Policies, call)
		if err != nil {
			return types.FunctionCall{}, err
		}
		if firstIteration {
			length := len(proof)
			data = append(data, fmt.Sprintf("0x%s", big.NewInt(int64(length)).Text(16)))
			firstIteration = false
		}
		data = append(data, proof...)
	}

	for _, signature := range plugin.token.signedSession.Signature {
		data = append(data, fmt.Sprintf("0x%s", signature.Text(16)))
	}
	return types.FunctionCall{
		ContractAddress:    plugin.accountAddress,
		EntryPointSelector: "use_plugin",
		CallData:           data,
	}, nil
}

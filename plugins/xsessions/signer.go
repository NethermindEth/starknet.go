package xsessions

import (
	"errors"
	"fmt"
	"math/big"

	starknetgo "github.com/NethermindEth/starknet.go"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/types"
	ctypes "github.com/NethermindEth/starknet.go/types"
	"github.com/NethermindEth/starknet.go/utils"
)

var (
	_ starknetgo.AccountPlugin = &SessionKeyPlugin{}
)

type SessionKeyPlugin struct {
	accountAddress *felt.Felt
	classHash      *big.Int
	token          *SessionKeyToken
}

func WithSessionKeyPlugin(pluginClassHash string, token *SessionKeyToken) starknetgo.AccountOptionFunc {
	return func(unused, address *felt.Felt) (starknetgo.AccountOption, error) {
		plugin, ok := big.NewInt(0).SetString(pluginClassHash, 0)
		if !ok {
			return starknetgo.AccountOption{}, errors.New("could not convert plugin class hash")
		}
		if !ok {
			return starknetgo.AccountOption{}, errors.New("could not convert plugin class hash")
		}
		return starknetgo.AccountOption{
			AccountPlugin: &SessionKeyPlugin{
				accountAddress: address,
				classHash:      plugin,
				token:          token,
			},
		}, nil
	}
}

// TODO: write get merkle proof
func getMerkleProof(policies []Policy, call ctypes.FunctionCall) ([]string, error) {
	leaves := []*big.Int{}
	for _, policy := range policies {
		leave, err := starknetgo.Curve.ComputeHashOnElements([]*big.Int{
			POLICY_TYPE_HASH,
			ctypes.HexToBN(policy.ContractAddress),
			ctypes.GetSelectorFromName(policy.Selector), // should we use felt??
		})
		if err != nil {
			return nil, err
		}
		leaves = append(leaves, leave)
	}
	tree, err := starknetgo.NewFixedSizeMerkleTree(leaves...)
	if err != nil {
		return nil, err
	}

	callkey, err := starknetgo.Curve.ComputeHashOnElements([]*big.Int{
		POLICY_TYPE_HASH,
		call.ContractAddress.BigInt(big.NewInt(0)),
		ctypes.GetSelectorFromName(call.EntryPointSelector.String()),
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

func (plugin *SessionKeyPlugin) PluginCall(calls []ctypes.FunctionCall) (ctypes.FunctionCall, error) {
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
			return ctypes.FunctionCall{}, err
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
	calldataFelt, err := utils.HexArrToFelt(data)
	if err != nil {
		return ctypes.FunctionCall{}, err
	}
	return ctypes.FunctionCall{
		ContractAddress:    plugin.accountAddress,
		EntryPointSelector: types.GetSelectorFromNameFelt("use_plugin"),
		Calldata:           calldataFelt,
	}, nil
}

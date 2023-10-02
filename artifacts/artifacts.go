package artifacts

import (
	_ "embed"
	"errors"
	"fmt"

	"github.com/NethermindEth/starknet.go/types"
)

//go:embed hello_starknet_compiled.casm.json
var ExampleWorldCasm []byte

//go:embed hello_starknet_compiled.sierra.json
var ExampleWorldSierra []byte

//go:embed starknet_hello_world_Balance.casm.json
var HelloWorldCasm []byte

//go:embed starknet_hello_world_Balance.sierra.json
var HelloWorldSierra []byte

//go:embed account.json
var AccountCompiled []byte

//go:embed account_plugin.json
var AccountWithPluginCompiled []byte

//go:embed starksigner.json
var StarkSignerPluginCompiled []byte

//go:embed sessionkey.json
var SessionKeyPluginCompiled []byte

//go:embed proxy.json
var ProxyCompiled []byte

//go:embed counter.json
var CounterCompiled []byte

//go:embed erc20.json
var ERC20Compiled []byte

type CalldataFormater func(accountHash, pluginHash, publicKey string) ([]string, error)

func AccountV0Formater(accountHash, pluginHash, publicKey string) ([]string, error) {
	calldata := []string{}
	if accountHash == "" {
		calldata = append(calldata, publicKey)
		if pluginHash != "" {
			calldata = append(calldata, pluginHash)
		}
		return calldata, nil
	}
	calldata = append(calldata, accountHash)
	initialize := fmt.Sprintf("0x%x", types.GetSelectorFromName("initialize"))
	calldata = append(calldata, initialize)
	paramLen := "0x1"
	if pluginHash != "" {
		paramLen = "0x2"
	}
	calldata = append(calldata, paramLen)
	calldata = append(calldata, publicKey)
	if pluginHash != "" {
		calldata = append(calldata, pluginHash)
	}
	return calldata, nil
}

func AccountFormater(accountHash, pluginHash, publicKey string) ([]string, error) {
	if pluginHash != "" {
		return nil, errors.New("plugin not supported")
	}
	calldata := []string{}
	if accountHash == "" {
		calldata = append(calldata, publicKey)
		if pluginHash != "" {
			calldata = append(calldata, pluginHash)
		}
		return calldata, nil
	}
	calldata = append(calldata, accountHash)
	initialize := fmt.Sprintf("0x%x", types.GetSelectorFromName("initialize"))

	calldata = append(calldata, initialize)
	calldata = append(calldata, "0x1")
	calldata = append(calldata, publicKey)
	return calldata, nil
}

func AccountPluginFormater(accountHash, pluginHash, publicKey string) ([]string, error) {
	if pluginHash == "" {
		return nil, errors.New("plugin is mandatory")
	}
	calldata := []string{}
	if accountHash != "" {
		calldata = append(calldata, accountHash)
		initialize := fmt.Sprintf("0x%x", types.GetSelectorFromName("initialize"))
		calldata = append(calldata, initialize)
		calldata = append(calldata, "0x4")
	}
	calldata = append(calldata, pluginHash)
	calldata = append(calldata, "0x2")
	calldata = append(calldata, "0x1")
	calldata = append(calldata, publicKey)
	return calldata, nil
}

type CompiledContract struct {
	AccountCompiled []byte
	PluginCompiled  []byte
	ProxyCompiled   []byte
	Formatter       CalldataFormater
}

type CompiledContracts map[string]map[bool]map[bool]CompiledContract

var AccountContracts = CompiledContracts{
	"v1": {
		false: {
			false: {
				AccountCompiled: AccountCompiled,
				Formatter:       AccountFormater,
			},
			true: {
				AccountCompiled: AccountWithPluginCompiled,
				PluginCompiled:  StarkSignerPluginCompiled,
				Formatter:       AccountPluginFormater,
			},
		},
		true: {
			false: {
				AccountCompiled: AccountCompiled,
				ProxyCompiled:   ProxyCompiled,
				Formatter:       AccountFormater,
			},
			true: {
				AccountCompiled: AccountWithPluginCompiled,
				PluginCompiled:  StarkSignerPluginCompiled,
				ProxyCompiled:   ProxyCompiled,
				Formatter:       AccountPluginFormater,
			},
		},
	},
}

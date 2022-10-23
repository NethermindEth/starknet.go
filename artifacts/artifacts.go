package artifacts

import (
	_ "embed"
	"errors"
	"fmt"

	"github.com/dontpanicdao/caigo/types"
)

//go:embed accountv0.json
var AccountV0Compiled []byte

//go:embed accountv0_plugin.json
var AccountV0WithPluginCompiled []byte

//go:embed pluginv0.json
var PluginV0Compiled []byte

//go:embed proxyv0.json
var ProxyV0Compiled []byte

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
	initialize := fmt.Sprintf("0x%s", types.GetSelectorFromName("initialize").Text(16))
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
	initialize := fmt.Sprintf("0x%s", types.GetSelectorFromName("initialize").Text(16))
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
		initialize := fmt.Sprintf("0x%s", types.GetSelectorFromName("initialize").Text(16))
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
	"v0": {
		false: {
			false: {
				AccountCompiled: AccountV0Compiled,
				Formatter:       AccountV0Formater,
			},
			true: {
				AccountCompiled: AccountV0WithPluginCompiled,
				PluginCompiled:  PluginV0Compiled,
				Formatter:       AccountV0Formater,
			},
		},
		true: {
			false: {
				AccountCompiled: AccountV0Compiled,
				ProxyCompiled:   ProxyV0Compiled,
				Formatter:       AccountV0Formater,
			},
			true: {
				AccountCompiled: AccountV0WithPluginCompiled,
				PluginCompiled:  PluginV0Compiled,
				ProxyCompiled:   ProxyV0Compiled,
				Formatter:       AccountV0Formater,
			},
		},
	},
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

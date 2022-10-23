package contracts

import (
	_ "embed"
)

//go:embed artifacts/accountv0.json
var AccountV0Compiled []byte

//go:embed artifacts/accountv0_plugin.json
var AccountV0WithPluginCompiled []byte

//go:embed artifacts/pluginv0.json
var PluginV0Compiled []byte

//go:embed artifacts/proxyv0.json
var ProxyV0Compiled []byte

//go:embed artifacts/account.json
var AccountCompiled []byte

//go:embed artifacts/account_plugin.json
var AccountWithPluginCompiled []byte

//go:embed artifacts/starksigner.json
var StarkSignerPluginCompiled []byte

//go:embed artifacts/sessionkey.json
var SessionKeyPluginCompiled []byte

//go:embed artifacts/proxy.json
var ProxyCompiled []byte

//go:embed artifacts/counter.json
var CounterCompiled []byte

//go:embed artifacts/erc20.json
var ERC20Compiled []byte

type CompiledContract struct {
	AccountCompiled []byte
	PluginCompiled  []byte
	ProxyCompiled   []byte
	Formatter       calldataFormater
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

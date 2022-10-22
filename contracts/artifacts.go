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

//go:embed artifacts/plugin.json
var PluginCompiled []byte

//go:embed artifacts/proxy.json
var ProxyCompiled []byte

//go:embed artifacts/counter.json
var CounterCompiled []byte

//go:embed artifacts/erc20.json
var ERC20Compiled []byte

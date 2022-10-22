
.EXPORT_ALL_VARIABLES:
CAIRO_PATH = ./starknet-plugin-account:./cairo-contracts/src:./src:

.PHONY: build
build: artifacts/account_plugin.json artifacts/plugin.json artifacts/counter.json ozlatest

artifacts/account_plugin.json: starknet-plugin-account/contracts/account/PluginAccount.cairo artifacts
	starknet-compile starknet-plugin-account/contracts/account/PluginAccount.cairo \
		--output artifacts/account_plugin.json \
		--abi artifacts/account_plugin_abi.json --account_contract || \
		rm -f artifacts/account_plugin.json

artifacts/plugin.json: starknet-plugin-account/contracts/plugins/SessionKey.cairo artifacts
	starknet-compile starknet-plugin-account/contracts/plugins/SessionKey.cairo \
		--output artifacts/plugin.json \
		--abi artifacts/plugin_abi.json || \
		rm -f artifacts/plugin.json

artifacts/counter.json: src/counter.cairo artifacts
	starknet-compile src/counter.cairo \
		--output artifacts/counter.json \
		--abi artifacts/counter_abi.json || \
		rm -f artifacts/counter.json

.PHONY: ozlatest
ozlatest: artifacts/account.json artifacts/erc20.json artifacts/proxy.json

artifacts/account.json: src/account.cairo artifacts
	starknet-compile src/account.cairo \
		--output artifacts/account.json \
		--abi artifacts/account_abi.json --account_contract || \
		rm -f artifacts/account.json

artifacts/proxy_plugin.json: src/proxy_plugin.cairo artifacts
	starknet-compile src/proxy_plugin.cairo \
		--output artifacts/proxy_plugin.json \
		--abi artifacts/proxy_plugin_abi.json || \
		rm -f artifacts/proxy_plugin.json

artifacts/proxy.json: cairo-contracts/src/openzeppelin/upgrades/presets/Proxy.cairo artifacts
	starknet-compile cairo-contracts/src/openzeppelin/upgrades/presets/Proxy.cairo \
		--output artifacts/proxy.json \
		--abi artifacts/proxy_abi.json || \
		rm -f artifacts/proxy.json

artifacts/erc20.json: cairo-contracts/src/openzeppelin/token/erc20/presets/ERC20.cairo artifacts
	starknet-compile cairo-contracts/src/openzeppelin/token/erc20/presets/ERC20.cairo \
		--output artifacts/erc20.json \
		--abi artifacts/erc20_abi.json || \
		rm -f artifacts/erc20.json

.PHONY: latest
latest: 
	cd cairo-contracts && git pull
	cd starknet-plugin-account && git pull

artifacts:
	mkdir -p artifacts

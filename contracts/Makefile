
.EXPORT_ALL_VARIABLES:
CAIRO_PATH = ./starknet-plugin-account:./cairo-contracts/src:./src:

.PHONY: build
build: ../artifacts/account_plugin.json ../artifacts/sessionkey.json ../artifacts/starksigner.json ../artifacts/counter.json ozlatest

# modified version of starknet-plugin-account/contracts/account/PluginAccount.cairo
../artifacts/account_plugin.json: src/account_plugin.cairo ../artifacts
	starknet-compile-deprecated src/account_plugin.cairo \
		--output ../artifacts/account_plugin.json \
		--abi ../artifacts/account_plugin_abi.json --account_contract || \
		rm -f ../artifacts/account_plugin.json

../artifacts/sessionkey.json: starknet-plugin-account/contracts/plugins/SessionKey.cairo ../artifacts
	starknet-compile-deprecated starknet-plugin-account/contracts/plugins/SessionKey.cairo \
		--output ../artifacts/sessionkey.json \
		--abi ../artifacts/sessionkey_abi.json || \
		rm -f ../artifacts/sessionkey.json

../artifacts/starksigner.json: starknet-plugin-account/contracts/plugins/signer/StarkSigner.cairo ../artifacts
	starknet-compile-deprecated starknet-plugin-account/contracts/plugins/signer/StarkSigner.cairo \
		--output ../artifacts/starksigner.json \
		--abi ../artifacts/starksigner_abi.json || \
		rm -f ../artifacts/starksigner.json

../artifacts/counter.json: src/counter.cairo ../artifacts
	starknet-compile-deprecated src/counter.cairo \
		--output ../artifacts/counter.json \
		--abi ../artifacts/counter_abi.json || \
		rm -f ../artifacts/counter.json

.PHONY: ozlatest
ozlatest: ../artifacts/account.json ../artifacts/erc20.json ../artifacts/proxy.json

../artifacts/account.json: src/account.cairo ../artifacts
	starknet-compile-deprecated src/account.cairo \
		--output ../artifacts/account.json \
		--abi ../artifacts/account_abi.json --account_contract || \
		rm -f ../artifacts/account.json

../artifacts/proxy_plugin.json: src/proxy_plugin.cairo ../artifacts
	starknet-compile-deprecated src/proxy_plugin.cairo \
		--output ../artifacts/proxy_plugin.json \
		--abi ../artifacts/proxy_plugin_abi.json || \
		rm -f ../artifacts/proxy_plugin.json

../artifacts/proxy.json: cairo-contracts/src/openzeppelin/upgrades/presets/Proxy.cairo ../artifacts
	starknet-compile-deprecated cairo-contracts/src/openzeppelin/upgrades/presets/Proxy.cairo \
		--output ../artifacts/proxy.json \
		--abi ../artifacts/proxy_abi.json || \
		rm -f ../artifacts/proxy.json

../artifacts/erc20.json: cairo-contracts/src/openzeppelin/token/erc20/presets/ERC20.cairo ../artifacts
	starknet-compile-deprecated cairo-contracts/src/openzeppelin/token/erc20/presets/ERC20.cairo \
		--output ../artifacts/erc20.json \
		--abi ../artifacts/erc20_abi.json || \
		rm -f ../artifacts/erc20.json

.PHONY: latest
latest: 
	cd cairo-contracts && git pull
	cd starknet-plugin-account && git pull

../artifacts:
	mkdir -p ../artifacts

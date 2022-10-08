
.EXPORT_ALL_VARIABLES:
CAIRO_PATH = ./cairo-contracts/src:./src

.PHONY: build
build: syncoz0 artifacts/account_plugin.json artifacts/counter.json artifacts/yeasayer.json artifacts/proxy.json artifacts/plugin.json

artifacts/account_plugin.json: src/account_plugin.cairo artifacts
	starknet-compile src/account_plugin.cairo \
		--output artifacts/account_plugin.json \
		--abi artifacts/account_plugin_abi.json --account_contract || \
		rm -f artifacts/account_plugin.json

artifacts/counter.json: src/counter.cairo artifacts
	starknet-compile src/counter.cairo \
		--output artifacts/counter.json \
		--abi artifacts/counter_abi.json || \
		rm -f artifacts/counter.json

artifacts/yeasayer.json: src/yeasayer.cairo artifacts
	starknet-compile src/yeasayer.cairo \
		--output artifacts/yeasayer.json \
		--abi artifacts/yeasayer_abi.json || \
		rm -f artifacts/yeasayer.json

artifacts/plugin.json: src/plugin.cairo artifacts
	starknet-compile src/plugin.cairo \
		--output artifacts/plugin.json \
		--abi artifacts/plugin_abi.json || \
		rm -f artifacts/plugin.json

artifacts/proxy.json: artifacts
	starknet-compile src/proxy.cairo \
		--output artifacts/proxy.json \
		--abi artifacts/proxy_abi.json || \
		rm -f artifacts/proxy.json

.PHONY: ozv0
ozv0: artifacts/account_v0.3.2.json artifacts/erc20.json

artifacts/account_v0.3.2.json: artifacts syncoz0
	cd cairo-contracts/src && \
		starknet-compile ./openzeppelin/account/presets/Account.cairo \
			--output ../../artifacts/account_v0.3.2.json \
			--abi ../../artifacts/account_v0.3.2_abi.json --account_contract || \
			rm -f ../../artifactsaccount_v0.3.2.json

artifacts/erc20.json: artifacts syncoz0
	cd cairo-contracts/src && \
		starknet-compile ./openzeppelin/token/erc20/presets/ERC20.cairo \
			--output ../../artifacts/erc20.json \
			--abi ../../artifacts/erc20_abi.json || \
			rm -f ../../artifacts/erc20.json

.PHONY: ozlatest
ozlatest: artifacts/account_v0.4.0b.json

artifacts/account_v0.4.0b.json: artifacts syncozlatest
	cd cairo-contracts/src && \
		starknet-compile ./openzeppelin/account/presets/Account.cairo \
			--output ../../artifacts/account_v0.4.0b.json \
			--abi ../../artifacts/account_v0.4.0b_abi.json --account_contract || \
			rm -f ../../artifacts/account_v0.4.0b.json

.PHONY: syncoz0
syncoz0:
	cd cairo-contracts && git pull && \
		git reset --hard 789ea11085a178960ce62a774157b6a4f0b85d7f

.PHONY: syncozlatest
syncozlatest: 
	cd cairo-contracts && git pull

artifacts:
	mkdir -p artifacts

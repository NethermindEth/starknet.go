
.EXPORT_ALL_VARIABLES:
CAIRO_PATH = ./cairo-contracts/src:./src

.PHONY: build
build: syncoz0 artifacts/accountv0_plugin.json artifacts/counter.json artifacts/proxy.json artifacts/plugin.json artifacts/accountv0.json

artifacts/accountv0_plugin.json: src/accountv0_plugin.cairo artifacts
	starknet-compile src/accountv0_plugin.cairo \
		--output artifacts/accountv0_plugin.json \
		--abi artifacts/accountv0_plugin_abi.json --account_contract || \
		rm -f artifacts/accountv0_plugin.json

artifacts/counter.json: src/counter.cairo artifacts
	starknet-compile src/counter.cairo \
		--output artifacts/counter.json \
		--abi artifacts/counter_abi.json || \
		rm -f artifacts/counter.json

artifacts/plugin.json: src/yeasayer.cairo artifacts
	starknet-compile src/yeasayer.cairo \
		--output artifacts/plugin.json \
		--abi artifacts/plugin_abi.json || \
		rm -f artifacts/plugin.json

artifacts/proxy.json: artifacts
	starknet-compile src/proxy.cairo \
		--output artifacts/proxy.json \
		--abi artifacts/proxy_abi.json || \
		rm -f artifacts/proxy.json

artifacts/accountv0.json: artifacts syncoz0
		starknet-compile src/accountv0.cairo \
			--output artifacts/accountv0.json \
			--abi artifacts/accountv0_abi.json --account_contract || \
			rm -f artifacts/accountv0.json

artifacts/erc20.json: artifacts syncozlatest
	cd cairo-contracts/src && \
		starknet-compile ./openzeppelin/token/erc20/presets/ERC20.cairo \
			--output ../../artifacts/erc20.json \
			--abi ../../artifacts/erc20_abi.json || \
			rm -f ../../artifacts/erc20.json

.PHONY: syncoz0
syncoz0:
	cd cairo-contracts && git pull && \
		git reset --hard 789ea11085a178960ce62a774157b6a4f0b85d7f

.PHONY: ozlatest
ozlatest: syncozlatest artifacts/accountv1.json artifacts/erc20.json

artifacts/accountv1.json: syncozlatest src/accountv1.cairo artifacts
	starknet-compile src/accountv1.cairo \
		--output artifacts/accountv1.json \
		--abi artifacts/accountv1_abi.json --account_contract || \
		rm -f artifacts/accountv1.json


.PHONY: syncozlatest
syncozlatest: 
	cd cairo-contracts && git pull

artifacts:
	mkdir -p artifacts

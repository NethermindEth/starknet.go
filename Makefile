test:
	@go test ./... -v

rpc-test:
	@go test -v ./rpc -env [mainnet|devnet|testnet|mock]

bench:
	@go test -bench=.

install-deps: | install-gofumpt install-mockgen install-golangci-lint

install-gofumpt:
	go install mvdan.cc/gofumpt@latest

install-mockgen:
	go install go.uber.org/mock/mockgen@latest

install-golangci-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.1

lint:
	@which golangci-lint || make install-golangci-lint
	golangci-lint run

tidy:
	 go mod tidy

format:
	gofumpt -l -w .

simple-call:
	@if [ ! -f ./examples/.env ]; then \
	    echo "This example calls two contract functions, with and without calldata. It uses an ERC20 token, but it can be any smart contract.\n"; \
	    echo "Steps:\n"; \
	    echo "	- Rename the '.env.template' file located at the root of the 'examples' folder to '.env'"; \
	    echo "	- Uncomment, and assign your Sepolia testnet endpoint to the RPC_PROVIDER_URL variable in the '.env' file"; \
	    echo "	- Uncomment, and assign your account address to the ACCOUNT_ADDRESS variable in the '.env' file"; \
	    echo "	- Execute <make simple-call>"; \
	    echo "The calls outputs will be returned at the end of the execution."; \
	else \
	    go run ./examples/simpleCall/main.go; \
	fi

deploy-account:
	@if [ ! -f ./examples/.env ]; then \
	    echo "This example uses a pre-existing class on the Sepolia network to deploy a new account contract. To successfully run this example, you will need: 1) a Sepolia endpoint, and 2) some Sepolia ETH to fund the precomputed address.\n"; \
	    echo "Steps:\n"; \
	    echo "	- Rename the '.env.template' file located at the root of the 'examples' folder to '.env'"; \
	    echo "	- Uncomment, and assign your Sepolia testnet endpoint to the RPC_PROVIDER_URL variable in the '.env' file"; \
	    echo "	- Execute <make deploy-account>"; \
	    echo "	- Fund the precomputed address using a starknet faucet, eg https://starknet-faucet.vercel.app/"; \
	    echo "	- Press any key, then enter"; \
	    echo "At this point your account should be deployed on testnet, and you can use a block explorer like Voyager to view your transaction using the transaction hash."; \
	else \
	    go run ./examples/deployAccount/main.go; \
	fi

simple-invoke:
	@if [ ! -f ./examples/.env ]; then \
	    echo "This example sends an invoke transaction with calldata. It uses an ERC20 token, but it can be any smart contract.\n"; \
	    echo "Steps:\n"; \
	    echo "	- Rename the '.env.template' file located at the root of the 'examples' folder to '.env'"; \
	    echo "	- Uncomment, and assign your Sepolia testnet endpoint to the RPC_PROVIDER_URL variable in the '.env' file"; \
	    echo "	- Uncomment, and assign your account address to the ACCOUNT_ADDRESS variable in the '.env' file (make sure to have a few ETH in it)"; \
	    echo "	- Uncomment, and assign your starknet public key to the PUBLIC_KEY variable in the '.env' file"; \
	    echo "	- Uncomment, and assign your private key to the PRIVATE_KEY variable in the '.env' file"; \
	    echo "	- Execute <make simple-invoke>; \
	    echo "The transaction hash and status will be returned at the end of the execution."; \
	else \
	    go run ./examples/simpleInvoke/main.go; \
	fi

deploy-contractUDC:
	@if [ ! -f ./examples/.env ]; then \
	    echo "This example deploys an ERC20 token using the UDC (Universal Deployer Contract) smart contract.\n"; \
	    echo "Steps:\n"; \
	    echo "	- Rename the '.env.template' file located at the root of the 'examples' folder to '.env'"; \
	    echo "	- Uncomment, and assign your Sepolia testnet endpoint to the RPC_PROVIDER_URL variable in the '.env' file"; \
	    echo "	- Uncomment, and assign your account address to the ACCOUNT_ADDRESS variable in the '.env' file (make sure to have a few ETH in it)"; \
	    echo "	- Uncomment, and assign your starknet public key to the PUBLIC_KEY variable in the '.env' file"; \
	    echo "	- Uncomment, and assign your private key to the PRIVATE_KEY variable in the '.env' file"; \
	    echo "	- Execute <make deploy-contractUDC>"; \
	    echo "The transaction hash and status will be returned at the end of the execution."; \
	else \
	    go run ./examples/deployContractUDC/main.go; \
	fi

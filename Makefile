GOLANGCI_LINT_VERSION := v1.61.0
MOCKGEN_VERSION := v0.5.0
GOFUMPT_VERSION := v0.8.0
GCI_VERSION := v0.13.6
.PHONY: test lint format

# You should first check the 'internal/.env.template' file to set the correct values for the variables.
# tip: use the '-j' flag to run tests in parallel. Example: 'make -j test'
test: clean-testcache mock-test devnet-test testnet-test mainnet-test ## Run all tests
spinup-test: clean-testcache mock-test spinup-devnet-test testnet-test mainnet-test ## Run all tests, but spin up devnet automatically (requires 'starknet-devnet' to be installed)

# small helpers to run 'rpc' and 'account' tests on a specified environment
test_rpc_on = go test -v ./rpc/... -env $(1)
test_account_on = go test -v ./account/... -env $(1)

clean-testcache: ## Clean Go test cache
	@go clean -testcache

mock-test: ## Run all mock tests
	@go test -v ./...

testnet-test: testnet-test-rpc testnet-test-account ## Run 'rpc' and 'account' tests on testnet environment
# splitted to best run in parallel
testnet-test-rpc:
	$(call test_rpc_on,testnet)
testnet-test-account:
	$(call test_account_on,testnet)

mainnet-test: mainnet-test-rpc mainnet-test-account ## Run 'rpc' and 'account' tests on mainnet environment
# splitted to best run in parallel
mainnet-test-rpc:
	$(call test_rpc_on,mainnet)
mainnet-test-account:
	$(call test_account_on,mainnet)

devnet-test: devnet-test-rpc devnet-test-account ## Run 'rpc' and 'account' tests on devnet environment
# splitted to best run in parallel
devnet-test-rpc:
	$(call test_rpc_on,devnet)
devnet-test-account:
	$(call test_account_on,devnet)

spinup-devnet-test: ## Spin up a 'starknet-devnet' instance, run devnet tests, and kill the instance (requires 'starknet-devnet' to be installed)
	@echo "Devnet starting..."
# start the devnet instance and save the pid to a file
	@starknet-devnet > /dev/null & echo $$! > devnet.pid
# wait a few seconds for the devnet instance to start
	@sleep 3
# run the devnet tests and kill the devnet instance after the tests are done.
# if the tests fail, it will remove the devnet.pid file anyhow.
	@$(MAKE) -j devnet-test; status=$$?; \
	if [ -f devnet.pid ]; then kill $$(cat devnet.pid) && rm devnet.pid; else echo "No devnet.pid file found."; fi; \
	exit $$status

lint: ## Run linting
	@echo "Running golangci-lint"
	@golangci-lint run

format: ## Format code
	@gofumpt -l -w .
	@gci write --skip-generated -s standard -s default .

# Install dependencies (Requires go => 1.23)
install-deps: install-gofumpt install-gci install-mockgen install-golangci-lint

install-gofumpt:
	which gofumpt || go install mvdan.cc/gofumpt@$(GOFUMPT_VERSION)

install-gci:
	which gci || go install github.com/daixiang0/gci@$(GCI_VERSION)

install-mockgen:
	which mockgen || go install go.uber.org/mock/mockgen@$(MOCKGEN_VERSION)

install-golangci-lint:
	which golangci-lint || go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

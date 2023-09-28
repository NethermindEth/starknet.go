
This example uses a pre-existing contract on the goerli network to deploy a new account contract. To successfully run this example, you will need: 1) a goerli endpoint, and 2) to fund the precomputed address.

Steps:
1. Rename the ".env.template" file to ".env.testnet"
2. Uncomment, and assign your testnet endpoint to the "INTEGRATION_BASE" variable
3. Execute `go mod tidy` (make sure you are in the "deployAccount" folder)
4. Execute `go run main.go`
5. Fund the precomputed address using a starknet faucet, eg https://faucet.goerli.starknet.io/
6. Press any key, then enter

At this point your account should be deployed on testnet, and you can use a block explorer like [Voyager](https://voyager.online/) to view your transaction using the transaction hash.


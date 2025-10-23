This example demonstrates how to send transactions on Starknet with a paymaster using the Starkent.go [SNIP-29](https://github.com/starknet-io/SNIPs/blob/dfd91b275ea65413f8c8aedb26677a8afff70f37/SNIPS/snip-29.md) implementation, allowing you to pay fees with tokens other than STRK.
It has three files: main.go, deploy.go, and deploy_and_invoke.go.

The main.go file shows how to send an invoke transaction using a paymaster with the "default" fee mode, where you pay fees using supported tokens (like STRK). It demonstrates the complete 3-step process: building the transaction via the paymaster, signing it with your account, and executing the transaction.

The deploy.go file demonstrates how to deploy a new account using a paymaster with the "sponsored" fee mode, where an entity covers the transaction fees.
The deploy_and_invoke.go file shows how to deploy an account and invoke a function in the same transaction using a paymaster, combining both deployment and execution in a single request.
Both of these `deploy...`examples require a valid paymaster API key.

All examples demonstrate integration with the AVNU paymaster service and require SNIP-9 compatible accounts.

Steps:
1. Rename the ".env.template" file located at the root of the "examples" folder to ".env"
2. Uncomment, and assign your Sepolia testnet endpoint to the `RPC_PROVIDER_URL` variable in the ".env" file
3. Uncomment, and assign your SNIP-9 compatible account address to the `ACCOUNT_ADDRESS` variable in the ".env" file (make sure to have some STRK tokens in it)
4. Uncomment, and assign your starknet public key to the `PUBLIC_KEY` variable in the ".env" file
5. Uncomment, and assign your private key to the `PRIVATE_KEY` variable in the ".env" file
6. Make sure you are in the "paymaster" directory
7. Execute `go run .` to run the basic paymaster invoke example
8. To run the deploy examples (requires API key), uncomment the function calls at the end of main.go and execute again. Also, uncomment, and assign your paymaster API key to the `AVNU_API_KEY` variable in the ".env" file

The transaction hashes, tracking IDs, and execution status will be returned at the end of each example.
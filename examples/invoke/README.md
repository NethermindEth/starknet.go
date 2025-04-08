This example sends an invoke transaction with calldata. It uses an ERC20 token, but it can be any smart contract.
It has two files: simpleInvoke.go and verboseInvoke.go.

The simpleInvoke.go file demonstrates a simplified approach to sending invoke transactions using the starknet.go library. It provides a straightforward function that handles building the transaction, estimating fees, signing, and waiting for the transaction receipt in a few lines of code.

The verboseInvoke.go file shows a more detailed, step-by-step approach to the same process, exposing each individual operation (getting the nonce, building the function call, formatting calldata, estimating fees, signing, and sending the transaction). This verbose approach gives you more control over each step of the transaction process.

Both examples demonstrate how to interact with a smart contract on Starknet by calling its functions with the appropriate parameters.
You can choose to run either example, or you can run both!

Steps:
1. Rename the ".env.template" file located at the root of the "examples" folder to ".env"
2. Uncomment, and assign your Sepolia testnet endpoint to the `RPC_PROVIDER_URL` variable in the ".env" file
3. Uncomment, and assign your account address to the `ACCOUNT_ADDRESS` variable in the ".env" file (make sure to have a few ETH in it)
4. Uncomment, and assign your starknet public key to the `PUBLIC_KEY` variable in the ".env" file
5. Uncomment, and assign your private key to the `PRIVATE_KEY` variable in the ".env" file
6. Make sure you are in the "invoke" directory
7. Execute `go run main.go simpleInvoke.go verboseInvoke.go`

The transactions hashes and status will be returned at the end of the execution.
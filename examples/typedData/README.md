This example shows how to sign and verify a typed data.

Steps:
1. Rename the ".env.template" file located at the root of the "examples" folder to ".env"
1. Uncomment, and assign your Sepolia testnet endpoint to the `RPC_PROVIDER_URL` variable in the ".env" file
1. Uncomment, and assign your account address to the `ACCOUNT_ADDRESS` variable in the ".env" file (make sure to have a few ETH in it)
1. Uncomment, and assign your starknet public key to the `PUBLIC_KEY` variable in the ".env" file
1. Uncomment, and assign your private key to the `PRIVATE_KEY` variable in the ".env" file
1. Make sure you are in the "typedData" directory
1. Execute `go run main.go`

The message hash, signature and the verification result will be printed at the end of the execution.

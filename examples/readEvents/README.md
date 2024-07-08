This example reads events from the logs for a smart contract.

Steps:
1. Rename the ".env.template" file located at the root of the "examples" folder to ".env"
1. Uncomment, and assign your Sepolia testnet endpoint to the `RPC_PROVIDER_URL` variable in the ".env" file
1. Uncomment, and assign your account address to the `ACCOUNT_ADDRESS` variable in the ".env" file
1. Make sure you are in the "readEvents" directory
1. Execute `go run main.go`

The events read will be returned at the end of the execution.

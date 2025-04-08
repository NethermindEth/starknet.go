This example calls two contract functions, with and without calldata. It uses an ERC20 token, but it can be any smart contract.

Steps:
1. Rename the ".env.template" file located at the root of the "examples" folder to ".env"
1. Uncomment, and assign your Sepolia testnet endpoint to the `RPC_PROVIDER_URL` variable in the ".env" file
1. Uncomment, and assign your account address to the `ACCOUNT_ADDRESS` variable in the ".env" file
1. Make sure you are in the "simpleCall" directory
1. Execute `go run main.go`

The calls outuputs will be returned at the end of the execution.
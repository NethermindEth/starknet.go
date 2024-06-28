
This example uses a pre-existing class on the Sepolia network to deploy a new account contract. To successfully run this example, you will need: 1) a Sepolia endpoint, and 2) some Sepolia ETH to fund the precomputed address.

Steps:
1. Rename the ".env.template" file located at the root of the "examples" folder to ".env"
1. Uncomment, and assign your Sepolia testnet endpoint to the `RPC_PROVIDER_URL` variable in the ".env" file
1. Make sure you are in the "deployAccount" directory
1. Execute `go run main.go`
1. Fund the precomputed address using a starknet faucet, eg https://starknet-faucet.vercel.app/
1. Press any key, then enter

At this point your account should be deployed on testnet, and you can use a block explorer like [Voyager](https://sepolia.voyager.online/) to view your transaction using the transaction hash.


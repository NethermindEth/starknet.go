This example demonstrates how to declare a contract on Starknet. It uses a simple Hello Starknet contract, but it can be any smart contract.

Steps:
1. Rename the ".env.template" file located at the root of the "examples" folder to ".env"
1. Uncomment, and assign your Sepolia testnet endpoint to the `RPC_PROVIDER_URL` variable in the ".env" file
1. Uncomment, and assign your account address to the `ACCOUNT_ADDRESS` variable in the ".env" file (make sure to have a few ETH in it)
1. Uncomment, and assign your starknet public key to the `PUBLIC_KEY` variable in the ".env" file
1. Uncomment, and assign your private key to the `PRIVATE_KEY` variable in the ".env" file
1. Make sure you are in the "simpleDeclare" directory
1. Ensure you have the contract files (`HelloStarknet.casm.json` and `HelloStarknet.sierra.json`) in the directory
1. Execute `go run main.go`

NOTE: you need to replace the `HelloStarknet` contract for another one. If not, this example WILL RETURN AN ERROR.
This is expected, since the `HelloStarknet` contract was already declared, and there can be only one contract class
in starknet( ref: https://docs.starknet.io/architecture-and-concepts/smart-contracts/contract-classes/#contract_classes_2).

After successful declaration, the transaction hash, status, and the class hash will be returned at the end of the execution.
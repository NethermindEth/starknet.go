# Deploy example

This directory provides a full example on how to use the gateway API to :
- Deploy an OpenZeppelin account
- Deploy an ERC20 contract
- Mint the ERC20 contract
- Transfer tokens from the deployed account to a third-party account

## Run the program
In a terminal in this directory, enter : `go run main.go`

You can choose to run the program with an instance of the devnet (local Starknet instance) or with the testnet by setting the `env` variable to `dev` for devnet or `testnet` for testnet

## Contracts

All used contracts can be found in `./contracts`

- The account is in `./contracts/account/`
- The erc20 is in `./contracts/erc20/`

You will find for each contract :  the Cairo version, the compiled version and the abi.

For the transfer operation, an account is already deployed on testnet at this address : `0x0024e9f35c5d6a14dcbb3f08be5fb7703e76611767266068c521fe8cba27983c`

Note:  

If you run the program with a devnet instance, you have to deploy an account manually and set the `predeployedContract` value with the deployed account address.

## Providing ethereum to the deployed account

When running the program, you will be prompted to add ethereum to the account.

This step has to be done with the testnet [faucet](https://faucet.goerli.starknet.io/)

Copy to the clipboard the address of the contract printed in the terminal and past it in the faucet. The transaction can take several minutes.

Once the transaction is accepted, go to [voyager](https://goerli.voyager.online/) to search for your contract. You should see that it has a little amount of ethereum.

Those ETHs are used to pay transaction fees.

NOTE: this operation has to be done too for devnet. See the devnet documentation to see the process.

## Useful links
- [voyager](https://goerli.voyager.online/): to explore deployed contracts and transactions
- [starknet faucet](https://faucet.goerli.starknet.io/): to provide ETH to accounts
- [devnet](https://github.com/Shard-Labs/starknet-devnet) : local starknet instance

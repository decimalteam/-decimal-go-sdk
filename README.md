# Decimal Go SDK

For detailed explanation on how things work, checkout the:

- [Decimal SDK docs](https://help.decimalchain.com/api-sdk/).
- [Decimal Console site](https://console.decimalchain.com/).

# Install
```
go get bitbucket.org/decimalteam/decimal-go-sdk
```

# Usage

## I. Actions

### Create wallet
```go
package main

import (
	"fmt"
    wallet "bitbucket.org/decimalteam/decimal-go-sdk/wallet"
)

const (
    testMnemonicWords              = "repair furnace west loud peasant false six hockey poem tube now alien service phone hazard winter favorite away sand fuel describe version tragic vendor"
	testMnemonicPassphrase         = ""
)

func main() {
    // Generate private key (account) by mnemonic words (bip39)
    account, err := wallet.NewAccountFromMnemonicWords(testMnemonicWords, testMnemonicPassphrase)
	if err != nil {
		panic(err)
	}
    // Output: dx12k95ukkqzjhkm9d94866r4d9fwx7tsd82r8pjd
    fmt.Println(account.Address())
}
```

### Bind wallet with API
```go
...

import (
    ...
	decapi "bitbucket.org/decimalteam/decimal-go-sdk/api"
)

const (
    ...
	hostURL = "https://testnet-gate.decimalchain.com/api"
	nodeURL = "http://localhost"
)

// for direct connection to node you need:
// 1. launch node with open tendermint RPC port (default: 26657), see https://help.decimalchain.com/masternode-launch/
// command example: decd start
// 2. launch deccli as rest-service with open REST port (default: 1317)
// command example: deccli rest-server --chain-id=...  --laddr=tcp://localhost:1317  --node tcp://localhost:26657 --trust-node=true --unsafe-cors
// In most cases direct connections responses are less informative agains gateway responses.
// Direct connection is useful for multiple transaction sending.

var directConnection = &decapi.DirectConn{PortRPC: ":26657", PortREST: ":1317"}

func main() {
    ...
	// Create Decimal API instance for gateway
	api := decapi.NewAPI(hostURL, nil)

	// Create Decimal API instance for direct connection to node
	api := decapi.NewAPI(nodeURL, directConnection)

	// Bind api-account
	err = bindWalletWithAPI(account, api)
	if err != nil {
		panic(err)
	}
}

func bindWalletWithAPI(account *wallet.Account, api *decapi.API) error {
	// Request chain ID
	chainID, err := api.ChainID()
	if err != nil {
		return err
	}
	account = account.WithChainID(chainID)

	// Request account number and sequence and update with received values
	an, s, err := api.AccountNumberAndSequence(account.Address())
	if err != nil {
		return err
	}
	account = account.WithAccountNumber(an).WithSequence(s)

	return nil
}
```

### Create transaction
```go
...

import (
    ...
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

const (
    ...
	testCoin               = "tdel"
	testReceiverAddress    = "dx1yzxrvpj807dzs5mapwpu77zuh4669lltjheqvv"
)

func main() {
    ...
	// Pack and sign transaction
	tx, err := createTransaction(
		api,
		account,
		testReceiverAddress,
		testCoin,
		sdk.NewInt(1500000000000000000), // 1.5
	)
	if err != nil {
		panic(err)
	}
    ...
}

func createTransaction(api *decapi.API, account *wallet.Account, recv string, ncoin string, val sdk.Int) (types.StdTx, error) {
	// Prepare message arguments
	sender, err := sdk.AccAddressFromBech32(account.Address())
	if err != nil {
		return types.StdTx{}, err
	}
	receiver, err := sdk.AccAddressFromBech32(recv)
	if err != nil {
		return types.StdTx{}, err
	}

	coin := sdk.NewCoin(ncoin, val)

	// Prepare message
	msg := decapi.NewMsgSendCoin(sender, coin, receiver)

	// Prepare transaction arguments
	msgs := []sdk.Msg{msg}
	feeCoins := sdk.NewCoins(sdk.NewCoin(ncoin, sdk.NewInt(0)))
	memo := ""

	// Create signed transaction
	tx, err := api.NewSignedTransaction(msgs, feeCoins, memo, account)
	if err != nil {
		return types.StdTx{}, err
	}

	return tx, nil
}
```

### Create NFT Transaction
```go
...

const (
    ...
    testNFTTokenId         = "n46sJWaSEgJ0Qyie3pelWci7jCI9mN1Wi0QFujHKSenbDAWuxFOjdCfhQmB02lR2"
)

func main() {
	...

  reserve, ok := sdk.NewIntFromString("100000000000000000000")
	if !ok {
		log.Println("invalid reserve")
	}

	// Pack and sign transaction
	tx, err := createTransactionNFT(
		api,
		account,
		testReceiverAddress,
		testCoin,
		testNFTTokenId,
		"", // denom
		fmt.Sprintf("%s/nfts/%s", hostURL, testNFTTokenId), // tokenURI
		sdk.NewInt(1), // quantity
		reserve, // reserve
		true,          // allowMint
	)
	if err != nil {
		panic(err)
	}
    ...
}

func createTransactionNFT(api *decapi.API, account *wallet.Account, recv, ncoin, tokenid, denom, uri string, quan, res sdk.Int, allow bool) (types.StdTx, error) {
	// Prepare message arguments
	sender, err := sdk.AccAddressFromBech32(account.Address())
	if err != nil {
		return types.StdTx{}, err
	}
	receiver, err := sdk.AccAddressFromBech32(recv)
	if err != nil {
		return types.StdTx{}, err
	}

	// Prepare message
	msg := decapi.NewMsgMintNFT(
		sender,
		receiver,
		tokenid,
		denom,
		uri,
		quan,
		res,
		allow,
	)

	// Prepare transaction arguments
	msgs := []sdk.Msg{msg}
	feeCoins := sdk.NewCoins(sdk.NewCoin(ncoin, sdk.NewInt(0)))
	memo := ""

	// Create signed transaction
	tx, err := api.NewSignedTransaction(msgs, feeCoins, memo, account)
	if err != nil {
		return types.StdTx{}, err
	}

	return tx, nil
}
``` 

## II. Views

### Coins information
```go
...

func main() {
    ...
	// Request information about all coins
	coins, err := api.Coins()
	if err != nil {
		panic(err)
	}
	printAsJSON("Coins response", coins)

	// Request information about coin with specific symbol
	symbol := coins[0].Symbol
	coin, err := api.Coin(symbol)
	if err != nil {
		panic(err)
	}
	printAsJSON(fmt.Sprintf("Coin %s response", symbol), coin)
}

// printAsJSON prints `obj` in JSON format.
func printAsJSON(msg string, obj interface{}) {
	objStr, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s:\n%s\n", msg, objStr)
}
```

### Transaction information
```go
...

const (
	testTxHash = "22EAE3E30713B1CC319FDDFCA0F47E94CC4BB94CC2052EBC1A255B53D27D05B7"
)

func main() {
    ...
	// Request information about transaction with specific hash
	tx, err := api.Transaction(testTxHash)
	if err != nil {
		panic(err)
	}
	printAsJSON("Transaction response", tx)
}

...
```

### Candidates information
```go
...

func main() {
    ...
	// Request information about all candidates
	candidates, err := api.Candidates()
	if err != nil {
		panic(err)
	}
	printAsJSON("Candidates response", candidates)
}

...
```

### Validators information
```go
...

func main() {
    ...
	// Request information about all validators
	validators, err := api.Validators()
	if err != nil {
		panic(err)
	}
	printAsJSON("Validators response", validators)

	// Request information about validator with specific address
	validator, err := api.Validator(validators[0].Address)
	if err != nil {
		panic(err)
	}
	printAsJSON("Validator response", validator)
}

...
```

### Stackes information
```go
...

const (
	testStakesMakerAddress = "dx1dqx544dw3gfc2q2n0yv0ghdsjq79zlaf9uflht"
)

func main() {
    ...
	// Request information about stakes from the account with specific address
	stakes, err := api.Stakes(testStakesMakerAddress)
	if err != nil {
		panic(err)
	}
	printAsJSON("Stakes response", stakes)
}

...
```

### Proposals information
```go
...

func main() {
    ...
	// Request information about all govs
	govs, err := api.Proposals()
	if err != nil {
		panic(err)
	}
	printAsJSON("Proposals transactions response", govs)

	govID := govs[0].ProposalID
	gov, err := api.Proposal(govID)
	if err != nil {
		panic(err)
	}
	printAsJSON(fmt.Sprintf("Proposal with ID = %d response", govID), gov)
}

...
```

### Multisig information
```go
...

func main() {
    ...
	// Request information about multisig wallets containing participant with specific address
	multisigWallets, err := api.MultisigWallets(testMultisigParticipantAddress)
	if err != nil {
		panic(err)
	}
	printAsJSON("Multisig wallets response", multisigWallets)

	// Request information about multisig wallet with specific address
	multisigWallet, err := api.MultisigWallet(multisigWallets[0].Address)
	if err != nil {
		panic(err)
	}
	printAsJSON("Multisig wallet response", multisigWallet)

	// Request information about transactions of multisig wallet with specific address
	multisigTransactions, err := api.MultisigTransactions(multisigWallets[0].Address)
	if err != nil {
		panic(err)
	}
	printAsJSON("Multisig transactions response", multisigTransactions)
}

...
```

This README would normally document whatever steps are necessary to get your application up and running.

### What is this repository for? ###

* Quick summary
* Version
* [Learn Markdown](https://bitbucket.org/tutorials/markdowndemo)

### How do I get set up? ###

* Summary of set up
* Configuration
* Dependencies
* Database configuration
* How to run tests
* Deployment instructions

### Contribution guidelines ###

* Writing tests
* Code review
* Other guidelines

### Who do I talk to? ###

* Repo owner or admin
* Other community or team contact

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	decapi "bitbucket.org/decimalteam/decimal-go-sdk/api"
	"bitbucket.org/decimalteam/decimal-go-sdk/wallet"
)

const (
	hostURL = "https://testnet-gate.decimalchain.com/api"

	testMnemonicWords              = "repair furnace west loud peasant false six hockey poem tube now alien service phone hazard winter favorite away sand fuel describe version tragic vendor"
	testMnemonicPassphrase         = ""
	testSenderAddress              = "dx12k95ukkqzjhkm9d94866r4d9fwx7tsd82r8pjd"
	testReceiverAddress            = "dx1yzxrvpj807dzs5mapwpu77zuh4669lltjheqvv"
	testStakesMakerAddress         = "dx1g0gf9gdcyyrqk9rhnm23l7dsm7v5ex8c3ga98p"
	testValidatorAddress           = "dxvaloper16rr3cvdgj8jsywhx8lfteunn9uz0xg2czw6gx5"
	testMultisigParticipantAddress = "dx173lnn7jjuym5rwp23aufhnwshylrdemcswtcg5"
	testMultisigAddress            = "dx1nkujpc7fj72cfdyrtj7f090wgdakjvnyy6dak5"
	testCoin                       = "tdel"
	testTxHash                     = "F7BEE024F6EECD0909EF90B5C6A46FE6AFD7AEF061EE02BB73F800960EF57326"
	testNFTTokenId                 = "rt7c7255cd002f1595a8d8a00ce11ffce25a315t"

	testWrongSenderAddress              = "dx12k95ukkqzjhkm9d94866r4d9fwx7tsd82r8p00"
	testWrongReceiverAddress            = "dx1yzxrvpj807dzs5mapwpu77zuh4669lltjheq00"
	testWrongStakesMakerAddress         = "dx1dqx544dw3gfc2q2n0yv0ghdsjq79zlaf000000"
	testWrongValidatorAddress           = "dxvaloper16rr3cvdgj8jsywhx8lfteunn9uz0xg2czw6g00"
	testWrongMultisigParticipantAddress = "dx173lnn7jjuym5rwp23aufhnwshylrdemcswtc00"
	testWrongMultisigAddress            = "dx1kgnzuwwgzhecyk0dn62sxmp4wyuk0000000000"
	testWrongCoin                       = "tdel0"
	testWrongTxHash                     = "22EAE3E30713B1CC319FDDFCA0F47E94CC4BB94CC2052EBC1A255B53D27D0500"
	testWrongNFTTokenId                 = "n46sJWaSEgJ0Qyie3pelWci7jCI9mN1Wi0QFujHKSenbDAWuxFOjdCfhQmB02l00"
	testWrongProposalID                 = 100500
)

var api *decapi.API
var account *wallet.Account

var err error

////////////////////////////////////////////////////////////////
// Decimal SDK example initializing
////////////////////////////////////////////////////////////////

func init() {
	rand.Seed(time.Now().UnixNano())

	// Create Decimal API instance
	api = decapi.NewAPI(hostURL)

	// Create account from the mnemonic
	account, err = wallet.NewAccountFromMnemonicWords(testMnemonicWords, testMnemonicPassphrase)
	if err != nil {
		log.Println("Error while creating account")
		panic(err)
	}

	// Request chain ID
	if chainID, err := api.ChainID(); err == nil {
		account = account.WithChainID(chainID)
		log.Printf("Current chain ID: %s\n", chainID)
	} else {
		log.Println("Error while requesting chain id")
		panic(err)
	}

	// Request account number and sequence and update with received values
	if an, s, err := api.AccountNumberAndSequence(testSenderAddress); err == nil {
		account = account.WithAccountNumber(an).WithSequence(s)
		log.Printf("Account %s number: %d, sequence: %d\n", account.Address(), an, s)
	} else {
		log.Println("Error while account number and sequence")
		panic(err)
	}
}

////////////////////////////////////////////////////////////////
// Decimal SDK example running
////////////////////////////////////////////////////////////////

func main() {

	// Request everything from the API
	exampleRequests()
	exampleWrongRequests()

	// Create and broadcast transactions
	exampleBroadcastMsgSendCoin()

	exampleBroadcastMsgMintNFT()
}

////////////////////////////////////////////////////////////////
// Decimal API requesting data
////////////////////////////////////////////////////////////////

func exampleRequests() {

	// Request information about the address
	fmt.Println("Request information about the address")
	address, err := api.Address(testSenderAddress)
	if err != nil {
		panic(err)
	}
	printAsJSON("Address response", address)

	// Request information about all coins
	fmt.Println("Request information about all coins")
	coins, err := api.Coins()
	if err != nil {
		panic(err)
	}
	printAsJSON("Coins response", coins)

	// Request information about coin with specific symbol
	fmt.Println("Request information about coin with specific symbol")
	symbol := coins[0].Symbol
	coin, err := api.Coin(symbol)
	if err != nil {
		panic(err)
	}
	printAsJSON(fmt.Sprintf("Coin %s response", symbol), coin)

	// Request information about transaction with specific hash
	fmt.Println("Request information about transaction with specific hash")
	tx, err := api.Transaction(testTxHash)
	if err != nil {
		panic(err)
	}
	printAsJSON("Transaction response", tx)

	// Request information about all candidates
	fmt.Println("Request information about all candidates")
	candidates, err := api.Candidates()
	if err != nil {
		panic(err)
	}
	printAsJSON("Candidates response", candidates)

	// Request information about all validators
	fmt.Println("Request information about all validators")
	validators, err := api.Validators()
	if err != nil {
		panic(err)
	}
	printAsJSON("Validators response", validators)

	// Request information about validator with specific address
	fmt.Println("Request information about validator with specific address")
	validator, err := api.Validator(testValidatorAddress)
	if err != nil {
		panic(err)
	}
	printAsJSON("Validator response", validator)

	// Request information about stakes from the account with specific address
	fmt.Println("Request information about stakes from the account with specific address")
	stakes, err := api.Stakes(testStakesMakerAddress)
	if err != nil {
		panic(err)
	}
	printAsJSON("Stakes response", stakes)

	// Request information about multisig wallets containing participant with specific address
	fmt.Println("Request information about multisig wallets containing participant with specific address")
	multisigWallets, err := api.MultisigWallets(testMultisigParticipantAddress)
	if err != nil {
		panic(err)
	}
	printAsJSON("Multisig wallets response", multisigWallets)

	// Request information about multisig wallet with specific address
	fmt.Println("Request information about multisig wallet with specific address")
	multisigWallet, err := api.MultisigWallet(testMultisigAddress)
	if err != nil {
		panic(err)
	}
	printAsJSON("Multisig wallet response", multisigWallet)

	// Request information about transactions of multisig wallet with specific address
	fmt.Println("Request information about transactions of multisig wallet with specific address")
	multisigTransactions, err := api.MultisigTransactions(testMultisigAddress)
	if err != nil {
		panic(err)
	}
	printAsJSON("Multisig transactions response", multisigTransactions)

	// Request information about all govs
	fmt.Println("Request information about all govs")
	govs, err := api.Proposals()
	if err != nil {
		panic(err)
		return
	}
	printAsJSON("Proposals transactions response", govs)

	fmt.Println("Get first proposal response")
	govID := govs[0].ProposalID
	gov, err := api.Proposal(govID)
	if err != nil {
		panic(err)
		return
	}
	printAsJSON(fmt.Sprintf("Proposal with ID = %d response", govID), gov)
}

func exampleWrongRequests() {

	// Request information about the address
	fmt.Println("(ERR)Request information about the address")
	_, err := api.Address(testWrongSenderAddress)
	if err == nil {
		panic("You should get error")
	}
	fmt.Println("You got error: ", err)

	// Request information about coin with specific symbol
	fmt.Println("(ERR)Request information about coin with specific symbol")
	_, err = api.Coin(testWrongCoin)
	if err == nil {
		panic("You should get error")
	}
	fmt.Println("You got error: ", err)

	// Request information about transaction with specific hash
	fmt.Println("(ERR)Request information about transaction with specific hash")
	_, err = api.Transaction(testWrongTxHash)
	if err == nil {
		panic("You should get error")
	}
	fmt.Println("You got error: ", err)

	// Request information about validator with specific address
	fmt.Println("(ERR)Request information about validator with specific address")
	_, err = api.Validator(testWrongValidatorAddress)
	if err == nil {
		panic("You should get error")
	}
	fmt.Println("You got error: ", err)

	// Request information about specific proposal
	fmt.Println("(ERR)Request proposal info")
	_, err = api.Proposal(testWrongProposalID)
	if err == nil {
		panic("You should get error")
	}
	fmt.Println("You got error: ", err)
}

////////////////////////////////////////////////////////////////
// Decimal API broadcasting transactions
////////////////////////////////////////////////////////////////

func exampleBroadcastMsgSendCoin() {

	// Prepare message arguments
	sender, err := sdk.AccAddressFromBech32(testSenderAddress)
	if err != nil {
		panic(err)
	}
	receiver, err := sdk.AccAddressFromBech32(testReceiverAddress)
	if err != nil {
		panic(err)
	}
	amount := sdk.NewInt(1500000000000000000) // 1.5
	coin := sdk.NewCoin(testCoin, amount)

	// Prepare message
	msg := decapi.NewMsgSendCoin(sender, coin, receiver)

	// Prepare transaction arguments
	msgs := []sdk.Msg{msg}
	feeCoins := sdk.NewCoins(sdk.NewCoin(testCoin, sdk.NewInt(0)))
	memo := ""

	// Create signed transaction
	tx, err := api.NewSignedTransaction(msgs, feeCoins, memo, account)
	if err != nil {
		panic(err)
	}

	// Broadcast signed transaction
	broadcastTxResult, err := api.BroadcastSignedTransactionJSON(tx, account)
	if err != nil {
		panic(err)
	}
	printAsJSON("Broadcast transaction (in JSON format) response", broadcastTxResult)

	// TODO: Block code executing until the transaction is placed in a block?
}

func exampleBroadcastMsgMintNFT() {
	// Prepare message arguments
	sender, err := sdk.AccAddressFromBech32(testSenderAddress)
	if err != nil {
		panic(err)
	}
	receiver, err := sdk.AccAddressFromBech32(testReceiverAddress)
	if err != nil {
		panic(err)
	}
	amount := sdk.NewInt(1500000000000000000) // 1.5
	_ = sdk.NewCoin(testCoin, amount)

	reserve, ok := sdk.NewIntFromString("100000000000000000000")
	if !ok {
		log.Println("invalid reserve")
	}

	// Prepare message
	msg := decapi.NewMsgMintNFT(
		sender,
		receiver,
		testNFTTokenId,
		"cat", // denom
		fmt.Sprintf("%s/nfts/%s", hostURL, testNFTTokenId), // tokenURI
		sdk.NewInt(1), // quantity
		reserve,       // reserve
		true,          // allowMint
	)

	// Prepare transaction arguments
	msgs := []sdk.Msg{msg}
	feeCoins := sdk.NewCoins(sdk.NewCoin(testCoin, sdk.NewInt(0)))
	memo := ""

	// Create signed transaction
	tx, err := api.NewSignedTransaction(msgs, feeCoins, memo, account)
	if err != nil {
		panic(err)
	}

	// Broadcast signed transaction
	broadcastTxResult, err := api.BroadcastSignedTransactionJSON(tx, account)
	if err != nil {
		panic(err)
	}
	printAsJSON("Broadcast transaction nft/mint_nft response", broadcastTxResult)
}

// printAsJSON prints `obj` in JSON format.
func printAsJSON(msg string, obj interface{}) {
	objStr, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s:\n%s\n", msg, objStr)
}

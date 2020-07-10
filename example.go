package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"

	"github.com/sethvargo/go-diceware/diceware"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"

	decapi "bitbucket.org/decimalteam/decimal-go-sdk/api"
	"bitbucket.org/decimalteam/decimal-go-sdk/wallet"
)

const (
	hostURL = "https://testnet-gate.decimalchain.com/api"

	testMnemonicWords      = "repair furnace west loud peasant false six hockey poem tube now alien service phone hazard winter favorite away sand fuel describe version tragic vendor"
	testMnemonicPassphrase = ""
	testSenderAddress      = "dx12k95ukkqzjhkm9d94866r4d9fwx7tsd82r8pjd"
	testReceiverAddress    = "dx1yzxrvpj807dzs5mapwpu77zuh4669lltjheqvv"
	testValidatorAddress   = "dxvaloper16rr3cvdgj8jsywhx8lfteunn9uz0xg2czw6gx5"
	testCoin               = "tdel"
	testTxHash             = "16B3284DD7ADCCCE74109A713B87D134C92089A0759297551BA2D4B4DA558B40"

	hugeGas = 16 * 1024
)

var api *decapi.API
var account *wallet.Account

var err error

////////////////////////////////////////////////////////////////
// Decimal SDK example initializing
////////////////////////////////////////////////////////////////

func init() {

	// Create Decimal API instance
	api = decapi.NewAPI(hostURL)

	// Create account from the mnemonic
	account, err = wallet.NewAccountFromMnemonicWords(testMnemonicWords, testMnemonicPassphrase)
	if err != nil {
		panic(err)
	}

	// Request chain ID
	if chainID, err := api.ChainID(); err == nil {
		account = account.WithChainID(chainID)
		fmt.Printf("Current chain ID: %s\n", chainID)
	} else {
		panic(err)
	}

	// Request account number and sequence and update with received values
	if an, s, err := api.AccountNumberAndSequence(testSenderAddress); err == nil {
		account = account.WithAccountNumber(an).WithSequence(s)
		fmt.Printf("Account %s number: %d, sequence: %d\n", account.Address(), an, s)
	} else {
		panic(err)
	}
}

////////////////////////////////////////////////////////////////
// Decimal SDK example running
////////////////////////////////////////////////////////////////

func main() {

	// Request everything from the API
	exampleRequests()

	// Create and broadcast transactions
	exampleBroadcastMsgSendCoin()
}

////////////////////////////////////////////////////////////////
// Decimal API requesting data
////////////////////////////////////////////////////////////////

func exampleRequests() {

	// Request information about the address
	address, err := api.Address(testSenderAddress)
	if err != nil {
		panic(err)
	}
	printAsJSON("Address response", address)

	// Request balance of the address
	balance, err := api.Balance(testSenderAddress)
	if err != nil {
		panic(err)
	}
	printAsJSON("Balance response", balance)

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

	// Request information about all candidates
	candidates, err := api.Candidates()
	if err != nil {
		panic(err)
	}
	printAsJSON("Candidates response", candidates)

	// Request information about all validators
	validators, err := api.Validators()
	if err != nil {
		panic(err)
	}
	printAsJSON("Validators response", validators)

	// Request information about validator with specific address
	validator, err := api.Validator(testValidatorAddress)
	if err != nil {
		panic(err)
	}
	printAsJSON("Validator response", validator)

	// Request information about transaction with specific hash
	tx, err := api.Transaction(testTxHash)
	if err != nil {
		panic(err)
	}
	printAsJSON("Transaction response", tx)
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
	fee := auth.NewStdFee(hugeGas, sdk.NewCoins(sdk.NewCoin(testCoin, sdk.NewInt(0))))
	memo := ""
	if words, _ := diceware.Generate(rand.Int() % 16); len(words) > 0 {
		memo = strings.Join(words, " ")
	}

	// Create and sign transaction
	tx := account.CreateTransaction(msgs, fee, memo)
	tx, err = account.SignTransaction(tx)
	if err != nil {
		panic(err)
	}

	// TODO: Estimate and adjust amount of gas wanted for the transaction

	// Broadcast signed transaction
	sendTxResult, err := api.SendTransactionJSON(tx)
	if err != nil {
		panic(err)
	}
	printAsJSON("Sent transaction in JSON format response", sendTxResult)

	// TODO: Block code executing until the transaction is placed in a block?
}

// printAsJSON prints `obj` in JSON format.
func printAsJSON(msg string, obj interface{}) {
	objStr, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s:\n%s\n", msg, objStr)
}

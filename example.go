package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"

	decapi "bitbucket.org/decimalteam/decimal-go-sdk/api"
	"bitbucket.org/decimalteam/decimal-go-sdk/wallet"
)

const (
	hostURL = "https://testnet-gate.decimalchain.com/api"

	testMnemonic         = "repair furnace west loud peasant false six hockey poem tube now alien service phone hazard winter favorite away sand fuel describe version tragic vendor"
	testSenderAddress    = "dx12k95ukkqzjhkm9d94866r4d9fwx7tsd82r8pjd"
	testReceiverAddress  = "dx1yzxrvpj807dzs5mapwpu77zuh4669lltjheqvv"
	testValidatorAddress = "dxvaloper16rr3cvdgj8jsywhx8lfteunn9uz0xg2czw6gx5"
)

func main() {

	// Create Decimal API instance
	api := decapi.NewAPI(hostURL)

	// Request information about the address
	address, err := api.Address(testSenderAddress)
	if err != nil {
		panic(err)
	}
	print("Address request result", address)

	// Request balance of the address
	balance, err := api.Balance(testSenderAddress)
	if err != nil {
		panic(err)
	}
	print("Balance request result", balance)

	// Request information about all coins
	coins, err := api.Coins()
	if err != nil {
		panic(err)
	}
	print("Coins request result", coins)

	// Request information about coin with specific symbol
	symbol := coins[0].Symbol
	coin, err := api.Coin(symbol)
	if err != nil {
		panic(err)
	}
	print(fmt.Sprintf("Coin %s request result", symbol), coin)

	// Request information about validator with specific address
	validator, err := api.Validator(testValidatorAddress)
	if err != nil {
		panic(err)
	}
	print("Validator request result", validator)

	// Create mnemonic from words
	mnemonic, err := wallet.NewMnemonicFromWords(testMnemonic, "")
	if err != nil {
		panic(err)
	}

	// Create extended key from the mnemonic
	extendedKey, err := wallet.NewExtendedKeyFromMnemonic(mnemonic)
	if err != nil {
		panic(err)
	}

	// Derive extended key
	extendedKey, err = extendedKey.GetChildAtPath("44'/60'/0'/0/0")
	if err != nil {
		panic(err)
	}

	// Create private key
	privateKey, err := extendedKey.GetECPrivateKey()
	if err != nil {
		panic(err)
	}

	// Create send coin transaction
	sender, err := sdk.AccAddressFromBech32(testSenderAddress)
	if err != nil {
		panic(err)
	}
	receiver, err := sdk.AccAddressFromBech32(testReceiverAddress)
	if err != nil {
		panic(err)
	}
	msg := decapi.MsgSendCoin{
		Sender:   sender,
		Coin:     sdk.NewCoin("tdel", sdk.NewInt(1000000000000000000)),
		Receiver: receiver,
	}
	fee := auth.NewStdFee(200000, sdk.NewCoins(sdk.NewCoin("tdel", sdk.NewInt(0))))
	memo := "Test sending coin..."
	tx, bytesToSign, err := api.NewTransaction(testSenderAddress, []sdk.Msg{msg}, fee, memo)
	if err != nil {
		panic(err)
	}
	tx, err = api.SignTransaction(tx, bytesToSign, privateKey)
	if err != nil {
		panic(err)
	}
	txBytesBinary, err := api.EncodeTransactionBinary(tx)
	txBytesJSON, err := api.EncodeTransactionJSON(tx)
	fmt.Printf("Binary encoded transaction:\n%s\n", hex.EncodeToString(txBytesBinary))
	fmt.Printf("JSON encoded transaction:\n%s\n", string(txBytesJSON))
}

// print prints `obj` in JSON format.
func print(msg string, obj interface{}) {
	objStr, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s:\n%s\n", msg, objStr)
}

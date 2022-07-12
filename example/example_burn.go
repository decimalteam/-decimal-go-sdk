package main

import (
	"log"
	"time"

	decapi "bitbucket.org/decimalteam/decimal-go-sdk/api"
	"bitbucket.org/decimalteam/decimal-go-sdk/wallet"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func testBurnCoin(api *decapi.API) {
	log.Printf("START test send")
	// make wallets
	mnemonic1 := "plug tissue today frown increase race brown sail post march trick coconut laptop churn call child question match also spend play credit already travel"
	acc1, err := wallet.NewAccountFromMnemonicWords(mnemonic1, "")
	if err != nil {
		log.Printf("ERROR: acc1 %s", err.Error())
	}
	// set chain id
	chainId, _ := api.ChainID()
	acc1 = acc1.WithChainID(chainId)

	bindAcc(api, acc1)
	//prepare transaction
	sender, err := sdk.AccAddressFromBech32(acc1.Address())
	if err != nil {
		log.Printf("ERROR: AccAddressFromBech32 %s->%s", acc1.Address(), err.Error())
	}

	//10^18
	amount := sdk.NewInt(1500000000000000000) // 1.5
	coin := sdk.NewCoin("crash", amount)      // invalid coin

	// Prepare message
	msg := decapi.NewMsgBurnCoin(sender, coin)

	// Prepare transaction arguments
	msgs := []sdk.Msg{msg}
	feeCoins := sdk.NewCoins(sdk.NewCoin(testCoin, sdk.NewInt(0)))
	memo := "test burn"

	// Create signed transaction
	tx, err := api.NewSignedTransaction(msgs, feeCoins, memo, acc1)
	if err != nil {
		log.Printf("ERROR: NewSignedTransaction(from %s) %s", acc1.Address(), err.Error())
	}
	log.Printf("SignedTransaction result: %s", formatAsJSON(tx))

	// Broadcast signed transaction
	broadcastTxResult, err := api.BroadcastSignedTransactionJSON(tx, acc1)
	if err != nil {
		log.Printf("ERROR: BroadcastSignedTransactionJSON(from %s) %s", acc1.Address(), err.Error())
	}
	log.Printf("BroadcastSignedTransactionJSON result: %s", formatAsJSON(broadcastTxResult))
	time.Sleep(time.Second * 5)

	// Get transaction result
	var txRes *decapi.TransactionResult = nil
	seconds := 0
	for txRes == nil { // txRes == nil mean transaction not in block yet
		txRes, err = api.Transaction(broadcastTxResult.TxHash)
		time.Sleep(time.Second)
		seconds++
		if seconds >= 5 { // 5 seconds above + 5 seconds ~ 2 blocks. No need to wait more
			break
		}
	}

	if err != nil {
		log.Printf("Final error: %s", err.Error())
		return
	}

	log.Printf("Tx result: %s", formatAsJSON(txRes))

}

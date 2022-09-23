package main

import (
	"log"
	"time"

	decapi "bitbucket.org/decimalteam/decimal-go-sdk/api"
	"bitbucket.org/decimalteam/decimal-go-sdk/wallet"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Test example to send coins
// It shows set up for wallets(accounts) and preparations for transactions
func testSend(api *decapi.API, baseCoin string) {
	log.Printf("START test send")
	// make wallets
	mnemonic1 := "plug tissue today frown increase race brown sail post march trick coconut laptop churn call child question match also spend play credit already travel"
	acc1, err := wallet.NewAccountFromMnemonicWords(mnemonic1, "")
	if err != nil {
		log.Printf("ERROR: acc1 %s", err.Error())
	}
	mnemonic2 := "layer pass tide basic raccoon olive trust satoshi coil harbor script shrimp health gadget few armed rival spread release welcome long dust almost banana"
	acc2, err := wallet.NewAccountFromMnemonicWords(mnemonic2, "")
	if err != nil {
		log.Printf("ERROR: acc2 %s", err.Error())
	}
	// set chain id
	chainId, _ := api.ChainID()
	acc1 = acc1.WithChainID(chainId)
	acc2 = acc2.WithChainID(chainId)

	//send in both directions
	testCases := []struct {
		accFrom *wallet.Account
		accTo   *wallet.Account
	}{
		{acc1, acc2},
		{acc2, acc1},
	}
	for _, tst := range testCases {
		bindAcc(api, acc1)
		bindAcc(api, acc2)
		//prepare transaction
		sender, err := sdk.AccAddressFromBech32(tst.accFrom.Address())
		if err != nil {
			log.Printf("ERROR: AccAddressFromBech32 %s->%s", tst.accFrom.Address(), err.Error())
		}
		receiver, err := sdk.AccAddressFromBech32(tst.accTo.Address())
		if err != nil {
			log.Printf("ERROR: AccAddressFromBech32 %s->%s", tst.accTo.Address(), err.Error())
		}
		//10^18
		amount := sdk.NewInt(1500000000000000000) // 1.5
		coin := sdk.NewCoin(baseCoin, amount)

		// Prepare message
		msg := decapi.NewMsgSendCoin(sender, coin, receiver)

		// Prepare transaction arguments
		msgs := []sdk.Msg{msg}
		feeCoins := sdk.NewCoins(sdk.NewCoin(testCoin, sdk.NewInt(0)))
		memo := "test message"

		// Create signed transaction
		tx, err := api.NewSignedTransaction(msgs, feeCoins, memo, tst.accFrom)
		if err != nil {
			log.Printf("ERROR: NewSignedTransaction(from %s) %s", tst.accTo.Address(), err.Error())
		}
		log.Printf("SignedTransaction result: %s", formatAsJSON(tx))

		// Broadcast signed transaction
		broadcastTxResult, err := api.BroadcastSignedTransactionJSON(tx, tst.accFrom)
		if err != nil {
			log.Printf("ERROR: BroadcastSignedTransactionJSON(from %s) %s", tst.accTo.Address(), err.Error())
		}
		log.Printf("BroadcastSignedTransactionJSON result: %s", formatAsJSON(broadcastTxResult))
		var txRes *decapi.TransactionResult = nil
		seconds := 0
		for txRes == nil {
			txRes, err = api.Transaction(broadcastTxResult.TxHash)
			time.Sleep(time.Second)
			seconds++
			if seconds > 30 {
				break
			}
		}

		log.Printf("Tx result: %s", formatAsJSON(txRes))
		log.Printf("seconds: %d", seconds)
	}
}

func testInvalidSendCoin(api *decapi.API) {
	log.Printf("START test send")
	// make wallets
	mnemonic1 := "plug tissue today frown increase race brown sail post march trick coconut laptop churn call child question match also spend play credit already travel"
	acc1, err := wallet.NewAccountFromMnemonicWords(mnemonic1, "")
	if err != nil {
		log.Printf("ERROR: acc1 %s", err.Error())
	}
	mnemonic2 := "layer pass tide basic raccoon olive trust satoshi coil harbor script shrimp health gadget few armed rival spread release welcome long dust almost banana"
	acc2, err := wallet.NewAccountFromMnemonicWords(mnemonic2, "")
	if err != nil {
		log.Printf("ERROR: acc2 %s", err.Error())
	}
	// set chain id
	chainId, _ := api.ChainID()
	acc1 = acc1.WithChainID(chainId)
	acc2 = acc2.WithChainID(chainId)

	//send in both directions
	testCases := []struct {
		accFrom *wallet.Account
		accTo   *wallet.Account
	}{
		{acc1, acc2},
		{acc2, acc1},
	}
	for _, tst := range testCases {
		bindAcc(api, acc1)
		bindAcc(api, acc2)
		//prepare transaction
		sender, err := sdk.AccAddressFromBech32(tst.accFrom.Address())
		if err != nil {
			log.Printf("ERROR: AccAddressFromBech32 %s->%s", tst.accFrom.Address(), err.Error())
		}
		receiver, err := sdk.AccAddressFromBech32(tst.accTo.Address())
		if err != nil {
			log.Printf("ERROR: AccAddressFromBech32 %s->%s", tst.accTo.Address(), err.Error())
		}
		//10^18
		amount := sdk.NewInt(1500000000000000000) // 1.5
		coin := sdk.NewCoin("del0", amount)       // invalid coin

		// Prepare message
		msg := decapi.NewMsgSendCoin(sender, coin, receiver)

		// Prepare transaction arguments
		msgs := []sdk.Msg{msg}
		feeCoins := sdk.NewCoins(sdk.NewCoin(testCoin, sdk.NewInt(0)))
		memo := "test message"

		// Create signed transaction
		tx, err := api.NewSignedTransaction(msgs, feeCoins, memo, tst.accFrom)
		if err != nil {
			log.Printf("ERROR: NewSignedTransaction(from %s) %s", tst.accTo.Address(), err.Error())
		}
		log.Printf("SignedTransaction result: %s", formatAsJSON(tx))

		// Broadcast signed transaction
		broadcastTxResult, err := api.BroadcastSignedTransactionJSON(tx, tst.accFrom)
		if err != nil {
			log.Printf("ERROR: BroadcastSignedTransactionJSON(from %s) %s", tst.accTo.Address(), err.Error())
		}
		log.Printf("BroadcastSignedTransactionJSON result: %s", formatAsJSON(broadcastTxResult))
		time.Sleep(time.Second * 5)
		txRes, err := api.Transaction(broadcastTxResult.TxHash)
		log.Printf("Tx result: %s", formatAsJSON(txRes))
	}
}

func testInvalidSendSignature(api *decapi.API) {
	log.Printf("START test send")
	// make wallets
	mnemonic1 := "plug tissue today frown increase race brown sail post march trick coconut laptop churn call child question match also spend play credit already travel"
	acc1, err := wallet.NewAccountFromMnemonicWords(mnemonic1, "")
	if err != nil {
		log.Printf("ERROR: acc1 %s", err.Error())
	}
	mnemonic2 := "blade curtain rather design icon choice autumn fence budget window border inform goose segment audit cross decrease actual flight daughter beef elder auction grunt"
	acc2, err := wallet.NewAccountFromMnemonicWords(mnemonic2, "")
	if err != nil {
		log.Printf("ERROR: acc2 %s", err.Error())
	}
	// set chain id
	chainId, _ := api.ChainID()
	acc1 = acc1.WithChainID(chainId)
	acc2 = acc2.WithChainID(chainId)

	//send in both directions
	testCases := []struct {
		accFrom *wallet.Account
		accTo   *wallet.Account
	}{
		{acc1, acc2},
		{acc2, acc1},
	}
	for _, tst := range testCases {
		bindAcc(api, acc1)
		bindAcc(api, acc2)
		//prepare transaction
		sender, err := sdk.AccAddressFromBech32(tst.accFrom.Address())
		if err != nil {
			log.Printf("ERROR: AccAddressFromBech32 %s->%s", tst.accFrom.Address(), err.Error())
		}
		receiver, err := sdk.AccAddressFromBech32(tst.accTo.Address())
		if err != nil {
			log.Printf("ERROR: AccAddressFromBech32 %s->%s", tst.accTo.Address(), err.Error())
		}
		//10^18
		amount := sdk.NewInt(1500000000000000000) // 1.5
		coin := sdk.NewCoin("tdel", amount)       // invalid coin

		// Prepare message
		msg := decapi.NewMsgSendCoin(sender, coin, receiver)

		// Prepare transaction arguments
		msgs := []sdk.Msg{msg}
		feeCoins := sdk.NewCoins(sdk.NewCoin(testCoin, sdk.NewInt(0)))
		memo := "test message"

		// Create signed transaction
		tx, err := api.NewSignedTransaction(msgs, feeCoins, memo, tst.accFrom)
		if err != nil {
			log.Printf("ERROR: NewSignedTransaction(from %s) %s", tst.accTo.Address(), err.Error())
		}
		log.Printf("SignedTransaction result: %s", formatAsJSON(tx))

		// Broadcast signed transaction
		broadcastTxResult, err := api.BroadcastSignedTransactionJSON(tx, tst.accFrom)
		if err != nil {
			log.Printf("ERROR: BroadcastSignedTransactionJSON(from %s) %s", tst.accTo.Address(), err.Error())
			continue
		}
		log.Printf("BroadcastSignedTransactionJSON result: %s", formatAsJSON(broadcastTxResult))
		time.Sleep(time.Second * 5)
		txRes, err := api.Transaction(broadcastTxResult.TxHash)
		log.Printf("Tx result: %s", formatAsJSON(txRes))
	}
}

func testGovProposal(api *decapi.API) {
	log.Printf("START gov proposal")
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
	//
	sender, err := sdk.AccAddressFromBech32(acc1.Address())
	if err != nil {
		log.Printf("ERROR: AccAddressFromBech32 %s->%s", acc1.Address(), err.Error())
	}

	msg := decapi.MsgSubmitProposal{}
	//enc := []byte(`{"content":{"title":"test title", "description":"test"}}`)
	//json.Unmarshal(enc, &msg)
	msg.Content.Title = "test title"
	msg.Content.Description = "test"
	msg.Proposer = sender
	msg.VotingStartBlock = 30100
	msg.VotingEndBlock = 30000
	msgs := []sdk.Msg{msg}
	feeCoins := sdk.NewCoins(sdk.NewCoin("del", sdk.NewInt(0)))
	memo := "test message"
	// Create signed transaction
	tx, err := api.NewSignedTransaction(msgs, feeCoins, memo, acc1)
	if err != nil {
		log.Printf("ERROR: NewSignedTransaction(from %s) %s", acc1, err.Error())
	} else {
		log.Printf("SignedTransaction result: %s", formatAsJSON(tx))
	}

	// Broadcast signed transaction
	broadcastTxResult, err := api.BroadcastSignedTransactionJSON(tx, acc1)
	if err != nil {
		log.Printf("ERROR: BroadcastSignedTransactionJSON(from %s) %s", acc1, err.Error())
	} else {
		log.Printf("BroadcastSignedTransactionJSON result: %s", formatAsJSON(broadcastTxResult))
		time.Sleep(time.Second * 5)
		txRes, _ := api.Transaction(broadcastTxResult.TxHash)
		log.Printf("Tx result: %s", formatAsJSON(txRes))
	}
}

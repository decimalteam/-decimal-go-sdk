package api

import (
	amino "github.com/tendermint/tendermint/crypto/encoding/amino"

	sdk "github.com/cosmos/cosmos-sdk/types"
	auth "github.com/cosmos/cosmos-sdk/x/auth"

	"bitbucket.org/decimalteam/go-node/x/coin"
	"bitbucket.org/decimalteam/go-node/x/multisig"
	"bitbucket.org/decimalteam/go-node/x/validator"

	"bitbucket.org/decimalteam/decimal-go-sdk/wallet"
)

type (

	// MsgCreateCoin .
	MsgCreateCoin = coin.MsgCreateCoin
	// MsgSendCoin .
	MsgSendCoin = coin.MsgSendCoin
	// MsgMultiSendCoin .
	MsgMultiSendCoin = coin.MsgMultiSendCoin
	// MsgBuyCoin .
	MsgBuyCoin = coin.MsgBuyCoin
	// MsgSellCoin .
	MsgSellCoin = coin.MsgSellCoin
	// MsgSellAllCoin .
	MsgSellAllCoin = coin.MsgSellAllCoin

	// MsgCreateWallet .
	MsgCreateWallet = multisig.MsgCreateWallet
	// MsgCreateTransaction .
	MsgCreateTransaction = multisig.MsgCreateTransaction
	// MsgSignTransaction .
	MsgSignTransaction = multisig.MsgSignTransaction

	// MsgDeclareCandidate .
	MsgDeclareCandidate = validator.MsgDeclareCandidate
	// MsgDelegate .
	MsgDelegate = validator.MsgDelegate
	// MsgSetOnline .
	MsgSetOnline = validator.MsgSetOnline
	// MsgSetOffline .
	MsgSetOffline = validator.MsgSetOffline
	// MsgUnbond .
	MsgUnbond = validator.MsgUnbond
	// MsgEditCandidate .
	MsgEditCandidate = validator.MsgEditCandidate
)

// NewTransaction creates new transaction with specified parameters and messages.
func (api *API) NewTransaction(sender string, msgs []sdk.Msg, fee auth.StdFee, memo string) (tx auth.StdTx, bytesToSign []byte, err error) {

	// Request chain ID from the Decimal API if necessary
	api.ensureChainID()

	// Create transaction
	tx = auth.NewStdTx(msgs, fee, make([]auth.StdSignature, 0, 1), memo)

	// Request account number and sequence of the sender address
	accountNumber, sequence, err := api.AccountNumberAndNonce(sender)

	// Retrieve transaction bytes required to sign
	bytesToSign = auth.StdSignBytes(api.chainID, accountNumber, sequence, fee, msgs, memo)

	return
}

// SignTransaction signs transactions and append signature to transaction signatures.
func (api *API) SignTransaction(tx auth.StdTx, bytesToSign []byte, privateKey *wallet.PrivateKey) (signedTx auth.StdTx, err error) {

	privKeyBytes := append([]byte{0xE1, 0xB0, 0xF7, 0x9B, 0x20}, privateKey.GetBytes()...)
	privKey, err := amino.PrivKeyFromBytes(privKeyBytes)
	if err != nil {
		return
	}

	signatureBytes, err := privKey.Sign(bytesToSign)
	if err != nil {
		return
	}

	signature := auth.StdSignature{
		PubKey:    privKey.PubKey(),
		Signature: signatureBytes,
	}

	signedTx = tx
	signedTx.Signatures = append(signedTx.Signatures, signature)

	return
}

// EncodeTransactionBinary encodes transaction to binary form using Amino.
func (api *API) EncodeTransactionBinary(tx auth.StdTx) (txBytes []byte, err error) {
	txBytes, err = api.codec.MarshalBinaryLengthPrefixed(tx)
	return
}

// EncodeTransactionJSON encodes transaction to binary form using Amino.
func (api *API) EncodeTransactionJSON(tx auth.StdTx) (txBytes []byte, err error) {
	txBytes, err = api.codec.MarshalJSON(tx)
	return
}

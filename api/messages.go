package api

import (
	"bitbucket.org/decimalteam/go-node/x/coin"
	"bitbucket.org/decimalteam/go-node/x/multisig"
	"bitbucket.org/decimalteam/go-node/x/validator"
)

////////////////////////////////////////////////////////////////
// Module: coin
////////////////////////////////////////////////////////////////

// Messages.
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
	// MsgRedeemCheck .
	MsgRedeemCheck = coin.MsgRedeemCheck
)

// Initializing functions.
var (
	// NewMsgCreateCoin creates MsgCreateCoin message.
	NewMsgCreateCoin = coin.NewMsgCreateCoin
	// NewMsgSendCoin creates MsgSendCoin message.
	NewMsgSendCoin = coin.NewMsgSendCoin
	// NewMsgMultiSendCoin creates MsgMultiSendCoin message.
	NewMsgMultiSendCoin = coin.NewMsgMultiSendCoin
	// NewMsgBuyCoin creates MsgBuyCoin message.
	NewMsgBuyCoin = coin.NewMsgBuyCoin
	// NewMsgSellCoin creates MsgSellCoin message.
	NewMsgSellCoin = coin.NewMsgSellCoin
	// NewMsgSellAllCoin creates MsgSellAllCoin message.
	NewMsgSellAllCoin = coin.NewMsgSellAllCoin
	// NewMsgRedeemCheck creates MsgRedeemCheck message.
	NewMsgRedeemCheck = coin.NewMsgRedeemCheck
)

////////////////////////////////////////////////////////////////
// Module: multisig
////////////////////////////////////////////////////////////////

// Messages.
type (
	// MsgCreateWallet .
	MsgCreateWallet = multisig.MsgCreateWallet
	// MsgCreateTransaction .
	MsgCreateTransaction = multisig.MsgCreateTransaction
	// MsgSignTransaction .
	MsgSignTransaction = multisig.MsgSignTransaction
)

// Initializing functions.
var (
	// NewMsgCreateWallet .
	NewMsgCreateWallet = multisig.NewMsgCreateWallet
	// NewMsgCreateTransaction .
	NewMsgCreateTransaction = multisig.NewMsgCreateTransaction
	// NewMsgSignTransaction .
	NewMsgSignTransaction = multisig.NewMsgSignTransaction
)

////////////////////////////////////////////////////////////////
// Module: validator
////////////////////////////////////////////////////////////////

// Messages.
type (
	// MsgDeclareCandidate .
	MsgDeclareCandidate = validator.MsgDeclareCandidate
	// MsgEditCandidate .
	MsgEditCandidate = validator.MsgEditCandidate
	// MsgDelegate .
	MsgDelegate = validator.MsgDelegate
	// MsgUnbond .
	MsgUnbond = validator.MsgUnbond
	// MsgSetOnline .
	MsgSetOnline = validator.MsgSetOnline
	// MsgSetOffline .
	MsgSetOffline = validator.MsgSetOffline
)

// Initializing functions.
var (
	// NewMsgDeclareCandidate .
	NewMsgDeclareCandidate = validator.NewMsgDeclareCandidate
	// NewMsgEditCandidate .
	NewMsgEditCandidate = validator.NewMsgEditCandidate
	// NewMsgDelegate .
	NewMsgDelegate = validator.NewMsgDelegate
	// NewMsgUnbond .
	NewMsgUnbond = validator.NewMsgUnbond
	// NewMsgSetOnline .
	NewMsgSetOnline = validator.NewMsgSetOnline
	// NewMsgSetOffline .
	NewMsgSetOffline = validator.NewMsgSetOffline
)

package api

import (
	"bitbucket.org/decimalteam/go-node/x/coin"
	"bitbucket.org/decimalteam/go-node/x/gov"
	"bitbucket.org/decimalteam/go-node/x/multisig"
	"bitbucket.org/decimalteam/go-node/x/nft"
	"bitbucket.org/decimalteam/go-node/x/swap"
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
	// MsgUpdateCoint
	MsgUpdateCoin = coin.MsgUpdateCoin
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
	// NewMsgUpdateCoin creates MsgUpdateCoin message.
	NewMsgUpdateCoin = coin.NewMsgUpdateCoin
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
// Module: gov
////////////////////////////////////////////////////////////////

// Messages.
type (
	// MsgSubmitProposal .
	MsgSubmitProposal = gov.MsgSubmitProposal
	// MsgVote .
	MsgVote = gov.MsgVote
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
	// MsgDelegateNFT
	MsgDelegateNFT = validator.MsgDelegateNFT
	// MsgUnbondNFT
	MsgUnbondNFT = validator.MsgUnbondNFT
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
	// NewMsgDelegateNFT
	NewMsgDelegateNFT = validator.NewMsgDelegateNFT
	// NewMsgUnbondNFT
	NewMsgUnbondNFT = validator.NewMsgUnbondNFT
)

////////////////////////////////////////////////////////////////
// Module: nft
////////////////////////////////////////////////////////////////

// Messages.
type (
	// MsgMintNFT
	MsgMintNFT = nft.MsgMintNFT
	// MsgBurnNFT
	MsgBurnNFT = nft.MsgBurnNFT
	// MsgTransferNFT
	MsgTransferNFT = nft.MsgTransferNFT
	// MsgEditNFTMetadata
	MsgEditNFTMetadata = nft.MsgEditNFTMetadata
)

// Initializing functions.
var (
	// NewMsgMintNFT
	NewMsgMintNFT = nft.NewMsgMintNFT
	// NewMsgBurnNFT
	NewMsgBurnNFT = nft.NewMsgBurnNFT
	// NewMsgTransferNFT
	NewMsgTransferNFT = nft.NewMsgTranfserNFT
	// NewMsgEditNFTMetadata
	NewMsgEditNFTMetadata = nft.NewMsgEditNFTMetadata
)

////////////////////////////////////////////////////////////////
// Module: swap
////////////////////////////////////////////////////////////////

// Messages.
type (
	// MsgSwapHtlt
	MsgSwapHTLT = swap.MsgHTLT
	// MsgSwapRedeem
	MsgSwapRedeem = swap.MsgRedeem
	// MsgSwapRefund
	MsgSwapRefund = swap.MsgRefund
	// MsgSwapInitialize
	MsgSwapInitialize = swap.MsgSwapInitialize
	// MsgChainDeactivate
	MsgChainDeactivate = swap.MsgChainDeactivate
	// MsgChainActivate
	MsgChainActivate = swap.MsgChainActivate
	// MsgRedeemV2
	MsgRedeemV2 = swap.MsgRedeemV2
)

// Initializing functions.
var (
	// NewMsgSwapHtlt
	NewMsgSwapHTLT = swap.NewMsgHTLT
	// NewMsgSwapRedeem
	NewMsgSwapRedeem = swap.NewMsgRedeem
	// NewMsgSwapRefund
	NewMsgSwapRefund = swap.NewMsgRefund
	// NewMsgSwapInitialize
	NewMsgSwapInitialize = swap.NewMsgSwapInitialize
	// NewMsgChainDeactivate
	NewMsgChainDeactivate = swap.NewMsgChainDeactivate
	// NewMsgChainActivate
	NewMsgChainActivate = swap.NewMsgChainActivate
	// NewMsgRedeemV2
	NewMsgRedeemV2 = swap.NewMsgRedeemV2
)

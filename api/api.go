package api

import (
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	auth "github.com/cosmos/cosmos-sdk/x/auth"

	"bitbucket.org/decimalteam/go-node/config"
	"bitbucket.org/decimalteam/go-node/x/coin"
	"bitbucket.org/decimalteam/go-node/x/multisig"
	"bitbucket.org/decimalteam/go-node/x/nft"
	"bitbucket.org/decimalteam/go-node/x/swap"
	"bitbucket.org/decimalteam/go-node/x/validator"
)

// BaseCoinSymbol is symbol of base coin in the network.
// TODO: Request it from a gateway instead?
const BaseCoinSymbol = "tdel"

// API is a struct implementing Decimal API iteraction.
type API struct {

	// Cosmos SDK
	config *sdk.Config
	codec  *codec.Codec

	// Resty
	client *resty.Client

	// Parameters
	chainID string
}

// NewAPI creates Decimal API instance.
func NewAPI(hostURL string) *API {
	return NewAPIWithClient(hostURL, resty.New().SetTimeout(time.Minute))
}

// NewAPIWithClient creates Decimal API instance with custom Resty client.
func NewAPIWithClient(hostURL string, client *resty.Client) *API {
	return &API{
		config: newConfig(),
		codec:  newCodec(),
		client: client.SetHostURL(hostURL),
	}
}

// Config returns Cosmos SDK config.
func (api *API) Config() *sdk.Config {
	return api.config
}

// Codec returns Cosmos SDK codec.
func (api *API) Codec() *codec.Codec {
	return api.codec
}

// ChainID retrieves chain ID.
func (api *API) ChainID() (chainID string, err error) {
	url := "/rpc/genesis/chain"
	res, err := api.client.R().Get(url)
	if err != nil {
		return
	}
	if res.IsError() {
		err = NewResponseError(res)
		return
	}
	api.chainID = string(res.Body())
	chainID = api.chainID
	return
}

// newConfig initializes new Cosmos SDK configuration.
func newConfig() *sdk.Config {
	cfg := sdk.GetConfig()
	cfg.SetCoinType(60)
	cfg.SetFullFundraiserPath("44'/60'/0'/0/0")
	cfg.SetBech32PrefixForAccount(config.DecimalPrefixAccAddr, config.DecimalPrefixAccPub)
	cfg.SetBech32PrefixForValidator(config.DecimalPrefixValAddr, config.DecimalPrefixValPub)
	cfg.SetBech32PrefixForConsensusNode(config.DecimalPrefixConsAddr, config.DecimalPrefixConsPub)
	cfg.Seal()
	return cfg
}

// newCodec initializes new Cosmos SDK codec.
func newCodec() *codec.Codec {
	cdc := codec.New()
	cdc.RegisterInterface((*sdk.Msg)(nil), nil)
	cdc.RegisterConcrete(coin.MsgCreateCoin{}, "coin/create_coin", nil)
	cdc.RegisterConcrete(coin.MsgSendCoin{}, "coin/send_coin", nil)
	cdc.RegisterConcrete(coin.MsgMultiSendCoin{}, "coin/multi_send_coin", nil)
	cdc.RegisterConcrete(coin.MsgBuyCoin{}, "coin/buy_coin", nil)
	cdc.RegisterConcrete(coin.MsgSellCoin{}, "coin/sell_coin", nil)
	cdc.RegisterConcrete(coin.MsgSellAllCoin{}, "coin/sell_all_coin", nil)
	cdc.RegisterConcrete(coin.MsgUpdateCoin{}, "coin/update_coin", nil)
	cdc.RegisterConcrete(coin.MsgRedeemCheck{}, "coin/redeem_check", nil)

	cdc.RegisterConcrete(validator.MsgDeclareCandidate{}, "validator/declare_candidate", nil)
	cdc.RegisterConcrete(validator.MsgDelegate{}, "validator/delegate", nil)
	cdc.RegisterConcrete(validator.MsgSetOnline{}, "validator/set_online", nil)
	cdc.RegisterConcrete(validator.MsgSetOffline{}, "validator/set_offline", nil)
	cdc.RegisterConcrete(validator.MsgUnbond{}, "validator/unbond", nil)
	cdc.RegisterConcrete(validator.MsgEditCandidate{}, "validator/edit_candidate", nil)
	cdc.RegisterConcrete(validator.MsgDelegateNFT{}, "validator/delegate_nft", nil)
	cdc.RegisterConcrete(validator.MsgUnbondNFT{}, "validator/unbond_nft", nil)

	cdc.RegisterConcrete(nft.MsgBurnNFT{}, "nft/msg_burn", nil)
	cdc.RegisterConcrete(nft.MsgMintNFT{}, "nft/msg_mint", nil)
	cdc.RegisterConcrete(nft.MsgEditNFTMetadata{}, "nft/msg_edit_metadata", nil)
	cdc.RegisterConcrete(nft.MsgTransferNFT{}, "nft/msg_transfer", nil)

	cdc.RegisterConcrete(swap.MsgHTLT{}, "swap/msg_htlt", nil)
	cdc.RegisterConcrete(swap.MsgRedeem{}, "swap/msg_redeem", nil)
	cdc.RegisterConcrete(swap.MsgRefund{}, "swap/msg_refund", nil)
	cdc.RegisterConcrete(swap.MsgRedeemV2{}, "swap/msg_redeemv2", nil)
	cdc.RegisterConcrete(swap.MsgChainDeactivate{}, "swap/msg_chain_deactivate", nil)
	cdc.RegisterConcrete(swap.MsgChainActivate{}, "swap/msg_chain_activate", nil)
	cdc.RegisterConcrete(swap.MsgSwapInitialize{}, "swap/msg_swap_initialize", nil)

	cdc.RegisterConcrete(multisig.MsgCreateWallet{}, "multisig/create_wallet", nil)
	cdc.RegisterConcrete(multisig.MsgCreateTransaction{}, "multisig/create_transaction", nil)
	cdc.RegisterConcrete(multisig.MsgSignTransaction{}, "multisig/sign_transaction", nil)
	auth.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	cdc.Seal()
	return cdc
}

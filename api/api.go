package api

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	auth "github.com/cosmos/cosmos-sdk/x/auth"

	"bitbucket.org/decimalteam/go-node/config"
	"bitbucket.org/decimalteam/go-node/x/coin"
	"bitbucket.org/decimalteam/go-node/x/gov"
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

	// Direct/Gate
	directConn *DirectConn

	// Resty (REST, RPC)
	client *clientConn

	// Parameters
	chainID string
}

// Ports for REST/RPC interfaces.
type DirectConn struct {
	// ":port"
	PortREST string
	PortRPC  string
}

type clientConn struct {
	rest *resty.Client
	rpc  *resty.Client
}

// NewAPI creates Decimal API instance.
// If directConn is nil then used gateway;
// If directConn is not nil AND directConn.PortREST = "" then used defaultPortREST (:1317);
// If directConn is not nil AND directConn.PortRPC  = "" then used defaultPortRPC (:26657);
func NewAPI(hostURL string, directConn *DirectConn) *API {
	return NewAPIWithClient(
		hostURL,
		resty.New().SetTimeout(time.Minute),
		resty.New().SetTimeout(time.Minute),
		directConn,
	)
}

// NewAPIWithClient creates Decimal API instance with custom Resty client.
func NewAPIWithClient(hostURL string, restClient *resty.Client, rpcClient *resty.Client, directConn *DirectConn) *API {
	const (
		defaultPortREST = ":1317"
		defaultPortRPC  = ":26657"
	)
	var (
		hostREST = hostURL
		hostRPC  = hostURL
	)

	if directConn != nil {
		if directConn.PortREST == "" {
			directConn.PortREST = defaultPortREST
		}
		if directConn.PortRPC == "" {
			directConn.PortRPC = defaultPortRPC
		}
		hostREST += directConn.PortREST
		hostRPC += directConn.PortRPC
	}

	return &API{
		config: newConfig(),
		codec:  newCodec(),
		client: &clientConn{
			rest: restClient.SetHostURL(hostREST),
			rpc:  rpcClient.SetHostURL(hostRPC),
		},
		directConn: directConn,
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
	type respDirectChainID struct {
		Result struct {
			NodeInfo struct {
				Network string `json:"network"`
			} `json:"node_info"`
		} `json:"result"`
	}

	var (
		dres = respDirectChainID{}
		url  = ""
	)

	if api.directConn == nil {
		url = "/rpc/genesis/chain"
	} else {
		url = "/status"
	}

	res, err := api.client.rpc.R().Get(url)
	if err != nil {
		return
	}
	if res.IsError() {
		err = NewResponseError(res)
		return
	}

	if api.directConn == nil {
		api.chainID = string(res.Body())
	} else {
		err = json.Unmarshal(res.Body(), &dres)
		if err != nil {
			return
		}
		api.chainID = dres.Result.NodeInfo.Network
	}

	chainID = api.chainID
	return
}

// Return current height (block number) of blockchain
func (api *API) GetHeight() (uint64, error) {
	//api: /api/blocks?limit=1&offset=0
	//rest: /status
	if api.directConn == nil {
		return api.apiGetHeight()
	} else {
		return api.restGetHeight()
	}
}

func (api *API) apiGetHeight() (uint64, error) {
	type responseType struct {
		OK     bool `json:"ok"`
		Result struct {
			Blocks []struct {
				Height uint64 `json:"height"`
			} `json:"blocks"`
		} `json:"result"`
	}
	res, err := api.client.rpc.R().Get("/blocks?limit=1&offset=0")
	if err = processConnectionError(res, err); err != nil {
		return 0, err
	}
	response := responseType{}
	err = json.Unmarshal(res.Body(), &response)
	if err != nil || !response.OK {
		responseError := JsonRPCError{}
		err = json.Unmarshal(res.Body(), &responseError)
		if err != nil {
			return 0, err
		}
		return 0, fmt.Errorf("received response containing error: %s", responseError.Error())
	}
	if len(response.Result.Blocks) > 0 {
		return response.Result.Blocks[0].Height, nil
	}
	return 0, fmt.Errorf("Response without blocks")
}

func (api *API) restGetHeight() (uint64, error) {
	type responseType struct {
		Result struct {
			SyncInfo struct {
				Height string `json:"latest_block_height"`
			} `json:"sync_info"`
		} `json:"result"`
	}
	res, err := api.client.rpc.R().Get("/status")
	if err = processConnectionError(res, err); err != nil {
		return 0, err
	}
	response := responseType{}
	err = json.Unmarshal(res.Body(), &response)
	if err != nil {
		responseError := JsonRPCError{}
		err = json.Unmarshal(res.Body(), &responseError)
		if err != nil {
			return 0, err
		}
		return 0, fmt.Errorf("received response containing error: %s", responseError.Error())
	}
	return strconv.ParseUint(response.Result.SyncInfo.Height, 10, 64)
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
	cdc.RegisterConcrete(nft.MsgUpdateReserveNFT{}, "nft/update_reserve", nil)

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

	cdc.RegisterConcrete(gov.MsgSubmitProposal{}, "cosmos-sdk/MsgSubmitProposal", nil)
	cdc.RegisterConcrete(gov.MsgVote{}, "cosmos-sdk/MsgVote", nil)
	cdc.RegisterConcrete(gov.MsgSoftwareUpgradeProposal{}, "cosmos-sdk/MsgSoftwareUpgradeProposal", nil)

	auth.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	cdc.Seal()
	return cdc
}

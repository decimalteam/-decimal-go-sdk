package api

import (
	"encoding/json"
	"fmt"
	"strings"

	"bitbucket.org/decimalteam/decimal-go-sdk/wallet"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

// const hugeGas = uint64(16 * 1024)
const hugeGas = uint64(16 << 10)

// BroadcastTxResponse contains API response.
type BroadcastTxResponse struct {
	OK     bool               `json:"ok"`
	Result *BroadcastTxResult `json:"result"`
}

// BroadcastTxResult contains API response fields.
type BroadcastTxResult struct {
	Height string `json:"height"`
	TxHash string `json:"txhash"`
	Code   int    `json:"code"`
	RawLog string `json:"raw_log"`
}

const (
	prefix       = `{"type":"cosmos-sdk/StdTx","value":`
	suffix       = `}`
	prefixLength = len(prefix)
	suffixLength = len(suffix)
)

// NOTE: To ensure that transaction was successfully committed to the blockchain,
// you need to find the transaction by the hash and ensure that the status code equals to 0.

// NewSignedTransaction creates and signs a transaction.
func (api *API) NewSignedTransaction(msgs []sdk.Msg, feeCoins sdk.Coins, memo string, account *wallet.Account) (tx auth.StdTx, err error) {
	// Adjust gas until it is equal to gasEstimated
	for gas, gasEstimated := hugeGas, uint64(0); gas != gasEstimated; {
		if gasEstimated != 0 {
			gas = gasEstimated
		}

		// Create and sign transaction
		fee := auth.NewStdFee(gas, feeCoins)
		tx = account.CreateTransaction(msgs, fee, memo)
		tx, err = account.SignTransaction(tx)
		if err != nil {
			return
		}

		// Estimate and adjust amount of gas wanted for the transaction
		gasEstimated, err = api.EstimateTransactionGasWanted(tx)
		if err != nil {
			return
		}
	}
	return
}

// BroadcastSignedTransactionJSON sends transaction (presented in JSON format) to the node and returns the result.
// If transaction is sucessful, it modified account sequence
func (api *API) BroadcastSignedTransactionJSON(tx auth.StdTx, acc *wallet.Account) (*BroadcastTxResult, error) {
	var (
		url = ""
	)

	// Marshal transaction to special JSON format
	txJSONBytes, err := api.codec.MarshalJSON(tx)
	if err != nil {
		return nil, err
	}
	txJSON := string(txJSONBytes)

	// Adjust format of broadcasting JSON object
	if strings.HasPrefix(txJSON, prefix) && strings.HasSuffix(txJSON, suffix) {
		txJSON = fmt.Sprintf(`{"tx":%s,"mode":"sync"}`, txJSON[prefixLength:len(txJSON)-suffixLength])
	}

	// Send POST request at path `/rpc/txs` and wait for the response
	if api.directConn == nil {
		url = "/rpc/txs"
	} else {
		url = "/txs"
	}

	// TODO: undefined /txs in RPC, but was found /txs in REST?
	res, err := api.client.rest.R().SetBody(txJSON).Post(url)
	if err = processConnectionError(res, err); err != nil {
		return nil, err
	}
	// Unmarshal response from JSON format
	response := BroadcastTxResult{}
	err = json.Unmarshal(res.Body(), &response)

	// Check transaction execution code (success or fail)
	if err != nil || response.Code != 0 {
		txError := TxError{}
		err = json.Unmarshal(res.Body(), &txError)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("received tx error: %s", txError.Error())
	}

	acc.WithSequence(uint64(acc.Sequence() + 1))

	return &response, nil
}

// BroadcastRawSignedTransaction sends raw transaction to the node and returns the result.
func (api *API) BroadcastRawSignedTransaction(tx auth.StdTx) (*BroadcastTxResult, error) {

	txBytes, err := api.codec.MarshalBinaryLengthPrefixed(tx)
	if err != nil {
		return nil, err
	}
	txHex := fmt.Sprintf("0x%x", txBytes)

	url := "" // TODO: url is nil.
	res, err := api.client.rest.R().SetQueryParam("tx", txHex).Get(url)
	if err = processConnectionError(res, err); err != nil {
		return nil, err
	}

	response := BroadcastTxResponse{}
	err = json.Unmarshal(res.Body(), &response)
	if err != nil || !response.OK {
		responseError := Error{}
		err = json.Unmarshal(res.Body(), &responseError)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("received response containing error: %s", responseError.Error())
	}

	return response.Result, nil
}

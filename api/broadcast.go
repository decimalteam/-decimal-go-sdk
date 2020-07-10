package api

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/x/auth"
)

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

// BroadcastTransactionJSON sends transaction (presented in JSON format) to the node and returns the result.
func (api *API) BroadcastTransactionJSON(tx auth.StdTx) (*BroadcastTxResult, error) {

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
	res, err := api.client.R().SetBody(txJSON).Post("/rpc/txs")
	if err != nil {
		return nil, err
	}
	if res.IsError() {
		return nil, NewResponseError(res)
	}

	// Unmarshal response from JSON format
	response := BroadcastTxResult{}
	err = json.Unmarshal(res.Body(), &response)
	if err != nil {
		return nil, err
	}

	// Check transaction execution code (success or fail)
	if response.Code != 0 {
		txError := TxError{}
		err = json.Unmarshal(res.Body(), &txError)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("received tx error: %s", txError.Error())
	}

	return &response, nil
}

// BroadcastRawTransaction sends raw transaction to the node and returns the result.
func (api *API) BroadcastRawTransaction(tx auth.StdTx) (*BroadcastTxResult, error) {

	txBytes, err := api.codec.MarshalBinaryLengthPrefixed(tx)
	if err != nil {
		return nil, err
	}
	txHex := fmt.Sprintf("0x%x", txBytes)

	url := "" // TODO
	res, err := api.client.R().SetQueryParam("tx", txHex).Get(url)
	if err != nil {
		return nil, err
	}
	if res.IsError() {
		return nil, NewResponseError(res)
	}

	response := BroadcastTxResponse{}
	err = json.Unmarshal(res.Body(), &response)
	if err != nil {
		return nil, err
	}

	if !response.OK {
		responseError := Error{}
		err = json.Unmarshal(res.Body(), &responseError)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("received response containing error: %s", responseError.Error())
	}

	return response.Result, nil
}

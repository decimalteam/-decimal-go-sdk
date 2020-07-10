package api

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/x/auth"
)

// SendTransactionResponse contains API response.
type SendTransactionResponse struct {
	OK     bool                   `json:"ok"`
	Result *SendTransactionResult `json:"result"`
}

// SendTransactionResult contains API response fields.
type SendTransactionResult struct {
	Height string `json:"height"`
	TxHash string `json:"txhash"`
	Code   int    `json:"code"`
	RawLog string `json:"raw_log"`
}

// NOTE: To ensure that transaction was successfully committed to the blockchain,
// you need to find the transaction by the hash and ensure that the status code equals to 0.

// SendTransactionJSON sends presented in JSON format transaction to the node and returns the result.
func (api *API) SendTransactionJSON(tx auth.StdTx) (*SendTransactionResult, error) {

	txJSONBytes, err := api.codec.MarshalJSON(tx)
	if err != nil {
		return nil, err
	}
	txJSON := string(txJSONBytes)

	// Fix format of resulting JSON object
	prefix, suffix := `{"type":"cosmos-sdk/StdTx","value":`, `}`
	prefixLength, suffixLength := len(prefix), len(suffix)
	if strings.HasPrefix(txJSON, prefix) && strings.HasSuffix(txJSON, suffix) {
		txJSON = fmt.Sprintf(`{"tx":%s,"mode":"sync"}`, txJSON[prefixLength:len(txJSON)-suffixLength])
	}

	url := "/rpc/txs"
	res, err := api.client.R().SetBody(txJSON).Post(url)
	if err != nil {
		return nil, err
	}
	if res.IsError() {
		return nil, NewResponseError(res)
	}

	response := SendTransactionResult{}
	err = json.Unmarshal(res.Body(), &response)
	if err != nil {
		return nil, err
	}

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

// // SendRawTransaction sends raw transaction to the node and returns the response.
// func (api *API) SendRawTransaction(tx auth.StdTx) (*SendTransactionResult, error) {

// 	txBytes, err := api.codec.MarshalBinaryLengthPrefixed(tx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	txHex := fmt.Sprintf("0x%x", txBytes)

// 	url := "" // TODO
// 	res, err := api.client.R().SetQueryParam("tx", txHex).Get(url)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if res.IsError() {
// 		return nil, NewResponseError(res)
// 	}

// 	response := SendTransactionResponse{}
// 	err = json.Unmarshal(res.Body(), &response)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if !response.OK {
// 		responseError := Error{}
// 		err = json.Unmarshal(res.Body(), &responseError)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return nil, fmt.Errorf("received response containing error: %s", responseError.Error())
// 	}

// 	return response.Result, nil
// }

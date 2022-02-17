package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// TransactionResponse contains API response.
type TransactionResponse struct {
	JSONRPC string             `json:"jsonrpc"`
	ID      int64              `json:"id"`
	Result  *TransactionResult `json:"result"`
}

// TransactionResult contains API response fields.
type TransactionResult struct {
	Hash     string    `json:"hash"`
	Height   string    `json:"height"`
	Index    uint64    `json:"index"`
	TxResult *TxResult `json:"tx_result"`
	Tx       string    `json:"tx"`
}

// TxResult contains API response fields.
type TxResult struct {
	Code      int64           `json:"code"`
	Data      string          `json:"data"`
	Log       string          `json:"log"`
	LogParsed []TxLog         `json:"-"`
	Info      string          `json:"info"`
	GasWanted string          `json:"gasWanted"`
	GasUsed   string          `json:"gasUsed"`
	Events    []TxEventBase64 `json:"events"`
	Codespace string          `json:"codespace"`
}

// TxLog contains API response fields.
type TxLog struct {
	MsgIndex uint64    `json:"msg_index"`
	Log      string    `json:"log"`
	Events   []TxEvent `json:"events"`
}

// TxEvent contains API response fields.
type TxEvent struct {
	Type       string        `json:"type"`
	Attributes []TxAttribute `json:"attributes"`
}

// TxAttribute contains API response fields.
type TxAttribute struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// TxEventBase64 contains API response fields.
type TxEventBase64 struct {
	Type       string              `json:"type"`
	Attributes []TxAttributeBase64 `json:"attributes"`
}

// TxAttributeBase64 contains API response fields.
type TxAttributeBase64 struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Transaction requests full information about transaction with specified hash.
// NOTE: It is expected that `txHash` encoded in hex format and written
// in capital letters and without "0x" at the beginning.
func (api *API) Transaction(txHash string) (*TransactionResult, error) {

	url := fmt.Sprintf("/rpc/tx?hash=%s", txHash)
	res, err := api.client.R().Get(url)
	if err != nil {
		return nil, err
	}
	if res.IsError() {
		return nil, NewResponseError(res)
	}

	response := TransactionResponse{}
	err = json.Unmarshal(res.Body(), &response)
	if err != nil || response.Result == nil {
		responseError := JsonRPCError{}
		err = json.Unmarshal(res.Body(), &responseError)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("received response containing error: %s", responseError.Error())
	}

	// Parse `log` value presented as string to []TxLog array
	txLogs := []TxLog{}
	json.Unmarshal([]byte(response.Result.TxResult.Log), &txLogs)
	response.Result.TxResult.LogParsed = txLogs

	return response.Result, nil
}

////////////////////////////////////////////////////////////////
// TxAttributeBase64
////////////////////////////////////////////////////////////////

// MarshalJSON implements Marshaler interface.
func (a *TxAttributeBase64) MarshalJSON() ([]byte, error) {
	attributeJSON := struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}{
		Key:   base64.StdEncoding.EncodeToString([]byte(a.Key)),
		Value: base64.StdEncoding.EncodeToString([]byte(a.Value)),
	}
	return json.Marshal(attributeJSON)
}

// UnmarshalJSON implements Marshaler interface.
func (a *TxAttributeBase64) UnmarshalJSON(b []byte) error {
	attributeJSON := struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}{}
	err := json.Unmarshal(b, &attributeJSON)
	if err != nil {
		return err
	}
	key, _ := base64.StdEncoding.DecodeString(attributeJSON.Key)
	a.Key = string(key)
	value, _ := base64.StdEncoding.DecodeString(attributeJSON.Value)
	a.Value = string(value)
	return nil
}

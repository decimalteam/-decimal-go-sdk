package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
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

// TxResponse by Hash contains API response.
type TxResponse struct {
	OK     bool     `json:"ok"`
	Result *TxCheck `json:"result"`
}

// TxCheck contains API response fields.
type TxCheck struct {
	Hash   string `json:"hash"`
	Status string `json:"status"`
}

// Transaction requests full information about transaction with specified hash.
// NOTE: It is expected that `txHash` encoded in hex format and written
// in capital letters and without "0x" at the beginning.
func (api *API) Transaction(txHash string) (*TransactionResult, error) {
	if api.directConn == nil {
		return api.apiTransaction(txHash)
	} else {
		return api.restTransaction(txHash)
	}
}

func (api *API) apiTransaction(txHash string) (*TransactionResult, error) {
	// request
	res, err := api.client.rpc.R().Get(fmt.Sprintf("/rpc/tx?hash=%s", txHash))
	if err = processConnectionError(res, err); err != nil {
		return nil, err
	}
	// json decode
	respValue, respErr := TransactionResponse{}, JsonRPCError{}
	err = universalJSONDecode(res.Body(), &respValue, &respErr, func() (bool, bool) {
		return respValue.Result != nil, respErr.InternalError.Code != 0
	})
	if err != nil {
		return nil, joinErrors(err, respErr)
	}
	//
	// Parse `log` value presented as string to []TxLog array
	txLogs := []TxLog{}
	json.Unmarshal([]byte(respValue.Result.TxResult.Log), &txLogs)
	respValue.Result.TxResult.LogParsed = txLogs
	return respValue.Result, nil
}

func (api *API) restTransaction(txHash string) (*TransactionResult, error) {
	// request
	res, err := api.client.rest.R().Get(fmt.Sprintf("/txs/%s", txHash))
	if err = processConnectionError(res, err); err != nil {
		return nil, err
	}
	// json decode
	respValue, respErr := TransactionResponse{}, JsonRPCError{}
	err = universalJSONDecode(res.Body(), &respValue, &respErr, func() (bool, bool) {
		return respValue.Result != nil, respErr.InternalError.Code != 0
	})
	if err != nil {
		return nil, joinErrors(err, respErr)
	}
	//
	txLogs := []TxLog{}
	json.Unmarshal([]byte(respValue.Result.TxResult.Log), &txLogs)
	respValue.Result.TxResult.LogParsed = txLogs
	return respValue.Result, nil
}

func (api *API) CheckTransaction(txHash string) (*TxCheck, error) {
	if api.directConn == nil {
		return api.apiCheckTransaction(txHash)
	} else {
		return api.restCheckTransaction(txHash)
	}
}

func (api *API) restCheckTransaction(txHash string) (*TxCheck, error) {
	url := fmt.Sprintf("/tx?hash=0x%s", txHash)
	res, err := api.client.rpc.R().Get(url)
	if err = processConnectionError(res, err); err != nil {
		return nil, err
	}

	respValue, respErr := TransactionResponse{}, JsonRPCError{}
	err = universalJSONDecode(res.Body(), &respValue, &respErr, func() (bool, bool) {
		return respValue.Result != nil, respErr.InternalError.Code != 0
	})
	if err != nil {
		return nil, joinErrors(err, respErr)
	}

	result := TxCheck{}
	result.Hash = respValue.Result.Hash
	switch respValue.Result.TxResult.Code {
	case 0:
		result.Status = "Success"
	default:
		result.Status = "Failure"
	}

	return &result, nil
}

func (api *API) apiCheckTransaction(txHash string) (*TxCheck, error) {
	url := fmt.Sprintf("/rpc/tx?hash=%s", txHash)
	res, err := api.client.rpc.R().Get(url)
	if err = processConnectionError(res, err); err != nil {
		return nil, err
	}

	respValue, respErr := TransactionResponse{}, JsonRPCError{}
	err = universalJSONDecode(res.Body(), &respValue, &respErr, func() (bool, bool) {
		return respValue.Result != nil, respErr.InternalError.Code != 0
	})
	if err != nil {
		return nil, joinErrors(err, respErr)
	}

	result := TxCheck{}
	result.Hash = respValue.Result.Hash
	switch respValue.Result.TxResult.Code {
	case 0:
		result.Status = "Success"
	default:
		result.Status = "Failure"
	}

	return &result, nil
}

// TransactionsByBlock return all transactions hashes in block
func (api *API) TransactionsByBlock(height uint64) ([]string, error) {
	if api.directConn == nil {
		return api.apiTransactionsByBlock(height)
	} else {
		return api.restTransactionsByBlock(height)
	}
}

func (api *API) apiTransactionsByBlock(height uint64) ([]string, error) {
	type responseType struct {
		OK     bool `json:"ok"`
		Result struct {
			Count int64 `json:"count"`
			Txs   []struct {
				Hash string `json:"hash"`
			} `json:"txs"`
		} `json:"result"`
	}
	// request
	res, err := api.client.rpc.R().Get(fmt.Sprintf("/block/%d/txs", height))
	if err = processConnectionError(res, err); err != nil {
		return nil, err
	}
	// json decode
	respValue, respErr := responseType{}, JsonRPCError{}
	err = universalJSONDecode(res.Body(), &respValue, &respErr, func() (bool, bool) {
		return respValue.OK, respErr.InternalError.Code != 0
	})
	if err != nil {
		return nil, joinErrors(err, respErr)
	}
	// process result
	result := make([]string, 0, respValue.Result.Count)
	for _, tx := range respValue.Result.Txs {
		result = append(result, tx.Hash)
	}
	return result, nil
}

func (api *API) restTransactionsByBlock(height uint64) ([]string, error) {
	type responseType struct {
		TotalCount string `json:"total_count"`
		Count      string `json:"count"`
		Txs        []struct {
			Hash string `json:"txhash"`
		} `json:"txs"`
	}
	// request
	res, err := api.client.rest.R().Get(fmt.Sprintf("/txs?tx.minheight=%d&tx.maxheight=%d&limit=%d", height, height, 1000))
	if err = processConnectionError(res, err); err != nil {
		return nil, err
	}
	// json decode
	respValue := responseType{}
	err = universalJSONDecode(res.Body(), &respValue, nil, func() (bool, bool) {
		return respValue.Count > "", false
	})
	if err != nil {
		return nil, err
	}
	// process result
	// TODO: pagination
	totalCount, _ := strconv.ParseUint(respValue.TotalCount, 10, 64)
	result := make([]string, 0, totalCount)
	for _, tx := range respValue.Txs {
		result = append(result, tx.Hash)
	}
	return result, nil
}

// //////////////////////////////////////////////////////////////
// TxAttributeBase64
// //////////////////////////////////////////////////////////////

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

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

// Transaction requests full information about transaction with specified hash.
// NOTE: It is expected that `txHash` encoded in hex format and written
// in capital letters and without "0x" at the beginning.
func (api *API) Transaction(txHash string) (*TransactionResult, error) {
	var (
		url = ""
	)

	if api.directConn == nil {
		url = fmt.Sprintf("/rpc/tx?hash=%s", txHash)
	} else {
		// TODO
		url = fmt.Sprintf("/tx?hash=0x%s", txHash)
	}

	res, err := api.client.rpc.R().Get(url)
	if err = processConnectionError(res, err); err != nil {
		return nil, err
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

//TransactionsByBlock return all transactions hashes in block
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
	res, err := api.client.rpc.R().Get(fmt.Sprintf("/block/%d/txs", height))
	if err = processConnectionError(res, err); err != nil {
		return []string{}, err
	}
	response := responseType{}
	err = json.Unmarshal(res.Body(), &response)
	if err != nil || !response.OK {
		responseError := JsonRPCError{}
		err = json.Unmarshal(res.Body(), &responseError)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("received response containing error: %s", responseError.Error())
	}
	result := make([]string, 0, response.Result.Count)
	for _, tx := range response.Result.Txs {
		result = append(result, tx.Hash)
	}
	return result, nil
}

func (api *API) restTransactionsByBlock(height uint64) ([]string, error) {
	type responseType struct {
		TotalCount string `json:"total_count"`
		Count      string `json:"count"`
		Txs        []struct {
			Hash string `json:"hash"`
		} `json:"txs"`
	}
	res, err := api.client.rest.R().Get(fmt.Sprintf("/txs?tx.minheight=%d&tx.maxheight=%d&limit=%d", height, height, 1000))
	if err = processConnectionError(res, err); err != nil {
		return []string{}, err
	}
	response := responseType{}
	err = json.Unmarshal(res.Body(), &response)
	if err != nil {
		return nil, fmt.Errorf("response unmarshaling error: %s", err.Error())
	}
	//TODO: pagination
	totalCount, _ := strconv.ParseUint(response.TotalCount, 10, 64)
	result := make([]string, 0, totalCount)
	for _, tx := range response.Txs {
		result = append(result, tx.Hash)
	}
	return result, nil
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

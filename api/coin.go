package api

import (
	"encoding/json"
	"fmt"
)

// CoinResponse contains API response.
type CoinResponse struct {
	OK     bool        `json:"ok,omitempty"`
	Result *CoinResult `json:"result,omitempty"`
}

// CoinResult contains API response fields.
type CoinResult struct {
	Symbol      string `json:"symbol"`
	Title       string `json:"title"`
	Crr         uint8  `json:"crr"`
	Reserve     string `json:"reserve"`
	Volume      string `json:"volume"`
	LimitVolume string `json:"limitVolume"`
	Creator     string `json:"creator"`           // Address of account created the coin
	TxHash      string `json:"txHash,omitempty"`  // Hash of transaction in which the coin was created
	BlockID     uint64 `json:"blockId,omitempty"` // Number of block in which the coin was created
	Avatar      string `json:"avatar,omitempty"`  // Optional avatar info presented in base64 format
}

// Coin requests full information about coin with specified symbol.
func (api *API) Coin(symbol string) (*CoinResult, error) {

	url := fmt.Sprintf("/coin/%s", symbol)
	res, err := api.client.R().Get(url)
	if err != nil {
		return nil, err
	}
	if res.IsError() {
		return nil, NewResponseError(res)
	}

	response := CoinResponse{}
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
	}

	return response.Result, nil
}

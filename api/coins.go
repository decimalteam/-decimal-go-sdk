package api

import (
	"encoding/json"
	"fmt"
)

// CoinResponse contains API response.
type CoinResponse struct {
	OK     bool        `json:"ok"`
	Result *CoinResult `json:"result"`
}

// CoinsResponse contains API response.
type CoinsResponse struct {
	OK     bool `json:"ok"`
	Result struct {
		Count uint64        `json:"count"`
		Coins []*CoinResult `json:"coins"`
	} `json:"result"`
}

// CoinResult contains API response fields.
type CoinResult struct {
	Symbol      string `json:"symbol"`
	Title       string `json:"title"`
	Crr         uint8  `json:"crr"`
	Reserve     string `json:"reserve"`
	Volume      string `json:"volume"`
	LimitVolume string `json:"limitVolume"`

	Creator string `json:"creator"` // Address of account created the coin
	TxHash  string `json:"txHash"`  // Hash of transaction in which the coin was created
	BlockID uint64 `json:"blockId"` // Number of block in which the coin was created
	Avatar  string `json:"avatar"`  // Optional avatar info presented in base64 format
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

// Coins requests full information about all coins.
func (api *API) Coins() ([]*CoinResult, error) {

	url := "/coin"
	res, err := api.client.R().Get(url)
	if err != nil {
		return nil, err
	}
	if res.IsError() {
		return nil, NewResponseError(res)
	}

	response := CoinsResponse{}
	err = json.Unmarshal(res.Body(), &response)
	if err != nil || !response.OK {
		responseError := Error{}
		err = json.Unmarshal(res.Body(), &responseError)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("received response containing error: %s", responseError.Error())
	}

	return response.Result.Coins, nil
}

package api

import (
	"encoding/json"
)

// CoinsResponse contains API response.
type CoinsResponse struct {
	OK     bool `json:"ok,omitempty"`
	Result struct {
		Count uint64        `json:"count,omitempty"`
		Coins []*CoinResult `json:"coins,omitempty"`
	} `json:"result,omitempty"`
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

	return response.Result.Coins, nil
}

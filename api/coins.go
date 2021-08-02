package api

import (
	"encoding/json"
	"errors"
	"fmt"
)

// CoinResult contains API response fields.
type CoinResult struct {
	Symbol      string `json:"symbol"`
	Title       string `json:"title"`
	Crr         uint8  `json:"crr"`
	Reserve     string `json:"reserve"`
	Volume      string `json:"volume"`
	LimitVolume string `json:"limitVolume"`

	// TODO: can't get creator, txHash, blockId, avatar from REST
	Creator string `json:"creator"` // Address of account created the coin
	TxHash  string `json:"txHash"`  // Hash of transaction in which the coin was created
	BlockID uint64 `json:"blockId"` // Number of block in which the coin was created
	Avatar  string `json:"avatar"`  // Optional avatar info presented in base64 format
}

// Coin requests full information about coin with specified symbol.
func (api *API) Coin(symbol string) (*CoinResult, error) {
	type respCoin struct {
		OK     bool        `json:"ok"`
		Result *CoinResult `json:"result"`
	}

	url := fmt.Sprintf("/coin/%s", symbol)

	res, err := api.client.rest.R().Get(url)
	if err != nil {
		return nil, err
	}
	if res.IsError() {
		return nil, NewResponseError(res)
	}

	response := respCoin{}
	err = json.Unmarshal(res.Body(), &response)
	if err != nil || (api.directConn == nil && !response.OK) {
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
	type respCoins struct {
		OK     bool `json:"ok"`
		Result struct {
			Count uint64        `json:"count"`
			Coins []*CoinResult `json:"coins"`
		} `json:"result"`
	}
	type respDirectCoins struct {
		Result []string
	}

	var (
		gres = respCoins{}
		dres = respDirectCoins{}
		url  = ""
	)

	if api.directConn == nil {
		url = "/coin"
	} else {
		url = "/coins"
	}

	res, err := api.client.rest.R().Get(url)
	if err != nil {
		return nil, err
	}
	if res.IsError() {
		return nil, NewResponseError(res)
	}

	if api.directConn == nil {
		err = json.Unmarshal(res.Body(), &gres)
	} else {
		err = json.Unmarshal(res.Body(), &dres)
	}

	if err != nil || (api.directConn == nil && !gres.OK) {
		responseError := Error{}
		err = json.Unmarshal(res.Body(), &responseError)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("received response containing error: %s", responseError.Error())
	}

	if api.directConn == nil {
		return gres.Result.Coins, nil
	}

	coins := []*CoinResult{}
	errstr := ""

	for _, val := range dres.Result {
		coin, err := api.Coin(val)
		if err != nil {
			errstr += err.Error()
			continue
		}
		coins = append(coins, coin)
	}

	if errstr != "" {
		err = errors.New(errstr)
	}

	return coins, err
}

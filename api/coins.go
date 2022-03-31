package api

import (
	"errors"
	"fmt"
	"strconv"
)

// CoinResult contains API response fields.
type CoinResult struct {
	Symbol      string `json:"symbol"`
	Title       string `json:"title"`
	Crr         uint8  `json:"crr"`
	Reserve     string `json:"reserve"`
	Volume      string `json:"volume"`
	LimitVolume string `json:"limitVolume"`
	Creator     string `json:"creator"` // Address of account created the coin

	// TODO: can't get price USD, txHash, blockId, avatar from REST
	PriceUSD        string `json:"priceUSD"`
	TxHash          string `json:"txHash"`  // Hash of transaction in which the coin was created
	BlockID         uint64 `json:"blockId"` // Number of block in which the coin was created
	Avatar          string `json:"avatar"`  // Optional avatar info presented in base64 format
	ContractAddress string `json:"contractAddress"`
}

// Coin requests full information about coin with specified symbol.
// Gateway: ok, REST/RPC: partial
func (api *API) Coin(symbol string) (*CoinResult, error) {
	if api.directConn == nil {
		return api.apiCoin(symbol)
	} else {
		return api.restCoin(symbol)
	}
}

func (api *API) apiCoin(symbol string) (*CoinResult, error) {
	type respCoin struct {
		OK     bool        `json:"ok"`
		Result *CoinResult `json:"result"`
	}
	//request
	res, err := api.client.rest.R().Get(fmt.Sprintf("/coin/%s", symbol))
	if err = processConnectionError(res, err); err != nil {
		return nil, err
	}
	//json decode
	respValue, respErr := respCoin{}, JsonRPCInternalError{}
	err = universalJSONDecode(res.Body(), &respValue, &respErr, func() (bool, bool) {
		return respValue.OK, respErr.Code != 0
	})
	if err != nil {
		return nil, joinErrors(err, respErr)
	}
	//process result
	return respValue.Result, nil
}
func (api *API) restCoin(symbol string) (*CoinResult, error) {
	type respDirectCoin struct {
		Result struct {
			Symbol      string `json:"symbol"`
			Title       string `json:"title"`
			Crr         string `json:"constant_reserve_ratio"`
			Reserve     string `json:"reserve"`
			LimitVolume string `json:"limit_volume"`
			Volume      string `json:"volume"`
			Creator     string `json:"creator"`
		} `json:"result"`
	}
	//request
	res, err := api.client.rest.R().Get(fmt.Sprintf("/coin/%s", symbol))
	if err = processConnectionError(res, err); err != nil {
		return nil, err
	}
	//json decode
	respValue := respDirectCoin{}
	err = universalJSONDecode(res.Body(), &respValue, nil, func() (bool, bool) {
		return respValue.Result.Symbol > "", false
	})
	if err != nil {
		return nil, err
	}
	//process result
	result := &CoinResult{}
	result.Title = respValue.Result.Title
	result.Symbol = respValue.Result.Symbol
	crr, _ := strconv.ParseUint(respValue.Result.Crr, 10, 8)
	result.Crr = uint8(crr)
	result.Reserve = respValue.Result.Reserve
	result.LimitVolume = respValue.Result.LimitVolume
	result.Volume = respValue.Result.Volume
	result.Creator = respValue.Result.Creator
	return result, nil
}

// Coins requests full information about all coins.
// Gateway: ok, REST/RPC: partial
func (api *API) Coins() ([]*CoinResult, error) {
	if api.directConn == nil {
		return api.apiCoins()
	} else {
		return api.restCoins()
	}
}

func (api *API) apiCoins() ([]*CoinResult, error) {
	type respCoins struct {
		OK     bool `json:"ok"`
		Result struct {
			Count uint64        `json:"count"`
			Coins []*CoinResult `json:"coins"`
		} `json:"result"`
	}
	//request
	res, err := api.client.rest.R().Get("/coin")
	if err = processConnectionError(res, err); err != nil {
		return nil, err
	}
	//json decode
	respValue, respErr := respCoins{}, Error{}
	err = universalJSONDecode(res.Body(), &respValue, &respErr, func() (bool, bool) {
		return respValue.OK, respErr.StatusCode != 0
	})
	if err != nil {
		return nil, joinErrors(err, respErr)
	}
	//process result
	return respValue.Result.Coins, nil
}

func (api *API) restCoins() ([]*CoinResult, error) {
	type respDirectCoins struct {
		Result []string
	}
	//request
	res, err := api.client.rest.R().Get("/coins")
	if err = processConnectionError(res, err); err != nil {
		return nil, err
	}
	//json decode
	respValue := respDirectCoins{}
	err = universalJSONDecode(res.Body(), &respValue, nil, func() (bool, bool) {
		return len(respValue.Result) > 0, false
	})
	if err != nil {
		return nil, err
	}
	//process result
	coins := []*CoinResult{}
	errstr := ""
	for _, val := range respValue.Result {
		coin, err := api.Coin(val)
		if err != nil {
			errstr += err.Error()
			continue
		}
		coins = append(coins, coin)
	}
	err = nil
	if errstr != "" {
		err = errors.New(errstr)
	}
	return coins, err
}

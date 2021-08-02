package api

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/auth"
)

// AddressResult contains API response fields.
type AddressResult struct {
	ID      uint64            `json:"id"`
	Address string            `json:"address"`
	Nonce   uint64            `json:"nonce"`
	Balance map[string]string `json:"balance"`
}

// Address requests full information about specified address.
func (api *API) Address(address string) (*AddressResult, error) {
	type respAddress struct {
		OK     bool `json:"ok"`
		Result struct {
			Address *AddressResult `json:"address"`
			Coins   []*CoinResult  `json:"coins"`
		} `json:"result"`
	}
	type respDirectAddress struct {
		Result struct {
			Value struct {
				AccountNumber uint64 `json:"account_number"`
				Address       string `json:"address"`
				Sequence      uint64 `json:"sequence"`
				Coins         []struct {
					Denom  string `json:"denom"`
					Amount string `json:"amount"`
				} `json:"coins"`
			} `json:"value"`
		} `json:"result"`
	}

	var (
		gres = respAddress{}
		dres = respDirectAddress{}
		url  = ""
	)

	if api.directConn == nil {
		url = fmt.Sprintf("/address/%s", address)
	} else {
		url = fmt.Sprintf("/auth/accounts/%s", address)
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
		return gres.Result.Address, nil
	}

	balance := make(map[string]string)
	for _, val := range dres.Result.Value.Coins {
		balance[val.Denom] = val.Amount
	}

	return &AddressResult{
		ID:      dres.Result.Value.AccountNumber,
		Address: dres.Result.Value.Address,
		Nonce:   dres.Result.Value.Sequence,
		Balance: balance,
	}, nil
}

// AccountNumberAndSequence requests account number and current sequence (nonce) of specified address.
func (api *API) AccountNumberAndSequence(address string) (uint64, uint64, error) {
	var (
		url = ""
	)

	if api.directConn == nil {
		url = fmt.Sprintf("/rpc/auth/accounts/%s", address)
	} else {
		url = fmt.Sprintf("/auth/accounts/%s", address)
	}

	res, err := api.client.rest.R().Get(url)
	if err != nil {
		return 0, 0, err
	}
	if res.IsError() {
		return 0, 0, NewResponseError(res)
	}

	response := struct {
		Height string          `json:"height"`
		Result json.RawMessage `json:"result"`
	}{}
	err = json.Unmarshal(res.Body(), &response)
	if err != nil {
		return 0, 0, err
	}

	account := auth.BaseAccount{}
	err = api.codec.UnmarshalJSON(response.Result, &account)
	if err != nil {
		return 0, 0, err
	}

	return account.AccountNumber, account.Sequence, nil
}

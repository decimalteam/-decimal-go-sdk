package api

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/auth"
)

// AddressResponse contains API response.
type AddressResponse struct {
	OK     bool `json:"ok"`
	Result struct {
		Address *AddressResult `json:"address"`
		Coins   []*CoinResult  `json:"coins"`
	} `json:"result"`
}

// AddressResult contains API response fields.
type AddressResult struct {
	ID         uint64            `json:"id"`
	Address    string            `json:"address"`
	Type       string            `json:"type"`
	Nonce      string            `json:"nonce"`
	Balance    map[string]string `json:"balance"`
	BalanceNft []string          `json:"balanceNft"`
	Txes       uint64            `json:"txes"`
}

// Address requests full information about specified address.
func (api *API) Address(address string) (*AddressResult, error) {

	url := fmt.Sprintf("/address/%s", address)
	res, err := api.client.R().Get(url)
	if err != nil {
		return nil, err
	}
	if res.IsError() {
		return nil, NewResponseError(res)
	}

	response := AddressResponse{}
	err = json.Unmarshal(res.Body(), &response)
	fmt.Printf("1111: %s\n", err)
	if err != nil || !response.OK {
		responseError := Error{}
		err = json.Unmarshal(res.Body(), &responseError)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("received response containing error: %s", responseError.Error())
	}

	return response.Result.Address, nil
}

// AccountNumberAndSequence requests account number and current sequence (nonce) of specified address.
func (api *API) AccountNumberAndSequence(address string) (uint64, uint64, error) {

	url := fmt.Sprintf("/rpc/auth/accounts/%s", address)
	res, err := api.client.R().Get(url)
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

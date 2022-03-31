package api

import (
	"fmt"
	"strconv"
)

// AddressResult contains API response fields.
type AddressResult struct {
	ID         uint64              `json:"id"`
	Address    string              `json:"address"`
	Type       string              `json:"type"`
	Nonce      string              `json:"nonce"`
	Balance    map[string]string   `json:"balance"`
	BalanceNft []*BalanceNftResult `json:"balanceNft"`
	Txes       uint64              `json:"txes"`
}

type BalanceNftResult struct {
	NftId      string                     `json:"nftId"`
	Collection string                     `json:"collection"`
	Amount     string                     `json:"amount"`
	NftReserve []*BalanceNftReserveResult `json:"nftReserve"`
	NftStake   []*BalanceNftStakeResult   `json:"nftStake"`
}

type BalanceNftReserveResult struct {
	SubTokenId  string `json:"subTokenId"`
	Reserve     string `json:"reserve"`
	Address     string `json:"address"`
	Delegated   bool   `json:"delegated"`
	ValidatorId string `json:"validatorId"`
	Unbonded    bool   `json:"unbonded"`
}

type BalanceNftStakeResult struct {
	NftId          string `json:"nftId"`
	SubTokenId     string `json:"subTokenId"`
	BaseQuantity   string `json:"baseQuantity"`
	UnbondQuantity string `json:"unbondQuantity"`
}

// Address requests full information about specified address.
// TODO: NFT for direct connection
// Gateway: ok, REST/RPC: partial
func (api *API) Address(address string) (*AddressResult, error) {
	if api.directConn == nil {
		return api.apiAddress(address)
	} else {
		return api.restAddress(address)
	}
}

func (api *API) apiAddress(address string) (*AddressResult, error) {
	type respAddress struct {
		OK     bool `json:"ok"`
		Result struct {
			Address *AddressResult `json:"address"`
			Coins   []*CoinResult  `json:"coins"`
		} `json:"result"`
	}
	//request
	res, err := api.client.rest.R().Get(fmt.Sprintf("/address/%s", address))
	if err = processConnectionError(res, err); err != nil {
		return nil, err
	}
	//json decode
	respValue, respErr := respAddress{}, Error{}
	err = universalJSONDecode(res.Body(), &respValue, &respErr, func() (bool, bool) {
		return respValue.OK, respErr.StatusCode != 0
	})
	if err != nil {
		return nil, joinErrors(err, respErr)
	}
	//process result
	return respValue.Result.Address, nil
}

func (api *API) restAddress(address string) (*AddressResult, error) {
	type respDirectAddress struct {
		Result struct {
			Value struct {
				AccountNumber string `json:"account_number"`
				Address       string `json:"address"`
				Sequence      string `json:"sequence"`
				Coins         []struct {
					Denom  string `json:"denom"`
					Amount string `json:"amount"`
				} `json:"coins"`
			} `json:"value"`
		} `json:"result"`
	}
	//request
	res, err := api.client.rest.R().Get(fmt.Sprintf("/auth/accounts/%s", address))
	if err = processConnectionError(res, err); err != nil {
		return nil, err
	}
	//json decode
	respValue, respErr := respDirectAddress{}, Error{}
	err = universalJSONDecode(res.Body(), &respValue, &respErr, func() (bool, bool) {
		return respValue.Result.Value.AccountNumber > "", respErr.StatusCode != 0
	})
	if err != nil {
		return nil, joinErrors(err, respErr)
	}
	//process result
	balance := make(map[string]string)
	for _, val := range respValue.Result.Value.Coins {
		balance[val.Denom] = val.Amount
	}

	accNumber, _ := strconv.ParseUint(respValue.Result.Value.AccountNumber, 10, 64)

	return &AddressResult{
		ID:      accNumber,
		Address: respValue.Result.Value.Address,
		Nonce:   respValue.Result.Value.Sequence,
		Balance: balance,
	}, nil
}

// AccountNumberAndSequence requests account number and current sequence (nonce) of specified address.
// Gateway: ok, REST/RPC: ok
func (api *API) AccountNumberAndSequence(address string) (uint64, uint64, error) {
	adrRes, err := api.Address(address)
	if err != nil {
		return 0, 0, err
	}
	seq, _ := strconv.ParseUint(adrRes.Nonce, 10, 64)
	return adrRes.ID, seq, nil
}

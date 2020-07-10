package api

import (
	"encoding/json"
	"fmt"
)

// StakesResponse contains API response.
type StakesResponse struct {
	OK     bool `json:"ok"`
	Result struct {
		Stakes []*StakesResult `json:"validators"`
		Total  string          `json:"total"`
	} `json:"result"`
}

// StakesResult contains API response fields.
type StakesResult struct {
	ValidatorID uint64           `json:"validatorId"`
	TotalStake  string           `json:"totalStake"`
	Validator   *ValidatorResult `json:"validator"`
	Stakes      []*Stake         `json:"stakes"`
}

// Stake contains API response fields.
type Stake struct {
	Coin         string `json:"coin"`
	Amount       string `json:"amount"`
	BaseAmount   string `json:"baseAmount"`
	UnbondAmount string `json:"unbondAmount"`
}

// Stakes requests full information about stakes from the account with specified address.
func (api *API) Stakes(address string) ([]*StakesResult, error) {

	url := fmt.Sprintf("/address/%s/stakes", address)
	res, err := api.client.R().Get(url)
	if err != nil {
		return nil, err
	}
	if res.IsError() {
		return nil, NewResponseError(res)
	}

	response := StakesResponse{}
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
		return nil, fmt.Errorf("received response containing error: %s", responseError.Error())
	}

	return response.Result.Stakes, nil
}

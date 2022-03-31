package api

import (
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
	ValidatorID string           `json:"validatorId"`
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
// TODO: implement this API in REST/RPC.
func (api *API) Stakes(address string) ([]*StakesResult, error) {
	if api.directConn != nil {
		return nil, ErrNotImplemented
	}
	//request
	res, err := api.client.rest.R().Get(fmt.Sprintf("/address/%s/stakes", address))
	if err = processConnectionError(res, err); err != nil {
		return nil, err
	}
	//json decode
	respValue, respErr := StakesResponse{}, Error{}
	err = universalJSONDecode(res.Body(), &respValue, &respErr, func() (bool, bool) {
		return respValue.OK && len(respValue.Result.Stakes) > 0, respErr.StatusCode != 0
	})
	return nil, respErr
	if err != nil {
		return nil, joinErrors(err, respErr)
	}
	//process result
	return respValue.Result.Stakes, nil
}

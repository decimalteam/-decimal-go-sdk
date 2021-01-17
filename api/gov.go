package api

import (
	"encoding/json"
	"fmt"

	"bitbucket.org/decimalteam/go-node/x/gov"
)

// GovResponse contains API response.
type GovResponse struct {
	OK     bool      `json:"ok"`
	Result GovResult `json:"result"`
}

// GovsResponse contains API response.
type GovsResponse struct {
	OK     bool `json:"ok"`
	Result struct {
		Count uint64      `json:"count"`
		Govs  []GovResult `json:"coins"`
	} `json:"result"`
}

// CoinResult contains API response fields.
type GovResult struct {
	gov.Proposal
}

// Govs requests full information about all govs.
func (api *API) Govs() ([]GovResult, error) {
	url := "/govs"

	res, err := api.client.R().Get(url)
	if err != nil {
		return nil, err
	}
	if res.IsError() {
		return nil, NewResponseError(res)
	}

	response := GovsResponse{}
	err = json.Unmarshal(res.Body(), &response)
	if err != nil || !response.OK {
		responseError := Error{}
		err = json.Unmarshal(res.Body(), &responseError)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("received response containing error: %s", responseError.Error())
	}

	return response.Result.Govs, nil
}

// Gov requests full information about gov with specified id.
func (api *API) Gov(id uint64) (GovResult, error) {

	url := fmt.Sprintf("/gov/%d", id)
	res, err := api.client.R().Get(url)
	if err != nil {
		return GovResult{}, err
	}
	if res.IsError() {
		return GovResult{}, NewResponseError(res)
	}

	response := GovResponse{}
	err = json.Unmarshal(res.Body(), &response)
	if err != nil || !response.OK {
		responseError := Error{}
		err = json.Unmarshal(res.Body(), &responseError)
		if err != nil {
			return GovResult{}, err
		}
		return GovResult{}, fmt.Errorf("received response containing error: %s", responseError.Error())
	}

	return response.Result, nil
}

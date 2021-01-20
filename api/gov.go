package api

import (
	"encoding/json"
	"fmt"
	"time"
)

// ProposalResponse contains API response.
type ProposalResponse struct {
	OK     bool           `json:"ok"`
	Result ProposalResult `json:"result"`
}

// ProposalsResponse contains API response.
type ProposalsResponse struct {
	OK     bool `json:"ok"`
	Result struct {
		Count     uint64           `json:"count"`
		Proposals []ProposalResult `json:"proposals"`
	} `json:"result"`
}

// CoinResult contains API response fields.
type ProposalResult struct {
	ProposalID       int64     `json:"proposalId"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	VotingStartBlock string    `json:"votingStartBlock"`
	VotingEndBlock   string    `json:"votingEndBlock"`
	Proposer         string    `json:"proposer"`
	StakesTotal      float64   `json:"stakesTotal"`
	StakesYes        float64   `json:"stakesYes"`
	StakesNo         float64   `json:"stakesNo"`
	StakesAbstain    float64   `json:"stakesAbstain"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
	PercentYes       string    `json:"percentYes"`
	PercentNo        string    `json:"percentNo"`
	PercentAbstain   string    `json:"percentAbstain"`
	PercentNone      string    `json:"percentNone"`
	Votes            struct {
		Count int `json:"count"`
		Votes []struct {
			ID          int64     `json:"id"`
			HashTx      string    `json:"hashTx"` // Hash of transaction in which the proposal was created
			ValidatorId string    `json:"validatorId"`
			Stake       float64   `json:"stake"`
			Vote        string    `json:"vote"`
			CreatedAt   time.Time `json:"createdAt"`
			UpdatedAt   time.Time `json:"updatedAt"`
		}
	} `json:"votes"`

	HashTx string `json:"hashTx"` // Hash of transaction in which the proposal was created
}

// Proposals requests full information about all govs.
func (api *API) Proposals() ([]ProposalResult, error) {
	url := "/proposals"

	res, err := api.client.R().Get(url)
	if err != nil {
		return nil, err
	}
	if res.IsError() {
		return nil, NewResponseError(res)
	}

	response := ProposalsResponse{}
	err = json.Unmarshal(res.Body(), &response)
	if err != nil || !response.OK {
		fmt.Println(err)
		responseError := Error{}
		err = json.Unmarshal(res.Body(), &responseError)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("received response containing error: %s", responseError.Error())
	}

	return response.Result.Proposals, nil
}

// Proposal requests full information about gov with specified id.
func (api *API) Proposal(id int64) (ProposalResult, error) {

	url := fmt.Sprintf("/proposalById/%d", id)
	res, err := api.client.R().Get(url)
	if err != nil {
		return ProposalResult{}, err
	}
	if res.IsError() {
		return ProposalResult{}, NewResponseError(res)
	}

	response := ProposalResponse{}
	err = json.Unmarshal(res.Body(), &response)
	if err != nil || !response.OK {
		fmt.Println(err)
		responseError := Error{}
		err = json.Unmarshal(res.Body(), &responseError)
		if err != nil {
			return ProposalResult{}, err
		}
		return ProposalResult{}, fmt.Errorf("received response containing error: %s", responseError.Error())
	}

	return response.Result, nil
}

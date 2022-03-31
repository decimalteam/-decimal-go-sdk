package api

import (
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
// Gateway: ok, RPC/REST: none
func (api *API) Proposals() ([]ProposalResult, error) {
	var (
		url = ""
	)

	// TODO: test with directConn.
	if api.directConn == nil {
		url = "/proposals"
	} else {
		return nil, ErrNotImplemented
		url = "/gov/proposals"
	}
	//request
	res, err := api.client.rest.R().Get(url)
	if err = processConnectionError(res, err); err != nil {
		return nil, err
	}
	//json decode
	respValue, respErr := ProposalsResponse{}, Error{}
	err = universalJSONDecode(res.Body(), &respValue, &respErr, func() (bool, bool) {
		return respValue.OK, respErr.StatusCode != 0
	})
	if err != nil {
		return nil, joinErrors(err, respErr)
	}
	//process result
	return respValue.Result.Proposals, nil
}

// Proposal requests full information about gov with specified id.
// Gateway: ok, RPC/REST: none
func (api *API) Proposal(id int64) (ProposalResult, error) {
	var (
		url = ""
	)

	// TODO: test with directConn.
	if api.directConn == nil {
		url = fmt.Sprintf("/proposalById/%d", id)
	} else {
		return ProposalResult{}, ErrNotImplemented
		url = fmt.Sprintf("/gov/proposals/%d", id)
	}

	//request
	res, err := api.client.rest.R().Get(url)
	if err = processConnectionError(res, err); err != nil {
		return ProposalResult{}, err
	}
	//json decode
	respValue, respErr := ProposalResponse{}, Error{}
	err = universalJSONDecode(res.Body(), &respValue, &respErr, func() (bool, bool) {
		return respValue.OK, respErr.StatusCode != 0
	})
	if err != nil {
		return ProposalResult{}, joinErrors(err, respErr)
	}
	// process result
	return respValue.Result, nil
}

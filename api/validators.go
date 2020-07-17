package api

import (
	"encoding/json"
	"fmt"
)

// ValidatorResponse contains API response.
type ValidatorResponse struct {
	OK     bool             `json:"ok"`
	Result *ValidatorResult `json:"result"`
}

// ValidatorsResponse contains API response.
type ValidatorsResponse struct {
	OK     bool `json:"ok"`
	Result struct {
		Count      int                `json:"count"`
		Online     int                `json:"online"`
		Validators []*ValidatorResult `json:"validators"`
	} `json:"result"`
}

// ValidatorResult contains API response fields.
type ValidatorResult struct {
	Address          string `json:"address"`          // Address with prefix "dxvaloper" to manage the validator
	RewardAddress    string `json:"rewardAddress"`    // Address with prefix "dx" to receive rewards for participating block producing and consensus
	ConsensusAddress string `json:"consensusAddress"` // Address with prefix "dxvalcons" used only for consensus

	Moniker         string `json:"moniker"`          // Specified by the validator operator
	Website         string `json:"website"`          // Specified by the validator operator
	Details         string `json:"details"`          // Specified by the validator operator
	Identity        string `json:"identity"`         // Specified by the validator operator
	SecurityContact string `json:"security_contact"` // Specified by the validator operator
	Comission       string `json:"fee"`              // Specified by the validator operator

	BlockID       uint64 `json:"blockId"`       // Number of block in which the validator was declared
	SkippedBlocks uint64 `json:"skippedBlocks"` // Amount of blocks missed to sign
	Stake         string `json:"stake"`         // Total stake of the validator
	Power         string `json:"power"`         // Voting power of the validator
	Slots         uint64 `json:"slots"`         // Amount of delegator slots used
	MinStake      string `json:"mins"`          // Minimum stake needed to get place in the delegators list
	Rating        string `json:"rating"`        // Rating of the validator
	Status        string `json:"status"`        // Current status of the validator (online, offline)
	Kind          string `json:"kind"`          // Kind of the validator (Validator, Candidate)
}

// Candidates requests full list of candidates (validators which does not participate in block production and consensus).
func (api *API) Candidates() ([]*ValidatorResult, error) {

	url := "/validators/candidate"
	res, err := api.client.R().Get(url)
	if err != nil {
		return nil, err
	}
	if res.IsError() {
		return nil, NewResponseError(res)
	}

	response := ValidatorsResponse{}
	err = json.Unmarshal(res.Body(), &response)
	if err != nil || !response.OK {
		responseError := Error{}
		err = json.Unmarshal(res.Body(), &responseError)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("received response containing error: %s", responseError.Error())
	}

	return response.Result.Validators, nil
}

// Validators requests full list of currently active validators.
func (api *API) Validators() ([]*ValidatorResult, error) {

	url := "/validators/validator"
	res, err := api.client.R().Get(url)
	if err != nil {
		return nil, err
	}
	if res.IsError() {
		return nil, NewResponseError(res)
	}

	response := ValidatorsResponse{}
	err = json.Unmarshal(res.Body(), &response)
	if err != nil || !response.OK {
		responseError := Error{}
		err = json.Unmarshal(res.Body(), &responseError)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("received response containing error: %s", responseError.Error())
	}

	return response.Result.Validators, nil
}

// Validator requests full information about validator with specified address.
func (api *API) Validator(address string) (*ValidatorResult, error) {

	url := fmt.Sprintf("/validator/%s", address)
	res, err := api.client.R().Get(url)
	if err != nil {
		return nil, err
	}
	if res.IsError() {
		return nil, NewResponseError(res)
	}

	response := ValidatorResponse{}
	err = json.Unmarshal(res.Body(), &response)
	if err != nil || !response.OK {
		responseError := Error{}
		err = json.Unmarshal(res.Body(), &responseError)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("received response containing error: %s", responseError.Error())
	}

	return response.Result, nil
}

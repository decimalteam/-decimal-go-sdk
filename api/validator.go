package api

import (
	"encoding/json"
	"fmt"
)

// ValidatorResponse contains API response.
type ValidatorResponse struct {
	OK     bool             `json:"ok,omitempty"`
	Result *ValidatorResult `json:"result,omitempty"`
}

// ValidatorResult contains API response fields.
type ValidatorResult struct {
	Address          string `json:"address"`          // Address with prefix "dxvaloper" to manage validator
	RewardAddress    string `json:"rewardAddress"`    // Address with prefix "dx" to receive reward
	ConsensusAddress string `json:"consensusAddress"` // Address with prefix "dxvalcons" used for the consensus
	Moniker          string `json:"moniker"`          // Specified by validator operator
	Website          string `json:"website"`          // Specified by validator operator
	Details          string `json:"details"`          // Specified by validator operator
	Identity         string `json:"identity"`         // Specified by validator operator
	SecurityContact  string `json:"security_contact"` // Specified by validator operator
	Fee              string `json:"fee"`              // Specified by validator operator
	BlockID          uint64 `json:"blockId"`          // Number of block in which the validator was declared
	SkippedBlocks    uint64 `json:"skippedBlocks"`    // Amount of blocks missed to sign
	Slots            uint64 `json:"slots"`            // Amount of delegator slots used
	Stake            string `json:"stake"`            // Total stake of the validator
	MinStake         string `json:"mins"`             // Minimum stake needed to get place in the delegators list
	Power            string `json:"power"`            // Voting power of the validator
	Rating           string `json:"rating"`           // Rating of the validator
	Status           string `json:"status"`           // Current status of validator (online, offline)
	Kind             string `json:"kind"`             // Kind of validator (Validator, Candidate)
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
	if err != nil {
		return nil, err
	}

	if !response.OK {
		responseError := Error{}
		err = json.Unmarshal(res.Body(), &responseError)
		if err != nil {
			return nil, err
		}
		// TODO
	}

	return response.Result, nil
}

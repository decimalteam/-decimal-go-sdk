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
	Stake           string `json:"stake"`            // Total stake of the validator
	Power           string `json:"power"`            // Voting power of the validator (propotional to stake)

	// TODO: this info missing in REST.
	BlockID       uint64 `json:"blockId"`       // Number of block in which the validator was declared
	SkippedBlocks uint64 `json:"skippedBlocks"` // Amount of blocks missed to sign
	Slots         uint64 `json:"slots"`         // Amount of delegator slots used
	MinStake      string `json:"mins"`          // Minimum stake needed to get place in the delegators list
	Rating        string `json:"rating"`        // Rating of the validator
	Status        string `json:"status"`        // Current status of the validator (online, offline)
	Kind          string `json:"kind"`          // Kind of the validator (Validator, Candidate)
}

type respDirectValidator struct {
	Address          string `json:"val_address"`    // Address with prefix "dxvaloper" to manage the validator
	RewardAddress    string `json:"reward_address"` // Address with prefix "dx" to receive rewards for participating block producing and consensus
	ConsensusAddress string `json:"pub_key"`        // Address with prefix "dxvalcons" used only for consensus
	Stake            string `json:"stake_coins"`    // Total stake of the validator
	Status           int    `json:"status"`         // TODO: kind of validator?
	Online           bool   `json:"online"`         // online/offline
	// TODO: comMission and comission?
	Comission   string `json:"commission"` // Specified by the validator operator
	Description struct {
		Moniker         string `json:"moniker"`          // Specified by the validator operator
		Identity        string `json:"identity"`         // Specified by the validator operator
		Website         string `json:"website"`          // Specified by the validator operator
		SecurityContact string `json:"security_contact"` // Specified by the validator operator
		Details         string `json:"details"`          // Specified by the validator operator
	} `json:"description"`
}

type respDirectValidators struct {
	Result []respDirectValidator `json:"result"`
}

func directResonse2Validator(dres respDirectValidator) *ValidatorResult {
	power := "0"
	if len(dres.Stake) > 18 {
		power = dres.Stake[0 : len(dres.Stake)-18]
	}
	onlineStatus := "online"
	if !dres.Online {
		onlineStatus = "offline"
	}
	return &ValidatorResult{
		Address:          dres.Address,
		RewardAddress:    dres.RewardAddress,
		ConsensusAddress: dres.ConsensusAddress,
		Comission:        dres.Comission,
		Moniker:          dres.Description.Moniker,
		Website:          dres.Description.Website,
		Details:          dres.Description.Details,
		Identity:         dres.Description.Identity,
		SecurityContact:  dres.Description.SecurityContact,
		Stake:            dres.Stake,
		Power:            power,
		Status:           onlineStatus,
		Kind:             "Validator", // TODO: replace after releasing candidates
	}

}

// Candidates requests full list of candidates (validators which does not participate in block production and consensus).
// Gateway: ok, REST/RPC: none
func (api *API) Candidates() ([]*ValidatorResult, error) {

	// TODO: undefined candidates in REST.
	url := "/validators/candidate"
	res, err := api.client.rest.R().Get(url)
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
// Gateway: ok, REST/RPC: partial
func (api *API) Validators() ([]*ValidatorResult, error) {
	var (
		url  = ""
		gres = ValidatorsResponse{}
		dres = respDirectValidators{}
	)

	if api.directConn == nil {
		url = "/validators/validator"
	} else {
		url = "/validator/validators"
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
		return gres.Result.Validators, nil
	}

	validators := []*ValidatorResult{}
	for _, val := range dres.Result {
		validators = append(validators, directResonse2Validator(val))
	}

	return validators, nil
}

// Validator requests full information about validator with specified address.
// Gateway: ok, REST/RPC: partial
func (api *API) Validator(address string) (*ValidatorResult, error) {
	type respDirect struct {
		Result respDirectValidator `json:"result"`
	}

	var (
		url  = ""
		gres = ValidatorResponse{}
		dres = respDirect{}
	)

	if api.directConn == nil {
		url = fmt.Sprintf("/validator/%s", address)
	} else {
		url = fmt.Sprintf("/validator/validators/%s", address)
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

	err = json.Unmarshal(res.Body(), &gres)
	if err != nil || (api.directConn == nil && !gres.OK) {
		responseError := Error{}
		err = json.Unmarshal(res.Body(), &responseError)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("received response containing error: %s", responseError.Error())
	}

	if api.directConn == nil {
		return gres.Result, nil
	}

	return directResonse2Validator(dres.Result), nil
}

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

	// TODO: this info missing in REST.
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
func (api *API) Validators() ([]*ValidatorResult, error) {
	type respDirectValidators struct {
		Result []struct {
			Address          string `json:"val_address"`    // Address with prefix "dxvaloper" to manage the validator
			RewardAddress    string `json:"reward_address"` // Address with prefix "dx" to receive rewards for participating block producing and consensus
			ConsensusAddress string `json:"pub_key"`        // Address with prefix "dxvalcons" used only for consensus
			// TODO: comMission and comission?
			Comission   string `json:"commission"` // Specified by the validator operator
			Description struct {
				Moniker         string `json:"moniker"`          // Specified by the validator operator
				Identity        string `json:"identity"`         // Specified by the validator operator
				Website         string `json:"website"`          // Specified by the validator operator
				SecurityContact string `json:"security_contact"` // Specified by the validator operator
				Details         string `json:"details"`          // Specified by the validator operator
			} `json:"description"`
		} `json:"result"`
	}

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
		validators = append(validators, &ValidatorResult{
			Address:          val.Address,
			RewardAddress:    val.RewardAddress,
			ConsensusAddress: val.ConsensusAddress,
			Comission:        val.Comission,
			Moniker:          val.Description.Moniker,
			Website:          val.Description.Website,
			Details:          val.Description.Details,
			Identity:         val.Description.Identity,
			SecurityContact:  val.Description.SecurityContact,
		})
	}

	return validators, nil
}

// Validator requests full information about validator with specified address.
func (api *API) Validator(address string) (*ValidatorResult, error) {
	type respDirectValidator struct {
		Result struct {
			Address          string `json:"val_address"`    // Address with prefix "dxvaloper" to manage the validator
			RewardAddress    string `json:"reward_address"` // Address with prefix "dx" to receive rewards for participating block producing and consensus
			ConsensusAddress string `json:"pub_key"`        // Address with prefix "dxvalcons" used only for consensus
			// TODO: comMission and comission?
			Comission   string `json:"commission"` // Specified by the validator operator
			Description struct {
				Moniker         string `json:"moniker"`          // Specified by the validator operator
				Identity        string `json:"identity"`         // Specified by the validator operator
				Website         string `json:"website"`          // Specified by the validator operator
				SecurityContact string `json:"security_contact"` // Specified by the validator operator
				Details         string `json:"details"`          // Specified by the validator operator
			} `json:"description"`
		} `json:"result"`
	}

	var (
		url  = ""
		gres = ValidatorResponse{}
		dres = respDirectValidator{}
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

	validator := &ValidatorResult{
		Address:          dres.Result.Address,
		RewardAddress:    dres.Result.RewardAddress,
		ConsensusAddress: dres.Result.ConsensusAddress,
		Comission:        dres.Result.Comission,
		Moniker:          dres.Result.Description.Moniker,
		Website:          dres.Result.Description.Website,
		Details:          dres.Result.Description.Details,
		Identity:         dres.Result.Description.Identity,
		SecurityContact:  dres.Result.Description.SecurityContact,
	}

	return validator, nil
}

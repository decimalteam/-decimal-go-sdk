package api

import (
	"encoding/json"
	"fmt"
)

type NFT struct {
	ID          int         `json:"id"`
	Slug        string      `json:"slug"`
	Headline    string      `json:"headline"`
	Description string      `json:"description"`
	Misc        interface{} `json:"misc"`
	Cover       string      `json:"cover"`
	Asset       string      `json:"asset"`
	Status      string      `json:"status"`
	CreatedAt   string      `json:"createdAt"`
	UpdatedAt   string      `json:"updatedAt"`
}

type NFTResponse struct {
	OK     bool `json:"ok"`
	Result struct {
		Token *NFT
	}
}

func (api *API) NFT(id string) (*NFT, error) {
	var (
		url = ""
	)

	// TODO: test with directConn.
	if api.directConn == nil {
		url = fmt.Sprintf("/nfts/%s", id)
	} else {
		url = fmt.Sprintf("/nft/collection/%s", id)
	}

	res, err := api.client.rest.R().Get(url)
	if err != nil {
		return nil, err
	}
	if res.IsError() {
		return nil, NewResponseError(res)
	}

	response := NFTResponse{}
	err = json.Unmarshal(res.Body(), &response)
	if err != nil || !response.OK {
		responseError := Error{}
		err = json.Unmarshal(res.Body(), &responseError)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("received response containing error: %s", responseError.Error())
	}

	return response.Result.Token, nil
}

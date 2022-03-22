package api

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Full NTF info
// TODO: replace replace small integers, serialized as string, to int64
type NFT struct {
	Id             string     `json:"nftId"`
	CollectionName string     `json:"nftCollection"`
	Quantity       string     `json:"quantity"`
	StartReserve   string     `json:"startReserve"`
	TotalReserve   string     `json:"totalReserve"`
	TokenURI       string     `json:"tokenUri"`
	AllowMint      bool       `json:"allowMint"`
	NonFungible    bool       `json:"nonFungible"`
	TxHash         string     `json:"txHash"`
	BlockId        int64      `json:"blockId"`
	CreatedAt      string     `json:"createdAt"`
	UpdatedAt      string     `json:"updatedAt"`
	Slug           string     `json:"slug"`
	Headline       string     `json:"headline"`
	Description    string     `json:"description"`
	NFTOwner       []NFTOwner `json:"nftOwner"`
	NFTReserve     []struct {
		SubTokenId  string `json:"subTokenId"`
		Reserve     string `json:"reserve"`
		Address     string `json:"address"`
		Delegated   bool   `json:"delegated"`
		ValidatorId string `json:"validatorId"`
	} `json:"nftReserve"`
	Misc struct {
		CoverHash      string `json:"coverHash"`
		CoverPath      string `json:"coverPath"`
		CoverExtension string `json:"coverExtension"`
	} `json:"misc"`
	Cover     string `json:"cover"`
	Status    string `json:"status"`
	IsPrivate bool   `json:"isPrivate"`
}

type NFTOwner struct {
	Address  string `json:"address"`
	Quantity string `json:"quantity"`
}

type NFTShort struct {
	Id             string
	CollectionName string
	Creator        string
	CreatedAt      string
	Quantity       int64
	Delegated      int64
	Unbound        int64
}

// Get all NFT in short format
func (api *API) NFTList() ([]*NFTShort, error) {
	if api.directConn == nil {
		data, err := api.apiNFTList()
		return data, err
	} else {
		data, err := api.restNFTList()
		return data, err
	}
}

func (api *API) apiNFTList() ([]*NFTShort, error) {
	var result []*NFTShort
	type responseType struct {
		OK     bool `json:"ok"`
		Result []struct {
			RawId             string `json:"nftId"`
			RawCollectionName string `json:"nftCollection"`
			RawCreator        string `json:"creator"`
			RawCreatedAt      string `json:"createdAt"`
			RawQuantity       string `json:"quantity"`
			RawDelegated      string `json:"delegated"`
			RawUnbound        string `json:"unbound"`
		} `json:"result"`
	}
	res, err := api.client.rest.R().Get("/nfts")
	// initial errors: connection/404/...
	if err != nil {
		return []*NFTShort{}, err
	}
	if res.IsError() {
		return []*NFTShort{}, NewResponseError(res)
	}
	response := responseType{}
	err = json.Unmarshal(res.Body(), &response)
	// errors in respose
	if err != nil || !response.OK {
		responseError := Error{}
		err = json.Unmarshal(res.Body(), &responseError)
		if err != nil {
			return []*NFTShort{}, err
		}
		return []*NFTShort{}, fmt.Errorf("received response containing error: %s", responseError.Error())
	}
	//
	err = nil
	for _, rec := range response.Result {
		q, err := strconv.ParseInt(rec.RawQuantity, 10, 64)
		if err != nil {
			return []*NFTShort{}, fmt.Errorf("Cannot convert Quantity '%s' to int", rec.RawQuantity)
		}
		d, err := strconv.ParseInt(rec.RawDelegated, 10, 64)
		if err != nil {
			return []*NFTShort{}, fmt.Errorf("Cannot convert Delegated '%s' to int", rec.RawDelegated)
		}
		u, err := strconv.ParseInt(rec.RawUnbound, 10, 64)
		if err != nil {
			return []*NFTShort{}, fmt.Errorf("Cannot convert Unbound '%s' to int", rec.RawUnbound)
		}
		result = append(result, &NFTShort{
			Id:             rec.RawId,
			CollectionName: rec.RawCollectionName,
			Creator:        rec.RawCreator,
			CreatedAt:      rec.RawCreatedAt,
			Quantity:       q,
			Delegated:      d,
			Unbound:        u,
		})
	}
	return result, nil
}

func (api *API) restNFTList() ([]*NFTShort, error) {
	var result []*NFTShort
	type responseDenomsType struct {
		Result []string `json:"result"`
	}
	type responseNFTType struct {
		Result map[string]struct {
			NFTs map[string]struct {
				RawId      string `json:"id"`
				RawCreator string `json:"creator"`
			} `json:"nfts"`
		} `json:"result"`
	}
	// 1 get collection list
	res, err := api.client.rest.R().Get("/nft/denoms")
	if err != nil {
		return []*NFTShort{}, err
	}
	response := responseDenomsType{}
	err = json.Unmarshal(res.Body(), &response)
	if err != nil {
		return []*NFTShort{}, err
	}
	// 2 get nfts from collections
	for _, denom := range response.Result {
		res, err := api.client.rest.R().Get(fmt.Sprintf("/nft/collection/%s", denom))
		if err != nil {
			return []*NFTShort{}, err
		}
		response := responseNFTType{}
		err = json.Unmarshal(res.Body(), &response)
		if err != nil {
			return []*NFTShort{}, err
		}
		for _, v1 := range response.Result {
			for _, v2 := range v1.NFTs {
				result = append(result, &NFTShort{
					Id:             v2.RawId,
					CollectionName: denom,
					Creator:        v2.RawCreator,
					CreatedAt:      "",
					Quantity:       0, //TODO
					Delegated:      0, //TODO
					Unbound:        0, //TODO
				})
			}
		}
	}

	return result, nil
}

// Get NFT owned by account
func (api *API) NFTByAddress(address string) ([]*NFT, error) {
	if api.directConn == nil {
		data, err := api.apiNFTByAddress(address)
		return data, err
	} else {
		data, err := api.restNFTByAddress(address)
		return data, err
	}
}

func (api *API) apiNFTByAddress(address string) ([]*NFT, error) {
	type responseType struct {
		OK     bool `json:"ok"`
		Result struct {
			Tokens []*NFT `json:"tokens"`
		}
	}
	res, err := api.client.rest.R().Get(fmt.Sprintf("/address/%s/nfts", address))
	// initial errors: connection/404/...
	if err != nil {
		return []*NFT{}, err
	}
	if res.IsError() {
		return []*NFT{}, NewResponseError(res)
	}
	response := responseType{}
	err = json.Unmarshal(res.Body(), &response)
	// errors in respose
	if err != nil || !response.OK {
		responseError := Error{}
		err = json.Unmarshal(res.Body(), &responseError)
		if err != nil {
			return []*NFT{}, err
		}
		return []*NFT{}, fmt.Errorf("received response containing error: %s", responseError.Error())
	}
	//
	return response.Result.Tokens, nil
}

//TODO: make full RPC/REST response
func (api *API) restNFTByAddress(address string) ([]*NFT, error) {
	type responseType struct {
		Result struct {
			Value struct {
				Collections []struct {
					Denom string   `json:"denom"`
					Ids   []string `json:"ids"`
				} `json:"idCollections"`
			} `json:"value"`
		} `json:"result"`
	}
	type responseNFTType struct {
		Result struct {
			Value struct {
				Owners struct {
					Owners []struct {
						Address string `json:"address"`
						//TODO "sub_token_ids": null
					} `json:"owners"`
				} `json:"owners"`
				Creator   string `json:"creator"`
				TokenURI  string `json:"token_uri"`
				Reserve   string `json:"reserve"`
				AllowMint bool   `json:"allow_mint"`
			} `json:"value"`
		} `json:"result"`
	}

	res, err := api.client.rest.R().Get(fmt.Sprintf("/nft/owner/%s", address))
	if err != nil {
		return []*NFT{}, err
	}
	response := responseType{}
	err = json.Unmarshal(res.Body(), &response)
	if err != nil {
		return []*NFT{}, err
	}
	result := []*NFT{}
	for _, col := range response.Result.Value.Collections {
		// TODO: add REST query by id
		for _, idd := range col.Ids {
			base := &NFT{Id: idd, CollectionName: col.Denom}
			res, err := api.client.rest.R().Get(fmt.Sprintf("/nft/collection/%s/nft/%s", col.Denom, idd))
			if err != nil {
				continue
			}
			response := responseNFTType{}
			err = json.Unmarshal(res.Body(), &response)
			if err != nil {
				continue
			}
			base.AllowMint = response.Result.Value.AllowMint
			base.TokenURI = response.Result.Value.TokenURI
			base.TotalReserve = response.Result.Value.Reserve
			for _, own := range response.Result.Value.Owners.Owners {
				base.NFTOwner = append(base.NFTOwner, NFTOwner{Address: own.Address})
			}
			result = append(result, base)
		}
	}
	return result, nil
}

// Get NTF by id
// TODO: implement for directConn.
func (api *API) NFT(id string) (*NFT, error) {
	type responseNFTType struct {
		OK     bool `json:"ok"`
		Result *NFT `json:"result"`
	}
	var url string

	if api.directConn == nil {
		url = fmt.Sprintf("/nfts/%s", id)
	} else {
		return nil, ErrNotImplemented
	}

	res, err := api.client.rest.R().Get(url)
	if err != nil {
		return nil, err
	}
	if res.IsError() {
		return nil, NewResponseError(res)
	}

	response := responseNFTType{}
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

package api

import (
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
		return api.apiNFTList()
	} else {
		return api.restNFTList()
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
	//request
	res, err := api.client.rest.R().Get("/nfts")
	if err = processConnectionError(res, err); err != nil {
		return nil, err
	}
	//json decode
	respValue, respErr := responseType{}, Error{}
	err = universalJSONDecode(res.Body(), &respValue, &respErr, func() (bool, bool) {
		return respValue.OK, respErr.StatusCode != 0
	})
	if err != nil {
		return nil, joinErrors(err, respErr)
	}
	//process
	err = nil
	for _, rec := range respValue.Result {
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
	//request
	res, err := api.client.rest.R().Get("/nft/denoms")
	if err = processConnectionError(res, err); err != nil {
		return []*NFTShort{}, err
	}
	//json decode
	response := responseDenomsType{}
	err = universalJSONDecode(res.Body(), &response, nil, func() (bool, bool) {
		return len(response.Result) > 0, false
	})
	if err != nil {
		return []*NFTShort{}, err
	}
	// 2 get nfts from collections
	for _, denom := range response.Result {
		//request
		res, err := api.client.rest.R().Get(fmt.Sprintf("/nft/collection/%s", denom))
		if err = processConnectionError(res, err); err != nil {
			return []*NFTShort{}, err
		}
		//json decode
		respValue := responseNFTType{}
		err = universalJSONDecode(res.Body(), &respValue, nil, func() (bool, bool) {
			return len(respValue.Result) > 0, false
		})
		if err != nil {
			return []*NFTShort{}, err
		}
		//process
		for _, v1 := range respValue.Result {
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
		return api.apiNFTByAddress(address)
	} else {
		return api.restNFTByAddress(address)
	}
}

func (api *API) apiNFTByAddress(address string) ([]*NFT, error) {
	type responseType struct {
		OK     bool `json:"ok"`
		Result struct {
			Tokens []*NFT `json:"tokens"`
		}
	}
	//request
	res, err := api.client.rest.R().Get(fmt.Sprintf("/address/%s/nfts", address))
	if err = processConnectionError(res, err); err != nil {
		return []*NFT{}, err
	}
	//json decode
	respValue, respErr := responseType{}, Error{}
	err = universalJSONDecode(res.Body(), &respValue, &respErr, func() (bool, bool) {
		return respValue.OK, respErr.StatusCode != 0
	})
	if err != nil {
		return nil, joinErrors(err, respErr)
	}
	//process result
	return respValue.Result.Tokens, nil
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
	if err = processConnectionError(res, err); err != nil {
		return []*NFT{}, err
	}
	response := responseType{}
	err = universalJSONDecode(res.Body(), &response, nil, func() (bool, bool) {
		return true, false
	})
	if err != nil {
		return []*NFT{}, err
	}
	result := []*NFT{}
	for _, col := range response.Result.Value.Collections {
		// TODO: add REST query by id
		for _, idd := range col.Ids {
			base := &NFT{Id: idd, CollectionName: col.Denom}
			res, err := api.client.rest.R().Get(fmt.Sprintf("/nft/collection/%s/nft/%s", col.Denom, idd))
			if err = processConnectionError(res, err); err != nil {
				continue
			}
			respValue := responseNFTType{}
			err = universalJSONDecode(res.Body(), &respValue, nil, func() (bool, bool) {
				return respValue.Result.Value.Reserve > "", false
			})
			if err != nil {
				continue
			}
			base.AllowMint = respValue.Result.Value.AllowMint
			base.TokenURI = respValue.Result.Value.TokenURI
			base.TotalReserve = respValue.Result.Value.Reserve
			for _, own := range respValue.Result.Value.Owners.Owners {
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
	if api.directConn != nil {
		return nil, ErrNotImplemented
	}
	type responseNFTType struct {
		OK     bool `json:"ok"`
		Result *NFT `json:"result"`
	}
	//request
	res, err := api.client.rest.R().Get(fmt.Sprintf("/nfts/%s", id))
	if err = processConnectionError(res, err); err != nil {
		return nil, err
	}
	//json decode
	respValue, respErr := responseNFTType{}, Error{}
	err = universalJSONDecode(res.Body(), &respValue, &respErr, func() (bool, bool) {
		return respValue.OK, respErr.StatusCode != 0
	})
	if err != nil {
		return nil, joinErrors(err, respErr)
	}
	//process result
	return respValue.Result, nil
}

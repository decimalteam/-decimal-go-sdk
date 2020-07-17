package api

import (
	"encoding/json"
	"fmt"
	"time"
)

// MultisigWalletResponse contains API response.
type MultisigWalletResponse struct {
	OK     bool                  `json:"ok"`
	Result *MultisigWalletResult `json:"result"`
}

// MultisigWalletsResponse contains API response.
type MultisigWalletsResponse struct {
	OK     bool `json:"ok"`
	Result struct {
		Count   int                      `json:"count"`
		Wallets []*MultisigWalletsResult `json:"wallets"`
	} `json:"result"`
}

// MultisigTransactionsResponse contains API response.
type MultisigTransactionsResponse struct {
	OK     bool `json:"ok"`
	Result struct {
		Count        int                          `json:"count"`
		Transactions []*MultisigTransactionResult `json:"txs"`
	} `json:"result"`
}

// MultisigWalletResult contains API response fields.
type MultisigWalletResult struct {
	MultisigWallet
	Owners       []*MultisigWalletOwner       `json:"owners"`
	Transactions []*MultisigTransactionResult `json:"txs"`
	Account      *MultisigAccount             `json:"account"`
}

// MultisigWalletsResult contains API response fields.
type MultisigWalletsResult struct {
	MultisigWalletOwner
	Wallet *MultisigWallet `json:"wallet"`
}

// MultisigTransactionResult contains API response fields.
type MultisigTransactionResult struct {
	Transaction   string `json:"transaction"`
	Confirmed     bool   `json:"confirmed"`
	Confirmations uint64 `json:"confirmations"`
	// TODO: field `data`
	Coins []*struct {
		Coin   string `json:"coin"`
		Amount string `json:"amount"`
	} `json:"coin"`
	Address   string    `json:"address"`
	To        string    `json:"to"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// MultisigWallet contains info about multisig wallet.
type MultisigWallet struct {
	Address   string    `json:"address"`
	Threshold uint64    `json:"threshold"`
	Creator   string    `json:"creator"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// MultisigWalletOwner contains info about multisig wallet participant.
type MultisigWalletOwner struct {
	ID        uint64    `json:"id"`
	Address   string    `json:"address"`
	Multisig  string    `json:"multisig"`
	Weight    uint64    `json:"weight"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// MultisigAccount contains info about multisig wallet underlying account.
type MultisigAccount struct {
	ID      uint64 `json:"id"`
	Address string `json:"address"`
	Type    string `json:"type"`
	// TODO: Field balance
	// TODO: Field nonce
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// MultisigWallets requests full list of multisig wallets which has account with specified address as participant.
func (api *API) MultisigWallets(address string) ([]*MultisigWalletsResult, error) {

	url := fmt.Sprintf("/address/%s/multisigs", address)
	res, err := api.client.R().Get(url)
	if err != nil {
		return nil, err
	}
	if res.IsError() {
		return nil, NewResponseError(res)
	}

	response := MultisigWalletsResponse{}
	err = json.Unmarshal(res.Body(), &response)
	if err != nil || !response.OK {
		responseError := Error{}
		err = json.Unmarshal(res.Body(), &responseError)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("received response containing error: %s", responseError.Error())
	}

	return response.Result.Wallets, nil
}

// MultisigWallet requests multisig wallet with specified address.
func (api *API) MultisigWallet(address string) (*MultisigWalletResult, error) {

	url := fmt.Sprintf("/multisig/%s", address)
	res, err := api.client.R().Get(url)
	if err != nil {
		return nil, err
	}
	if res.IsError() {
		return nil, NewResponseError(res)
	}

	response := MultisigWalletResponse{}
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

// MultisigTransactions requests full list of transactions in multisig wallet with specified address.
func (api *API) MultisigTransactions(address string) ([]*MultisigTransactionResult, error) {

	url := fmt.Sprintf("/multisig/%s/txs", address)
	res, err := api.client.R().Get(url)
	if err != nil {
		return nil, err
	}
	if res.IsError() {
		return nil, NewResponseError(res)
	}

	response := MultisigTransactionsResponse{}
	err = json.Unmarshal(res.Body(), &response)
	if err != nil || !response.OK {
		responseError := Error{}
		err = json.Unmarshal(res.Body(), &responseError)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("received response containing error: %s", responseError.Error())
	}

	return response.Result.Transactions, nil
}

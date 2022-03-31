package api

import (
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
	Data          map[string]struct {
		SignerWeight uint64    `json:"signer_weight"`
		Timestamp    time.Time `json:"timestamp"`
	} `json:"data"`
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
	ID        uint64                  `json:"id"`
	Address   string                  `json:"address"`
	Type      string                  `json:"type"`
	Balance   map[string]BalanceEntry `json:"balance"`
	CreatedAt time.Time               `json:"createdAt"`
	UpdatedAt time.Time               `json:"updatedAt"`
}

type BalanceEntry struct {
	Value  string `json:"value"`
	Avatar string `json:"avatar"`
}

// TODO: undefined interfaces with miltisigs in REST/RPC.
// https://bitbucket.org/decimalteam/go-node/src/master/x/multisig/client/rest/query.go

// MultisigWallets requests full list of multisig wallets which has account with specified address as participant.
func (api *API) MultisigWallets(address string) ([]*MultisigWalletsResult, error) {
	if api.directConn != nil {
		return nil, ErrNotImplemented
	}
	//request
	res, err := api.client.rest.R().Get(fmt.Sprintf("/address/%s/multisigs", address))
	if err = processConnectionError(res, err); err != nil {
		return nil, err
	}
	//json decode
	respValue, respErr := MultisigWalletsResponse{}, Error{}
	err = universalJSONDecode(res.Body(), &respValue, &respErr, func() (bool, bool) {
		return respValue.OK, respErr.StatusCode != 0
	})
	if err != nil {
		return nil, joinErrors(err, respErr)
	}
	//process result
	return respValue.Result.Wallets, nil
}

// MultisigWallet requests multisig wallet with specified address.
func (api *API) MultisigWallet(address string) (*MultisigWalletResult, error) {
	if api.directConn != nil {
		return nil, ErrNotImplemented
	}
	//request
	res, err := api.client.rest.R().Get(fmt.Sprintf("/multisig/%s", address))
	if err = processConnectionError(res, err); err != nil {
		return nil, err
	}
	//response
	respValue, respErr := MultisigWalletResponse{}, Error{}
	err = universalJSONDecode(res.Body(), &respValue, &respErr, func() (bool, bool) {
		return respValue.OK, respErr.StatusCode != 0
	})
	if err != nil {
		return nil, joinErrors(err, respErr)
	}
	//process result
	return respValue.Result, nil
}

// MultisigTransactions requests full list of transactions in multisig wallet with specified address.
func (api *API) MultisigTransactions(address string) ([]*MultisigTransactionResult, error) {
	if api.directConn != nil {
		return nil, ErrNotImplemented
	}
	//request
	res, err := api.client.rest.R().Get(fmt.Sprintf("/multisig/%s/txs", address))
	if err = processConnectionError(res, err); err != nil {
		return nil, err
	}
	//json decode
	respValue, respErr := MultisigTransactionsResponse{}, Error{}
	err = universalJSONDecode(res.Body(), &respValue, &respErr, func() (bool, bool) {
		return respValue.OK, respErr.StatusCode != 0
	})
	if err != nil {
		return nil, joinErrors(err, respErr)
	}
	//process result
	return respValue.Result.Transactions, nil
}

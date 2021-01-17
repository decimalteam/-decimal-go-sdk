package api

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

////////////////////////////////////////////////////////////////
// For each transaction sender should pay fee.
// Fees are measured in "units":
//     1 unit  =  10^15 pip  =  0.001 DEL
// E.g. MsgSendCoin fee is 10 units.
// Also sender should pay extra 2 units per byte in raw tx.
////////////////////////////////////////////////////////////////

// Fee is an alias of the type of variables containing amount of units in a value.
type Fee float64

// Unit defines cost of 1 unit in DEL.
const unit Fee = 0.001

// Fees for `coin/*` messages.
const (
	FeeCoinCreate      Fee = 100
	FeeCoinSend        Fee = 10
	FeeCoinMultiSend   Fee = 10 // +5 for each send except first one
	FeeCoinBuy         Fee = 100
	FeeCoinSell        Fee = 100
	FeeCoinSellAll     Fee = 100
	FeeCoinRedeemCheck Fee = 30
)

// Fees for `multisig/*` messages.
const (
	FeeMultisigCreateWallet      Fee = 100
	FeeMultisigCreateTransaction Fee = 100
	FeeMultisigSignTransaction   Fee = 100
)

// Fees for `validator/*` messages.
const (
	FeeValidatorDeclareCandidate Fee = 10000
	FeeValidatorEditCandidate    Fee = 10000
	FeeValidatorDelegate         Fee = 200
	FeeValidatorUnbond           Fee = 200
	FeeValidatorSetOnline        Fee = 100
	FeeValidatorSetOffline       Fee = 100
)

// EstimateTransactionGasWanted counts complete set of different fees and
// returns exact gas wanted to successfully execute specified transaction.
func (api *API) EstimateTransactionGasWanted(tx auth.StdTx) (uint64, error) {
	gasWanted, err := api.getTransactionFee(tx)
	if err != nil {
		return 0, err
	}
	for _, msg := range tx.Msgs {
		fee, err := api.getMessageFee(msg)
		if err != nil {
			return 0, err
		}
		gasWanted += fee
	}
	return uint64(gasWanted), nil
}

// getMessageFee returns amount of fixed units needed to pay for the specified message.
func (api *API) getMessageFee(msg sdk.Msg) (fee Fee, err error) {
	switch r := msg.Route(); r {
	case "coin":
		switch t := msg.Type(); t {
		case "create_coin":
			fee = FeeCoinCreate
		case "send_coin":
			fee = FeeCoinSend
		case "multi_send_coin":
			msgMultiSendCoin, ok := msg.(*MsgMultiSendCoin)
			if !ok {
				err = fmt.Errorf("unable to cast message to type *MsgMultiSendCoin")
				return
			}
			fee = FeeCoinMultiSend + Fee(len(msgMultiSendCoin.Sends)-1)*5
		case "buy_coin":
			fee = FeeCoinBuy
		case "sell_coin":
			fee = FeeCoinSell
		case "sell_all_coin":
			fee = FeeCoinSellAll
		case "redeem_check":
			fee = FeeCoinRedeemCheck
		default:
			err = fmt.Errorf(`unexpected message "coin/%s"`, t)
		}
	case "multisig":
		switch t := msg.Type(); t {
		case "create_wallet":
			fee = FeeMultisigCreateWallet
		case "create_transaction":
			fee = FeeMultisigCreateTransaction
		case "sign_transaction":
			fee = FeeMultisigSignTransaction
		default:
			err = fmt.Errorf(`unexpected message "multisig/%s"`, t)
		}
	case "validator":
		switch t := msg.Type(); t {
		case "declare_candidate":
			fee = FeeValidatorDeclareCandidate
		case "edit_candidate":
			fee = FeeValidatorEditCandidate
		case "delegate":
			fee = FeeValidatorDelegate
		case "unbond":
			fee = FeeValidatorUnbond
		case "set_online":
			fee = FeeValidatorSetOnline
		case "set_offline":
			fee = FeeValidatorSetOffline
		default:
			err = fmt.Errorf(`unexpected message "validator/%s"`, t)
		}
	case "gov":
		fee = 0
	default:
		err = fmt.Errorf(`unexpected message route "%s"`, r)
	}
	return
}

// e18 is a multiply factor to convert value in coin (like DEL) to value in pip.
var e18 = sdk.NewInt(1000000000000000000)

// getMessageSpecialFee returns amount of coins needed to pay for the specified message.
// NOTE: This payment is not counted as gas or units. It is a fee in base or custom coin.
// It is spent only when transaction was successfully executed. In case of failure this
// payment will not be charged.
// TODO: Implement calculating fee in custom coins.
func (api *API) getMessageSpecialFee(msg sdk.Msg) (payment sdk.Coin, err error) {
	msgID := fmt.Sprintf("%s/%s", msg.Route(), msg.Type())

	// Special fee for "coin/create_coin" message
	if msgID == "coin/create_coin" {
		msgCreateCoin, ok := msg.(*MsgCreateCoin)
		if !ok {
			err = fmt.Errorf("unable to cast message to type *MsgCreateCoin")
			return
		}
		var amountInBaseCoin sdk.Int
		switch len(msgCreateCoin.Symbol) {
		case 3:
			amountInBaseCoin = sdk.NewInt(1_000_000)
		case 4:
			amountInBaseCoin = sdk.NewInt(100_000)
		case 5:
			amountInBaseCoin = sdk.NewInt(10_000)
		case 6:
			amountInBaseCoin = sdk.NewInt(1_000)
		default:
			amountInBaseCoin = sdk.NewInt(100)
		}
		payment = sdk.NewCoin(BaseCoinSymbol, amountInBaseCoin.Mul(e18))
	}

	return
}

// getTransactionFee returns amount of units needed to pay for transaction bytes (2 units for each byte).
func (api *API) getTransactionFee(tx auth.StdTx) (fee Fee, err error) {
	txBytes, err := api.codec.MarshalBinaryLengthPrefixed(tx)
	if err != nil {
		return
	}
	fee = Fee(len(txBytes)) * 2
	return
}

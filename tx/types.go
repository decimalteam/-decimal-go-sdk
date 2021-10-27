package tx

const (
	COIN_BUY                    = "coin/buy_coin"
	COIN_CREATE                 = "coin/create_coin"
	COIN_UPDATE                 = "coin/update_coin"
	COIN_SELL                   = "coin/sell_coin"
	COIN_SEND                   = "coin/send_coin"
	COIN_MULTISEND              = "coin/multi_send_coin"
	COIN_SELL_ALL               = "coin/sell_all_coin"
	COIN_REDEEM_CHECK           = "coin/redeem_check"
	COIN_ISSUE_CHECK            = "coin/issue_check"
	VALIDATOR_CANDIDATE_DECLARE = "validator/declare_candidate"
	VALIDATOR_DELEGATE          = "validator/delegate"
	VALIDATOR_SET_ONLINE        = "validator/set_online"
	VALIDATOR_SET_OFFLINE       = "validator/set_offline"
	VALIDATOR_UNBOND            = "validator/unbond"
	VALIDATOR_CANDIDATE_EDIT    = "validator/edit_candidate"
	MULTISIG_CREATE_WALLET      = "multisig/create_wallet"
	MULTISIG_CREATE_TX          = "multisig/create_transaction"
	MULTISIG_SIGN_TX            = "multisig/sign_transaction"

	PROPOSAL_SUBMIT           = "cosmos-sdk/MsgSubmitProposal"
	PROPOSAL_SOFTWARE_UPGRADE = "cosmos-sdk/MsgSoftwareUpgradeProposal"
	PROPOSAL_VOTE             = "cosmos-sdk/MsgVote"

	SWAP_HTLT   = "swap/msg_htlt"
	SWAP_REDEEM = "swap/msg_redeem"
	SWAP_REFUND = "swap/msg_refund"

	NFT_MINT           = "nft/msg_mint"
	NFT_BURN           = "nft/msg_burn"
	NFT_EDIT_METADATA  = "nft/msg_edit_metadata"
	NFT_TRANSFER       = "nft/msg_transfer"
	NFT_DELEGATE       = "nft/msg_delegate"
	NFT_UNBOND         = "nft/msg_unbond"
	NFT_UPDATE_RESERVE = "nft/update_reserve"
)

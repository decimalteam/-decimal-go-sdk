package wallet

import (
	"crypto/sha256"
	"errors"
	"math/big"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/base58"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/bech32"

	sdk "github.com/cosmos/cosmos-sdk/types"
	auth "github.com/cosmos/cosmos-sdk/x/auth"
)

const (
	mnemonicBits   = 256
	derivationPath = "m/44'/60'/0'/0/0"
	addressPrefix  = "dx"
)

// Account contains private key of the account that allows to sign transactions to broadcast to the blockchain.
type Account struct {
	privateKey   *PrivateKey
	privateKeyTM *secp256k1.PrivKeySecp256k1
	publicKeyTM  *secp256k1.PubKeySecp256k1
	address      string

	// These fields are used only for signing transactions:
	chainID       string
	accountNumber int64
	sequence      int64
}

// NewAccount creates new account with random mnemonic.
func NewAccount(password string) (*Account, error) {
	mnemonic, err := NewMnemonic(mnemonicBits, password)
	if err != nil {
		return nil, err
	}
	return NewAccountFromMnemonic(mnemonic)
}

// NewAccountFromMnemonicWords creates account from mnemonic presented as set of words.
func NewAccountFromMnemonicWords(words string, password string) (*Account, error) {
	mnemonic, err := NewMnemonicFromWords(words, password)
	if err != nil {
		return nil, err
	}
	return NewAccountFromMnemonic(mnemonic)
}

// NewAccountFromMnemonic creates account from mnemonic.
func NewAccountFromMnemonic(mnemonic *Mnemonic) (*Account, error) {

	// Create HD derivation seed and root extended key from specified mnemonic
	extendedKey, err := hdkeychain.NewMaster(mnemonic.Seed(), &chaincfg.MainNetParams)
	if err != nil {
		return nil, err
	}

	// Derive root extended key to `derivationPath`
	extendedKey, err = extendedKeyDerivePath(extendedKey, derivationPath)
	if err != nil {
		return nil, err
	}

	// Calculate private key from derrived extended key
	privateKey, err := extendedKeyToPrivateKey(extendedKey)
	if err != nil {
		return nil, err
	}

	// Prepare private key as secp256k1.PrivKeySecp256k1 object which can be used to sign transactions
	privateKeyTM := &secp256k1.PrivKeySecp256k1{}
	copy(privateKeyTM[:], privateKey.Bytes())

	// Prepare public key as secp256k1.PubKeySecp256k1 object
	publicKeyTM := &secp256k1.PubKeySecp256k1{}
	copy(publicKeyTM[:], privateKey.PublicKey().BytesCompressed())

	// Calculate address from private key
	address := privateKey.PublicKey().AddressString()

	// Create and return account
	result := &Account{
		privateKey:    privateKey,
		privateKeyTM:  privateKeyTM,
		publicKeyTM:   publicKeyTM,
		address:       address,
		accountNumber: -1,
		sequence:      -1,
	}

	return result, nil
}

// WithChainID sets chain ID of network.
func (acc *Account) WithChainID(chainID string) *Account {
	acc.chainID = chainID
	return acc
}

// WithAccountNumber sets accounts's number.
func (acc *Account) WithAccountNumber(accountNumber uint64) *Account {
	acc.accountNumber = int64(accountNumber)
	return acc
}

// WithSequence sets accounts's sequence (last used nonce).
func (acc *Account) WithSequence(sequence uint64) *Account {
	acc.sequence = int64(sequence)
	return acc
}

// PrivateKey returns accounts's private key.
func (acc *Account) PrivateKey() *PrivateKey {
	return acc.privateKey
}

// Address returns accounts's address in bech32 format.
func (acc *Account) Address() string {
	return acc.address
}

// ChainID returns chain ID of network.
func (acc *Account) ChainID() string {
	return acc.chainID
}

// AccountNumber returns accounts's number.
func (acc *Account) AccountNumber() int64 {
	return acc.accountNumber
}

// Sequence returns accounts's sequence (last used nonce).
func (acc *Account) Sequence() int64 {
	return acc.sequence
}

// CreateTransaction creates new transaction with specified messages and parameters.
func (acc *Account) CreateTransaction(msgs []sdk.Msg, fee auth.StdFee, memo string) auth.StdTx {
	return auth.NewStdTx(msgs, fee, nil, memo)
}

// SignTransaction signs transaction and appends signature to transaction signatures.
func (acc *Account) SignTransaction(tx auth.StdTx) (auth.StdTx, error) {

	// Check chain ID, account number and sequence
	if len(acc.chainID) == 0 {
		return tx, errors.New("chain ID is not set up")
	}
	if acc.accountNumber < 0 || acc.sequence < 0 {
		return tx, errors.New("account number or sequence is not set up")
	}

	// Retrieve transaction bytes required to sign
	bytesToSign := auth.StdSignBytes(
		acc.chainID, uint64(acc.accountNumber), uint64(acc.sequence),
		tx.Fee, tx.Msgs, tx.Memo,
	)

	// Sign bytes prepared to sign
	signatureBytes, err := acc.privateKeyTM.Sign(bytesToSign)
	if err != nil {
		return tx, err
	}

	// Prepare auth.StdSignature object
	signature := auth.StdSignature{
		PubKey:    acc.publicKeyTM,
		Signature: signatureBytes,
	}

	// Copy input transaction and append signature to the list
	tx.Signatures = append(tx.Signatures, signature)

	return tx, err
}

// IssueCheck issues check and returns it as base58 string.
func (acc *Account) IssueCheck(coinSymbol string, amount sdk.Int, nonce sdk.Int, dueBlock uint64, passphrase string) (string, error) {

	// TODO: Check if coin exists?
	// TODO: Check amount?

	// Prepare private key from passphrase
	passphraseHash := sha256.Sum256([]byte(passphrase))
	passphrasePrivKey, _ := crypto.ToECDSA(passphraseHash[:])

	// Prepare check without lock
	check := &Check{
		ChainID:  acc.chainID,
		Coin:     coinSymbol,
		Amount:   amount.BigInt(),
		Nonce:    nonce.BigInt().Bytes(),
		DueBlock: dueBlock,
	}

	// Prepare check lock
	checkHash := check.HashWithoutLock()
	lock, _ := crypto.Sign(checkHash[:], passphrasePrivKey)

	// Fill check with prepared lock
	check.Lock = big.NewInt(0).SetBytes(lock)

	// Retrieve private key from the keybase account
	privKeyECDSA, err := crypto.ToECDSA(acc.privateKeyTM[:])
	if err != nil {
		panic(err)
	}
	// address := sdk.AccAddress(privKey.PubKey().Address())

	// Sign check by check issuer
	checkHash = check.Hash()
	signature, err := crypto.Sign(checkHash[:], privKeyECDSA)
	if err != nil {
		panic(err)
	}
	check.SetSignature(signature)

	// Return issued raw check encoded to base64 format to the issuer
	checkBytes, err := rlp.EncodeToBytes(check)
	if err != nil {
		panic(err)
	}

	return base58.Encode(checkBytes), nil
}

func (acc *Account) RedeemCheck() (string, error) {
	// todo
	_ = bech32.Decode(acc.address)

	_, _ = crypto.Sign([]byte{}, "")

	return "", nil
}

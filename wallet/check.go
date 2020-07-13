package wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// HashLength represents fixed hash length.
const HashLength = 32

// Hash represents the 32 byte Keccak256 hash of arbitrary data.
type Hash [HashLength]byte

// Errors.
var (
	ErrInvalidSig = errors.New("invalid transaction v, r, s values")
)

// Check is a struct providing helper methods for check issuing.
type Check struct {
	ChainID  string
	Coin     string
	Amount   *big.Int
	Nonce    []byte
	DueBlock uint64
	Lock     *big.Int
	V        *big.Int
	R        *big.Int
	S        *big.Int
}

// Sender returns address of check issuer.
func (check *Check) Sender() (sdk.AccAddress, error) {
	return recoverPlain(check.Hash(), check.R, check.S, check.V)
}

// LockPubKey recovers public key from check's lock and returns it.
func (check *Check) LockPubKey() ([]byte, error) {
	sig := check.Lock.Bytes()

	if len(sig) < 65 {
		sig = append(make([]byte, 65-len(sig)), sig...)
	}

	hash := check.HashWithoutLock()

	pub, err := crypto.Ecrecover(hash[:], sig)
	if err != nil {
		return nil, err
	}
	if len(pub) == 0 || pub[0] != 4 {
		return nil, errors.New("invalid public key")
	}

	return pub, nil
}

// HashWithoutLock returns hash from check fields (witout check's lock).
func (check *Check) HashWithoutLock() Hash {
	return rlpHash([]interface{}{
		check.ChainID,
		check.Coin,
		check.Amount,
		check.Nonce,
		check.DueBlock,
	})
}

// Hash returns hash from check fields.
func (check *Check) Hash() Hash {
	return rlpHash([]interface{}{
		check.ChainID,
		check.Coin,
		check.Amount,
		check.Nonce,
		check.DueBlock,
		check.Lock,
	})
}

// HashFull returns hash full from check fields (including check's signature fields).
func (check *Check) HashFull() Hash {
	return rlpHash([]interface{}{
		check.ChainID,
		check.Coin,
		check.Amount,
		check.Nonce,
		check.DueBlock,
		check.Lock,
		check.V,
		check.R,
		check.S,
	})
}

// Sign signs the check and sets check signature fields.
func (check *Check) Sign(prv *ecdsa.PrivateKey) error {
	h := check.Hash()
	sig, err := crypto.Sign(h[:], prv)
	if err != nil {
		return err
	}

	check.SetSignature(sig)

	return nil
}

// SetSignature sets check signature fields.
func (check *Check) SetSignature(sig []byte) {
	check.R = new(big.Int).SetBytes(sig[:32])
	check.S = new(big.Int).SetBytes(sig[32:64])
	check.V = new(big.Int).SetBytes([]byte{sig[64] + 27})
}

// String returns string representation of the check.
func (check *Check) String() string {
	sender, _ := check.Sender()

	return fmt.Sprintf("Check sender: %s nonce: %x, dueBlock: %d, value: %s %s", sender.String(), check.Nonce,
		check.DueBlock, check.Amount.String(), check.Coin)
	// return fmt.Sprintf("Check nonce: %x, dueBlock: %d, value: %s %s", check.Nonce, check.DueBlock, check.Amount.String(), check.Coin)
}

// ParseCheck parses check from bytes.
func ParseCheck(buf []byte) (*Check, error) {
	var check Check
	err := rlp.Decode(bytes.NewReader(buf), &check)
	if err != nil {
		return nil, err
	}

	if check.S == nil || check.R == nil || check.V == nil {
		return nil, errors.New("incorrect tx signature")
	}

	return &check, nil
}

func rlpHash(x interface{}) (h Hash) {
	hw := sha3.NewLegacyKeccak256()
	err := rlp.Encode(hw, x)
	if err != nil {
		panic(err)
	}
	hw.Sum(h[:0])
	return h
}

func recoverPlain(sighash Hash, R, S, Vb *big.Int) (sdk.AccAddress, error) {
	if Vb.BitLen() > 8 {
		return sdk.AccAddress{}, ErrInvalidSig
	}
	V := byte(Vb.Uint64() - 27)
	if !crypto.ValidateSignatureValues(V, R, S, true) {
		return sdk.AccAddress{}, ErrInvalidSig
	}
	// encode the snature in uncompressed format
	r, s := R.Bytes(), S.Bytes()
	sig := make([]byte, 65)
	copy(sig[32-len(r):32], r)
	copy(sig[64-len(s):64], s)
	sig[64] = V
	// recover the public key from the snature
	pub, err := crypto.Ecrecover(sighash[:], sig)
	if err != nil {
		return sdk.AccAddress{}, err
	}
	if len(pub) == 0 || pub[0] != 4 {
		return sdk.AccAddress{}, errors.New("invalid public key")
	}
	pub2, err := crypto.UnmarshalPubkey(pub)
	if err != nil {
		return sdk.AccAddress{}, err
	}
	pub3 := crypto.CompressPubkey(pub2)
	hasherSHA256 := sha256.New()
	hasherSHA256.Write(pub3)
	sha := hasherSHA256.Sum(nil)
	hasherRIPEMD160 := ripemd160.New()
	hasherRIPEMD160.Write(sha)
	return sdk.AccAddress(hasherRIPEMD160.Sum(nil)), nil
}

package wallet

import (
	"crypto/ecdsa"
	"encoding/base64"
	"math/big"

	"github.com/btcsuite/btcd/btcec"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tendermint/tendermint/libs/bech32"
)

// PublicKey contains public key object.
type PublicKey struct {
	data           []byte
	ecdsaPublicKey *ecdsa.PublicKey
}

// NewPublicKeyFromPrivateKeyBytes creates instance of public key object from private key as byte array.
func NewPublicKeyFromPrivateKeyBytes(privateKey []byte) (*PublicKey, error) {
	x, y := btcec.S256().ScalarBaseMult(privateKey)
	ecdsaPublicKey := &ecdsa.PublicKey{
		Curve: btcec.S256(),
		X:     x,
		Y:     y,
	}
	xb := ecdsaPublicKey.X.Bytes()
	if len(xb) < 32 {
		xb = append(make([]byte, (32-len(xb))), xb...)
	}
	yb := ecdsaPublicKey.Y.Bytes()
	if len(yb) < 32 {
		yb = append(make([]byte, (32-len(yb))), yb...)
	}
	data := make([]byte, 65)
	data[0] = 4
	copy(data[1:33], xb)
	copy(data[33:65], yb)
	return &PublicKey{
		data:           data,
		ecdsaPublicKey: ecdsaPublicKey,
	}, nil
}

// NewPublicKeyFromBytes creates instance of public key object from public key as byte array.
func NewPublicKeyFromBytes(data []byte) (*PublicKey, error) {
	ecdsaPublicKey := &ecdsa.PublicKey{
		Curve: btcec.S256(),
		X:     new(big.Int).SetBytes(data[1:33]),
		Y:     new(big.Int).SetBytes(data[33:65]),
	}
	return &PublicKey{
		data:           data,
		ecdsaPublicKey: ecdsaPublicKey,
	}, nil
}

// String returns string representation of the public key.
// It is actually just public key (in compressed form) presented
// as byte array with 33 bytes length encoded to base64 format.
func (key *PublicKey) String() string {
	return base64.StdEncoding.EncodeToString(key.BytesCompressed())
}

// Bytes returns public key as byte array with 65 bytes length.
func (key *PublicKey) Bytes() []byte {
	return key.data
}

// BytesCompressed returns public key (in compressed form) as byte array with 33 bytes length.
func (key *PublicKey) BytesCompressed() []byte {
	data := make([]byte, 33)
	data[0] = key.data[64]%2 + 2
	copy(data[1:33], key.data[1:33])
	return data
}

// ECDSA returns pointer to base ecdsa.PublicKey.
func (key *PublicKey) ECDSA() *ecdsa.PublicKey {
	return key.ecdsaPublicKey
}

// AddressString returns string representation of the address (calculated from public key).
// It is actually just address presented as byte array with
// 20 bytes length encoded to bech32 format with prefix `addressPrefix`.
func (key *PublicKey) AddressString() string {
	address, err := bech32.ConvertAndEncode(addressPrefix, key.AddressBytes())
	if err != nil {
		return ""
	}
	return address
}

// AddressBytes returns address (calculated from public key) as byte array with 20 bytes length.
func (key *PublicKey) AddressBytes() []byte {
	hash := crypto.Keccak256Hash(key.data[1:]).Bytes()
	return hash[12:]
}

package wallet

import (
	"crypto/ecdsa"
	"encoding/base64"
	"math/big"
)

// PrivateKey contains private key object.
type PrivateKey struct {
	data            []byte
	ecdsaPrivateKey *ecdsa.PrivateKey
	publicKey       *PublicKey
}

// NewPrivateKeyFromBytes creates instance of private key object from byte array.
func NewPrivateKeyFromBytes(data []byte) (*PrivateKey, error) {
	publicKey, err := NewPublicKeyFromPrivateKeyBytes(data)
	if err != nil {
		return nil, err
	}
	ecdsaPrivateKey := &ecdsa.PrivateKey{
		PublicKey: *publicKey.ECDSA(),
		D:         new(big.Int).SetBytes(data),
	}
	return &PrivateKey{
		data:            data,
		ecdsaPrivateKey: ecdsaPrivateKey,
		publicKey:       publicKey,
	}, nil
}

// String returns string representation of the private key.
// It is actually just private key presented as byte array
// with 32 bytes length encoded to base64 format.
func (key *PrivateKey) String() string {
	return base64.StdEncoding.EncodeToString(key.data)
}

// Bytes returns private key as byte array with 32 bytes length.
func (key *PrivateKey) Bytes() []byte {
	return key.data
}

// ECDSA returns pointer to base ecdsa.PrivateKey.
func (key *PrivateKey) ECDSA() *ecdsa.PrivateKey {
	return key.ecdsaPrivateKey
}

// PublicKey returns pointer to public key.
func (key *PrivateKey) PublicKey() *PublicKey {
	return key.publicKey
}

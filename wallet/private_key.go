package wallet

import (
	"crypto/ecdsa"
	"math/big"
)

// PrivateKey contains private key object.
type PrivateKey struct {
	data            []byte
	ecdsaPrivateKey *ecdsa.PrivateKey
	publicKey       *PublicKey
}

// NewPrivateKeyFromBytes creates instance of private key object from byte array.
func NewPrivateKeyFromBytes(privateKey []byte) (*PrivateKey, error) {
	publicKey, err := NewPublicKeyFromPrivateKeyBytes(privateKey)
	if err != nil {
		return nil, err
	}
	ecdsaPrivateKey := &ecdsa.PrivateKey{
		PublicKey: *publicKey.GetECDSA(),
		D:         new(big.Int).SetBytes(privateKey),
	}
	return &PrivateKey{
		data:            privateKey,
		ecdsaPrivateKey: ecdsaPrivateKey,
		publicKey:       publicKey,
	}, nil
}

// GetBytes returns private key as byte array with 32 bytes length.
func (key *PrivateKey) GetBytes() []byte {
	return key.data
}

// GetECDSA returns pointer to base ecdsa.PrivateKey.
func (key *PrivateKey) GetECDSA() *ecdsa.PrivateKey {
	return key.ecdsaPrivateKey
}

// GetPublicKey returns pointer to public key.
func (key *PrivateKey) GetPublicKey() *PublicKey {
	return key.publicKey
}

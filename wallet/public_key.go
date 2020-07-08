package wallet

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/btcsuite/btcd/btcec"
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
	publicKey := make([]byte, 65)
	publicKey[0] = 4
	copy(publicKey[1:33], xb)
	copy(publicKey[33:65], yb)
	return &PublicKey{
		data:           publicKey,
		ecdsaPublicKey: ecdsaPublicKey,
	}, nil
}

// NewPublicKeyFromBytes creates instance of public key object from public key as byte array.
func NewPublicKeyFromBytes(publicKey []byte) (*PublicKey, error) {
	ecdsaPublicKey := &ecdsa.PublicKey{
		Curve: btcec.S256(),
		X:     new(big.Int).SetBytes(publicKey[1:33]),
		Y:     new(big.Int).SetBytes(publicKey[33:65]),
	}
	return &PublicKey{
		data:           publicKey,
		ecdsaPublicKey: ecdsaPublicKey,
	}, nil
}

// GetBytes returns public key as byte array with 65 bytes length.
func (key *PublicKey) GetBytes() []byte {
	return key.data
}

// GetECDSA returns pointer to base ecdsa.PublicKey.
func (key *PublicKey) GetECDSA() *ecdsa.PublicKey {
	return key.ecdsaPublicKey
}

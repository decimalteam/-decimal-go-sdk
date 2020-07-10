package wallet

import (
	"crypto/ecdsa"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/btcsuite/btcutil/hdkeychain"
)

// extendedKeyDerive returns a derived child extended key at the given index.
func extendedKeyDerive(key *hdkeychain.ExtendedKey, index uint32, hardened bool) (*hdkeychain.ExtendedKey, error) {
	if hardened {
		index = index + hdkeychain.HardenedKeyStart
	}
	keyDerrived, err := key.Child(index)
	if err != nil {
		return nil, err
	}
	return keyDerrived, nil
}

// extendedKeyDerivePath returns a derived child extended key at the given path like `m/44'/60'/0'/0/0`.
func extendedKeyDerivePath(key *hdkeychain.ExtendedKey, path string) (*hdkeychain.ExtendedKey, error) {
	pathMatch, err := regexp.MatchString("^(m\\/)?(\\d+'?\\/)*\\d+'?$", path)
	if err != nil {
		return nil, err
	} else if !pathMatch {
		return nil, fmt.Errorf("Specified path %q is invalid", path)
	}
	splitPath := strings.Split(path, "/")
	offset := 0
	if splitPath[0] == "m" {
		if key.ParentFingerprint() != 0 {
			return nil, fmt.Errorf("Required master extended key to get child extended key at path %q", path)
		}
		offset = 1
	}
	keyDerrived := key
	for _, partPath := range splitPath[offset:] {
		partLength := len(partPath)
		if partLength > 0 {
			partHardened := partPath[partLength-1] == '\''
			partEndOffset := 0
			if partHardened {
				partEndOffset = 1
			}
			index, err := strconv.ParseUint(partPath[:partLength-partEndOffset], 10, 32)
			if err != nil {
				return nil, err
			}
			keyDerrived, err = extendedKeyDerive(keyDerrived, uint32(index), partHardened)
			if err != nil {
				return nil, err
			}
		}
	}
	return keyDerrived, nil
}

// extendedKeyToPrivateKey returns a new private key from the extended key.
func extendedKeyToPrivateKey(key *hdkeychain.ExtendedKey) (*PrivateKey, error) {
	ecPrivateKey, err := key.ECPrivKey()
	if err != nil {
		return nil, err
	}
	ecdsaPrivateKey := (*ecdsa.PrivateKey)(ecPrivateKey)
	privateKeyBytes := ecdsaPrivateKey.D.Bytes()
	if len(privateKeyBytes) < 32 {
		privateKeyBytes = append(make([]byte, (32-len(privateKeyBytes))), privateKeyBytes...)
	}
	privateKey, err := NewPrivateKeyFromBytes(privateKeyBytes)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

// extendedKeyToPublicKey returns a new public key from the extended key.
func extendedKeyToPublicKey(key *hdkeychain.ExtendedKey) (*PublicKey, error) {
	ecPublicKey, err := key.ECPubKey()
	if err != nil {
		return nil, err
	}
	ecdsaPublicKey := (*ecdsa.PublicKey)(ecPublicKey)
	publicKeyCompressed := make([]byte, 33)
	publicKeyCompressed[0] = byte(ecdsaPublicKey.Y.Bit(0) + 2)
	x := ecdsaPublicKey.X.Bytes()
	if len(x) < 32 {
		x = append(make([]byte, (32-len(x))), x...)
	}
	copy(publicKeyCompressed[1:33], x)
	publicKey, err := NewPublicKeyFromBytes(publicKeyCompressed)
	if err != nil {
		return nil, err
	}
	return publicKey, nil
}

package wallet

import (
	"crypto/ecdsa"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
)

// ExtendedKey contains hierarchical deterministic extended keys (BIP-0032).
type ExtendedKey struct {
	extendedKey *hdkeychain.ExtendedKey
	masterKey   bool
}

// NewExtendedKeyFromMnemonic creates extended key from mnemonic.
// Created extended key is master extended tree (root node in hierarchical deterministic).
func NewExtendedKeyFromMnemonic(mnemonic *Mnemonic) (*ExtendedKey, error) {
	seed := mnemonic.Seed()
	extendedKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return nil, err
	}
	result := &ExtendedKey{
		extendedKey: extendedKey,
		masterKey:   true,
	}
	return result, nil
}

// NewExtendedKeyFromString creates extended key from private or public key in string format.
func NewExtendedKeyFromString(key string) (*ExtendedKey, error) {
	extendedKey, err := hdkeychain.NewKeyFromString(key)
	if err != nil {
		return nil, err
	}
	result := &ExtendedKey{
		extendedKey: extendedKey,
		masterKey:   extendedKey.ParentFingerprint() == 0,
	}
	return result, nil
}

// IsMaster returns whether or not the extended key is a root extended key in HD keys tree.
func (key *ExtendedKey) IsMaster() bool {
	return key.masterKey
}

// IsPrivate returns whether or not the extended key is a private extended key.
func (key *ExtendedKey) IsPrivate() bool {
	return key.extendedKey.IsPrivate()
}

// GetDepth returns the current derivation level with respect to the root.
func (key *ExtendedKey) GetDepth() uint8 {
	return key.extendedKey.Depth()
}

// GetChild returns a derived child extended key at the given index.
func (key *ExtendedKey) GetChild(i uint32, hardened bool) (*ExtendedKey, error) {
	var index uint32
	if hardened {
		index = hdkeychain.HardenedKeyStart + i
	} else {
		index = i
	}
	extendedKey, err := key.extendedKey.Child(index)
	if err != nil {
		return nil, err
	}
	result := &ExtendedKey{
		extendedKey: extendedKey,
	}
	return result, nil
}

// GetChildAtPath returns a derived child extended key at the given path like `m/44'/0'/0'/0/0` or `0'/0/0`.
func (key *ExtendedKey) GetChildAtPath(path string) (*ExtendedKey, error) {
	pathMatch, err := regexp.MatchString("^(m\\/)?(\\d+'?\\/)*\\d+'?$", path)
	if err != nil {
		return nil, err
	} else if !pathMatch {
		return nil, fmt.Errorf("Specified path %q is invalid", path)
	}
	splitPath := strings.Split(path, "/")
	offset := 0
	if splitPath[0] == "m" {
		if !key.IsMaster() {
			return nil, fmt.Errorf("Required master extended key to get child extended key at path %q", path)
		}
		offset = 1
	}
	extendedKey := key
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
			extendedKey, err = extendedKey.GetChild(uint32(index), partHardened)
			if err != nil {
				return nil, err
			}
		}
	}
	return extendedKey, nil
}

// GetAccountRoot returns a derived child extended key at the specific path `m/44'/60'/account'/0`.
func (key *ExtendedKey) GetAccountRoot(account uint32) (*ExtendedKey, error) {
	path := fmt.Sprintf("m/44'/60'/%d'/0", account)
	return key.GetChildAtPath(path)
}

// GetPublicKey returns a new public extended key from this private extended key.
func (key *ExtendedKey) GetPublicKey() (*ExtendedKey, error) {
	extendedKey, err := key.extendedKey.Neuter()
	if err != nil {
		return nil, err
	}
	result := &ExtendedKey{
		extendedKey: extendedKey,
	}
	return result, nil
}

// GetECPrivateKey returns a new private key from the private extended key.
func (key *ExtendedKey) GetECPrivateKey() (*PrivateKey, error) {
	if !key.IsPrivate() {
		return nil, fmt.Errorf("Unable to get private key from public extended key")
	}
	ecPrivateKey, err := key.extendedKey.ECPrivKey()
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

// GetECPublicKey returns a new public key from the private/public extended key.
func (key *ExtendedKey) GetECPublicKey() (*PublicKey, error) {
	ecPublicKey, err := key.extendedKey.ECPubKey()
	if err != nil {
		return nil, err
	}
	ecdsaPublicKey := (*ecdsa.PublicKey)(ecPublicKey)
	publicKeyUncompressed := make([]byte, 33)
	publicKeyUncompressed[0] = byte(ecdsaPublicKey.Y.Bit(0) + 2)
	x := ecdsaPublicKey.X.Bytes()
	if len(x) < 32 {
		x = append(make([]byte, (32-len(x))), x...)
	}
	copy(publicKeyUncompressed[1:33], x)
	publicKey, err := NewPublicKeyFromBytes(publicKeyUncompressed)
	if err != nil {
		return nil, err
	}
	return publicKey, nil
}

// GetString returns the extended key as a human-readable base58-encoded string.
func (key *ExtendedKey) GetString() string {
	return key.extendedKey.String()
}

func (key *ExtendedKey) String() string {
	return key.GetString()
}

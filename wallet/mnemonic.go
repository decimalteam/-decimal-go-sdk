package wallet

import (
	"github.com/tyler-smith/go-bip39"
)

// Mnemonic contains entropy, seed and mnemonic words array which can be used for hierarchical deterministic extended keys.
type Mnemonic struct {
	entropy []byte
	words   string
	seed    []byte
}

// NewMnemonicRandom creates a new random (crypto safe) Mnemonic. Use 128 bits for a 12 words code or 256 bits for a 24 words.
func NewMnemonicRandom(bits int, password string) (*Mnemonic, error) {
	entropy, err := bip39.NewEntropy(bits)
	if err != nil {
		return nil, err
	}
	words, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return nil, err
	}
	return &Mnemonic{
		entropy: entropy,
		words:   words,
		seed:    bip39.NewSeed(words, password),
	}, nil
}

// NewMnemonicFromEntropy creates a Mnemonic based on a known entropy.
func NewMnemonicFromEntropy(entropy []byte, password string) (*Mnemonic, error) {
	words, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return nil, err
	}
	return &Mnemonic{
		entropy: entropy,
		words:   words,
		seed:    bip39.NewSeed(words, password),
	}, nil
}

// NewMnemonicFromWords creates a Mnemonic based on a known list of words.
func NewMnemonicFromWords(words string, password string) (*Mnemonic, error) {
	entropy, err := bip39.MnemonicToByteArray(words, true)
	if err != nil {
		return nil, err
	}
	return &Mnemonic{
		entropy: entropy,
		words:   words,
		seed:    bip39.NewSeed(words, password),
	}, nil
}

// Entropy returns the entropy of the Mnemonic.
func (m *Mnemonic) Entropy() []byte {
	return m.entropy
}

// Words returns the words from the Mnemonic.
func (m *Mnemonic) Words() string {
	return m.words
}

// Seed returns the seed of the Mnemonic.
func (m *Mnemonic) Seed() []byte {
	return m.seed
}

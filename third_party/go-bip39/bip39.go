// Package bip39 re-exports the luxfi fork (github.com/luxfi/go-bip39) under the
// original github.com/tyler-smith/go-bip39 module path. See go.mod for why this
// shim exists. The full public surface of the upstream package is forwarded so
// the shim is a faithful drop-in, not a partial one tied to today's callers.
package bip39

import bip39 "github.com/luxfi/go-bip39"

// Mirror the upstream sentinel errors so callers can compare against them.
var (
	ErrInvalidMnemonic             = bip39.ErrInvalidMnemonic
	ErrEntropyLengthInvalid        = bip39.ErrEntropyLengthInvalid
	ErrValidatedSeedLengthMismatch = bip39.ErrValidatedSeedLengthMismatch
	ErrChecksumIncorrect           = bip39.ErrChecksumIncorrect
)

// SetWordList sets the word list used package-wide for mnemonics.
func SetWordList(list []string) { bip39.SetWordList(list) }

// GetWordList returns the current word list.
func GetWordList() []string { return bip39.GetWordList() }

// GetWordIndex returns the index of a word in the current word list.
func GetWordIndex(word string) (int, bool) { return bip39.GetWordIndex(word) }

// NewEntropy returns randomly generated entropy of the given bit size.
func NewEntropy(bitSize int) ([]byte, error) { return bip39.NewEntropy(bitSize) }

// EntropyFromMnemonic recovers the entropy that produced the given mnemonic.
func EntropyFromMnemonic(mnemonic string) ([]byte, error) {
	return bip39.EntropyFromMnemonic(mnemonic)
}

// NewMnemonic returns a mnemonic for the given entropy.
func NewMnemonic(entropy []byte) (string, error) { return bip39.NewMnemonic(entropy) }

// MnemonicToByteArray converts a mnemonic into its byte-array representation.
func MnemonicToByteArray(mnemonic string, raw ...bool) ([]byte, error) {
	return bip39.MnemonicToByteArray(mnemonic, raw...)
}

// NewSeedWithErrorChecking creates a hashed seed, validating the mnemonic first.
func NewSeedWithErrorChecking(mnemonic, password string) ([]byte, error) {
	return bip39.NewSeedWithErrorChecking(mnemonic, password)
}

// NewSeed creates a hashed seed from a mnemonic and passphrase.
func NewSeed(mnemonic, password string) []byte { return bip39.NewSeed(mnemonic, password) }

// IsMnemonicValid reports whether a mnemonic is valid for the current word list.
func IsMnemonicValid(mnemonic string) bool { return bip39.IsMnemonicValid(mnemonic) }

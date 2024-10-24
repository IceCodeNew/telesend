package aead

import (
	"strconv"

	"github.com/IceCodeNew/telesend/pkg/crypto"
	"github.com/cloudflare/circl/cipher/ascon"
)

var additionalData = []byte("IceCode/TELESEND")

func EncAscon128a(key, nonce, plaintext []byte) ([]byte, error) {
	block, err := ascon.New(key, ascon.Ascon128a)
	if err != nil {
		return nil, err
	}
	return block.Seal(plaintext[:0], nonce, plaintext, additionalData), nil
}

func DecAscon128a(key, nonce, ciphertext []byte) ([]byte, error) {
	block, err := ascon.New(key, ascon.Ascon128a)
	if err != nil {
		return nil, err
	}
	return block.Open(ciphertext[:0], nonce, ciphertext, additionalData)
}

func DeriveKeyAndNonce(passphrase, salt []byte) (key, nonce []byte) {
	kn := crypto.DeriveKey(passphrase, salt, ascon.KeySize+ascon.NonceSize)
	key, nonce = kn[:ascon.KeySize], kn[ascon.KeySize:]
	return
}

func PredictableSeed(seed1, seed2 uint64) [32]byte {
	if seed1 > seed2 {
		seed1, seed2 = seed2, seed1
	}

	// It does not matter if the product overflowed
	//
	// The final result of the product would length at least 16 digits, at most 20 digits.
	// (18446744073709551615)
	product := seed1 * seed2
	for product < uint64(1000000000000000) {
		product *= seed2
	}
	_strProduct := strconv.FormatUint(product, 10)

	seed := append(make([]byte, 0, 32), _strProduct...)

	padLen := (32 - len(_strProduct))
	// let compiler aware of the fact that
	// padLen <= len(additionalData)
	padLen %= len(additionalData) + 1

	seed = append(seed, additionalData[:padLen]...)
	return [32]byte(seed)
}

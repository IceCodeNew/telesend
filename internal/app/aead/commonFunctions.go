package aead

import (
	"strconv"

	"golang.org/x/crypto/argon2"
)

func DeriveKey(passphrase, salt []byte, keyLen uint32) (key []byte) {
	return argon2.IDKey(passphrase, salt, 1, 64*1024, 4, keyLen)
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

var additionalData = []byte("IceCode/TELESEND")

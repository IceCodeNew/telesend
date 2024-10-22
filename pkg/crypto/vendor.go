package crypto

import (
	"crypto/rand"
)

// Imported from https://stackoverflow.com/a/44477359/13631331, with modifications

// len(encodeStr) = 8+5+11+11+2+11+10+6 = 64.
//
// This allows (0 <= x <= 255) x % 64 to have an even distribution.
const encodeStr = "ABCDEFGH" + "JKLMN" + "PQRSTUVWXYZ" +
	"abcdefghijk" + "mn" + "pqrstuvwxyz" + "0123456789" +
	"-_+=,."

// A helper function create and fill a slice of length n with characters from the following string:
//
// "ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnpqrstuvwxyz0123456789-_+=,."
func RandAsciiBytes(n int) ([]byte, error) {
	// We will take n bytes, one byte for each character of output.
	output, randomness := make([]byte, n), make([]byte, n)

	// read all random
	_, err := rand.Read(randomness)
	if err != nil {
		return nil, err
	}

	for i, _lenEncodeStr := 0, uint8(len(encodeStr)); i < n; i++ {
		randomPos := uint8(randomness[i]) % _lenEncodeStr
		output[i] = encodeStr[randomPos]
	}
	randomness = nil

	return output, nil
}

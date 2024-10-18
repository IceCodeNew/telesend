package crypto

import (
	"crypto/rand"
)

// Imported from https://stackoverflow.com/a/44477359/13631331, with modifications

// len(encodeURL) == 64. This allows (x <= 256) x % 64 to have an even
// distribution.
const encodeURL = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"

// A helper function create and fill a slice of length n with characters from
// a-zA-Z0-9_-. It panics if there are any problems getting random bytes.
func RandAsciiBytes(n int) ([]byte, error) {
	// We will take n bytes, one byte for each character of output.
	output, randomness := make([]byte, n), make([]byte, n)

	// read all random
	_, err := rand.Read(randomness)
	if err != nil {
		return nil, err
	}

	for i := range output {
		randomPos := uint8(randomness[i]) % uint8(len(encodeURL))
		output[i] = encodeURL[randomPos]
	}

	return output, nil
}

package aead

import (
	"fmt"
	"math/rand/v2"
	"strconv"
	"strings"

	"github.com/IceCodeNew/telesend/internal/app/config"
	"github.com/IceCodeNew/telesend/pkg/crypto"
	"github.com/cloudflare/circl/cipher/ascon"
)

var (
	additionalData     = []byte("IceCode/TELESEND")
	block              *ascon.Cipher
	fakeKey, fakeNonce []byte
)
var (
	asconCipherNotInitErr = fmt.Errorf("ERROR: [Internal] Ascon cipher not initialized")
	invalidBotTokenErr    = fmt.Errorf("ERROR: [Internal] Invalid bot token format")
)

// TeleSend Use the specified telegram bot token to encrypt/decrypt the Bark Sender info
// just for convenience.
//
// The encryption is merely meant to prevent the database to be scanned easily.
func EncAscon128a(plaintext []byte) []byte {
	if block == nil {
		panic(asconCipherNotInitErr)
	}
	return block.Seal(plaintext[:0], fakeNonce, plaintext, additionalData)
}

// TeleSend Use the specified telegram bot token to encrypt/decrypt the Bark Sender info
// just for convenience.
//
// The encryption is merely meant to prevent the database to be scanned easily.
func DecAscon128a(ciphertext []byte) ([]byte, error) {
	if block == nil {
		return nil, asconCipherNotInitErr
	}
	plaintext, err := block.Open(ciphertext[:0], fakeNonce, ciphertext, additionalData)
	return plaintext, err
}

func InitCipher() error {
	_seed, passphrase, found := strings.Cut(config.TSConfig.BotToken, ":")
	if !found {
		return invalidBotTokenErr
	}
	r := rand.NewChaCha8(genFixedSeed(_seed))

	var err error
	fakeKey, fakeNonce = deriveKeyAndNonce(r, []byte(passphrase))
	if block, err = ascon.New(fakeKey, ascon.Ascon128a); err != nil {
		return err
	}
	return nil
}

func deriveKeyAndNonce(r *rand.ChaCha8, passphrase []byte) (key, nonce []byte) {
	salt := make([]byte, 8)
	r.Read(salt)

	kn := crypto.DeriveKey(passphrase, salt, ascon.KeySize+ascon.NonceSize)
	key, nonce = kn[:ascon.KeySize], kn[ascon.KeySize:]
	return
}

const seed2 = uint64(198964)

func genFixedSeed(s string) [32]byte {
	seed1, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		panic(err)
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

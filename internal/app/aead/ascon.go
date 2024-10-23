package aead

import (
	"fmt"

	"github.com/IceCodeNew/telesend/internal/app/config"
	"github.com/IceCodeNew/telesend/pkg/crypto"
	"github.com/cloudflare/circl/cipher/ascon"
)

var (
	additionalData     = []byte("IceCode/TELESEND")
	fakeKey, fakeNonce []byte
)
var (
	asconCipherNotInitErr = fmt.Errorf("ERROR: [Internal] Ascon cipher not initialized")
	block                 *ascon.Cipher
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

func deriveKeyAndNonce() error {
	var (
		botToken   = config.TSConfig.BotToken
		passphrase []byte
	)
	for i := 1; i < len(botToken); i++ {
		if botToken[i-1] == ':' {
			passphrase = []byte(botToken[i:])
			break
		}
	}
	if passphrase == nil {
		return fmt.Errorf("ERROR: [User Input] Invalid bot token")
	}

	// if lenPassphrase is multiples of ascon.NonceSize
	fakeSalt := additionalData

	lenPassphrase, _allDigitsOne := len(passphrase), ascon.NonceSize-1
	if lenPassphrase&_allDigitsOne != 0 {
		padLen := (lenPassphrase + _allDigitsOne) &^ _allDigitsOne
		padLen -= lenPassphrase
		fakeSalt = additionalData[:padLen]
	}

	kn := crypto.DeriveKey(passphrase, fakeSalt, ascon.KeySize+ascon.NonceSize)
	fakeKey, fakeNonce = kn[:ascon.KeySize], kn[ascon.KeySize:]
	return nil
}

func init() {
	err := deriveKeyAndNonce()
	if err != nil {
		panic(err)
	}

	block, err = ascon.New(fakeKey, ascon.Ascon128a)
	if err != nil {
		panic(err)
	}
}

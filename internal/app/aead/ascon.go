package aead

import (
	"github.com/cloudflare/circl/cipher/ascon"
)

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

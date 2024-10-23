package crypto

import (
	"golang.org/x/crypto/argon2"
)

func DeriveKey(passphrase, salt []byte, keyLen uint32) []byte {
	return argon2.IDKey(passphrase, salt, 1, 64*1024, 4, keyLen)
}

package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"

	"github.com/IceCodeNew/telesend/third_party/pkcs7"
)

const (
	KeySizeAES128 = 16
	KeySizeAES192 = 24
	KeySizeAES256 = 32

	maxLenOfWord = 9
	numWords     = 4
)

func verifyKeySize(key []byte) error {
	var err error

	_len := len(key)
	switch _len {
	case KeySizeAES128, KeySizeAES192, KeySizeAES256:
		// assertion ok
	default:
		err = fmt.Errorf(`
The length of AES key can ONLY be 16, 24, or 32 bytes, indicating AES-128, AES-192, or AES-256 accordingly.
But here we got direved key %s, which length is %d.`, key, _len)
	}
	return err
}

func EncryptWithAESCBC(asciiIV, asciiKey, msg []byte) (string, error) {
	if len(asciiIV) != aes.BlockSize {
		return "", fmt.Errorf("IV length must be %d", aes.BlockSize)
	}
	if err := verifyKeySize(asciiKey); err != nil {
		return "", err
	}
	block, err := aes.NewCipher(asciiKey)
	if err != nil {
		return "", err
	}

	// _allDigitsOne is 0b1111 (15)
	_allDigitsOne, lenMsg := aes.BlockSize-1, len(msg)
	// make the length of msg a multiple of aes.BlockSize
	// only works when all binary digits of _allDigitsOne is 1
	_cap := (lenMsg + _allDigitsOne) &^ _allDigitsOne

	data := make([]byte, lenMsg, _cap)
	copy(data, msg)

	// CBC mode works on blocks so plaintexts may need to be padded to the
	// next whole block. For an example of such padding, see
	// https://tools.ietf.org/html/rfc5246#section-6.2.3.2.
	data = pkcs7.Pad(aes.BlockSize, data)

	cbc := cipher.NewCBCEncrypter(block, asciiIV)
	cbc.CryptBlocks(data, data)

	return base64.StdEncoding.EncodeToString(data), nil
}

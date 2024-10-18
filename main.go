package main

import (
	"crypto/aes"
	"fmt"

	"github.com/IceCodeNew/telesend/pkg/crypto"
)

func main() {
	goCrypt()
}

func goCrypt() {
	asciiIV, err := crypto.RandAsciiBytes(aes.BlockSize)
	if err != nil {
		panic(err)
	}
	asciiKey, err := crypto.RandAsciiBytes(crypto.KeySizeAES256)
	if err != nil {
		panic(err)
	}

	plaintext := []byte(`{"body": "test", "sound": "birdsong"}`)
	enc, err := crypto.EncryptWithAESCBC(asciiIV, asciiKey, plaintext)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Encrypted text: %s\n", enc)
	fmt.Printf("IV: %s\n", string(asciiIV))
	fmt.Printf("Key: %s\n", string(asciiKey))
}

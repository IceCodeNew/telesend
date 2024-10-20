package main

import (
	"crypto/aes"
	"fmt"

	"github.com/IceCodeNew/telesend/internal/app/config"
	"github.com/IceCodeNew/telesend/internal/app/lotusdb"
	"github.com/IceCodeNew/telesend/pkg/crypto"
)

func init() {
	if err := config.ReadConfig(); err != nil {
		panic(err)
	}
}

func main() {
	db, err := lotusdb.InitDB(config.TSConfig.DbPath)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = db.Close()
	}()

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
	asciiIV, asciiKey = []byte("HHR67jpPwFIYl6h0"), []byte("Q3w0K0FJyQv7xwhZyQ9xJcQJVjZKRcVU")

	plaintext := `{"body": "test", "sound": "birdsong"}`
	enc, err := crypto.EncryptWithAESCBC(asciiIV, asciiKey, plaintext)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Encrypted text: %s\n", enc)
	fmt.Printf("IV: %s\n", string(asciiIV))
	fmt.Printf("Key: %s\n", string(asciiKey))
}

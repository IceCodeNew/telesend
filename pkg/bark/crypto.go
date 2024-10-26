package bark

import (
	"encoding/json"
	"net/url"

	"github.com/IceCodeNew/telesend/internal/app/aead"
	"github.com/IceCodeNew/telesend/pkg/crypto"
)

func (sender *BarkSender) queryFactor(msg *BarkMessage) (string, error) {
	plaintext, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}

	iv, key := sender.PreSharedSHA256IV, sender.PreSharedSHA256Key
	ciphertext, err := crypto.EncryptWithAESCBC(iv, key, plaintext)
	if err != nil {
		return "", err
	}

	params := url.Values{}
	params.Add("ciphertext", ciphertext)
	params.Add("iv", string(iv))

	iv, key = nil, nil
	return params.Encode(), nil
}

// TeleSend Use the specified telegram bot token to encrypt/decrypt the Bark Sender info
// just for convenience.
//
// The encryption is merely meant to prevent the database to be scanned easily.
func (sender *BarkSender) SelfEncrypt() error {
	deterministicKey, deterministicNonce, err := sender.deterministicKeyAndNonce()
	if err != nil {
		return err
	}

	deviceKey, err := aead.EncAscon128a(deterministicKey, deterministicNonce, sender.DeviceKey)
	if err != nil {
		return err
	}
	iv, err := aead.EncAscon128a(deterministicKey, deterministicNonce, sender.PreSharedSHA256IV)
	if err != nil {
		return err
	}
	key, err := aead.EncAscon128a(deterministicKey, deterministicNonce, sender.PreSharedSHA256Key)
	if err != nil {
		return err
	}

	_,
		sender.DeviceKey,
		sender.PreSharedSHA256IV,
		sender.PreSharedSHA256Key =
		0,
		deviceKey, iv, key

	deviceKey, iv, key = nil, nil, nil
	return nil
}

// TeleSend Use the specified telegram bot token to encrypt/decrypt the Bark Sender info
// just for convenience.
//
// The encryption is merely meant to prevent the database to be scanned easily.
func (sender *BarkSender) SelfDecrypt() error {
	deterministicKey, deterministicNonce, err := sender.deterministicKeyAndNonce()
	if err != nil {
		return err
	}

	deviceKey, err := aead.DecAscon128a(deterministicKey, deterministicNonce, sender.DeviceKey)
	if err != nil {
		return err
	}
	iv, err := aead.DecAscon128a(deterministicKey, deterministicNonce, sender.PreSharedSHA256IV)
	if err != nil {
		return err
	}
	key, err := aead.DecAscon128a(deterministicKey, deterministicNonce, sender.PreSharedSHA256Key)
	if err != nil {
		return err
	}

	_,
		sender.DeviceKey,
		sender.PreSharedSHA256IV,
		sender.PreSharedSHA256Key =
		0,
		deviceKey, iv, key

	deviceKey, iv, key = nil, nil, nil
	return nil
}

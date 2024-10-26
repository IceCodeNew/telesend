package bark

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/url"
	"strconv"
	"strings"

	"github.com/IceCodeNew/telesend/internal/app/aead"
	"github.com/IceCodeNew/telesend/internal/app/config"
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

func (sender *BarkSender) deterministicKeyAndNonce() (key, nonce []byte, err error) {
	_seed, _token, found := strings.Cut(config.TSConfig.BotToken, ":")
	if !found {
		return nil, nil, fmt.Errorf("ERROR: [Internal] Invalid bot token format")
	}

	passphrase := append(
		make([]byte, 0, len(_token)+len(sender.ID)),
		_token...,
	)
	passphrase = append(passphrase, sender.ID...)

	seed1, err := strconv.ParseUint(_seed, 10, 64)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"ERROR: [Internal] Failed to parse seed into uint64: %s", _seed,
		)
	}
	r := rand.NewChaCha8(
		aead.PredictableSeed(seed1, uint64(sender.Creator)),
	)
	salt := make([]byte, crypto.KeySizeAES128)
	r.Read(salt)

	_kn := aead.DeriveKey(passphrase, salt, crypto.KeySizeAES128<<1)
	key, nonce, _kn = _kn[:crypto.KeySizeAES128], _kn[crypto.KeySizeAES128:], nil
	return key, nonce, nil
}

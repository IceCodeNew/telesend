package bark

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/IceCodeNew/telesend/internal/app/aead"
	"github.com/IceCodeNew/telesend/internal/app/config"
	"github.com/IceCodeNew/telesend/pkg/crypto"
	"github.com/IceCodeNew/telesend/pkg/httpHelper"
	"github.com/samber/lo"
)

func (sender *BarkSender) Send(msg *BarkMessage, verbose bool) error {
	url, err := url.JoinPath(sender.Server, string(sender.DeviceKey))
	if err != nil {
		return err
	}
	body, err := sender.queryFactor(msg)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var resp *http.Response
	_, _, err = lo.AttemptWhileWithDelay(3, time.Second*10,
		func(int, time.Duration) (error, bool) {
			// BE AWARE as the resp is NOT GUARANTEED to be non-nil
			// resp, err = httpHelper.HttpReqHelper(http.MethodPost, url, nil, verbose)
			// do the http request, with additional headers & HTTP POST data
			resp, err = httpHelper.HttpReqHelper(req, verbose)
			if err != nil {
				return err, true
			}
			return nil, false
		})
	if err != nil {
		return fmt.Errorf("FATAL: failed to send message after 3 attempts, the last error was:\n %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	return nil
}

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
	predictableKey, predictableNonce, err := sender.predictableKeyAndNonce()
	if err != nil {
		return err
	}

	deviceKey, err := aead.EncAscon128a(predictableKey, predictableNonce, sender.DeviceKey)
	if err != nil {
		return err
	}
	iv, err := aead.EncAscon128a(predictableKey, predictableNonce, sender.PreSharedSHA256IV)
	if err != nil {
		return err
	}
	key, err := aead.EncAscon128a(predictableKey, predictableNonce, sender.PreSharedSHA256Key)
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
	predictableKey, predictableNonce, err := sender.predictableKeyAndNonce()
	if err != nil {
		return err
	}

	deviceKey, err := aead.DecAscon128a(predictableKey, predictableNonce, sender.DeviceKey)
	if err != nil {
		return err
	}
	iv, err := aead.DecAscon128a(predictableKey, predictableNonce, sender.PreSharedSHA256IV)
	if err != nil {
		return err
	}
	key, err := aead.DecAscon128a(predictableKey, predictableNonce, sender.PreSharedSHA256Key)
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

func (sender *BarkSender) predictableKeyAndNonce() (key, nonce []byte, err error) {
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

	key, nonce = aead.DeriveKeyAndNonce(passphrase, salt)
	return key, nonce, nil
}

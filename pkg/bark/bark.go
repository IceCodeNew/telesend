package bark

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/IceCodeNew/telesend/internal/app/aead"
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

func (sender *BarkSender) SelfEncrypt() {
	deviceKey, iv, key :=
		sender.DeviceKey,
		sender.PreSharedSHA256IV,
		sender.PreSharedSHA256Key

	_,
		sender.DeviceKey,
		sender.PreSharedSHA256IV,
		sender.PreSharedSHA256Key =
		0,
		aead.EncAscon128a(deviceKey),
		aead.EncAscon128a(iv),
		aead.EncAscon128a(key)

	deviceKey, iv, key = nil, nil, nil
}

func (sender *BarkSender) SelfDecrypt() error {
	deviceKey, err := aead.DecAscon128a(sender.DeviceKey)
	if err != nil {
		return err
	}
	iv, err := aead.DecAscon128a(sender.PreSharedSHA256IV)
	if err != nil {
		return err
	}
	key, err := aead.DecAscon128a(sender.PreSharedSHA256Key)
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

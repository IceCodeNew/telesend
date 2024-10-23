package notificator

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

type Sender struct {
	// Required; telegram user ID
	Creator int64
	// Required
	DeviceKey []byte
	// Required; generate by function UniqueID of pkg/uniqueID
	ID string
	// Required
	PreSharedSHA256IV []byte
	// Required
	PreSharedSHA256Key []byte
	// Required; default: "https://api.day.app/"
	Server string
}

type Message struct {
	// The number displayed next to App icon
	// Number greater than 9999 will be displayed as 9999+
	Badge int `json:"badge,omitempty"`
	// The content of the notification
	Body string `json:"body,omitempty"`
	// The value to be copied
	Copy string `json:"copy,omitempty"`
	// The group of the notification
	Group string `json:"group,omitempty"`
	// An url to the icon, available only on iOS 15 or later
	Icon string `json:"icon,omitempty"`
	// Value from https://github.com/Finb/Bark/tree/master/Sounds
	Sound string `json:"sound,omitempty"`
	// Notification title, optionally set by the sender
	Title string `json:"title,omitempty"`
	// Url that will jump when click notification
	URL string `json:"url,omitempty"`
}

func (sender *Sender) Send(msg *Message, verbose bool) error {
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

func (sender *Sender) queryFactor(msg *Message) (string, error) {
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

func (sender *Sender) SelfEncrypt() {
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

func (sender *Sender) SelfDecrypt() error {
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

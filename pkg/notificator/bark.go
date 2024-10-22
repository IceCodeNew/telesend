package notificator

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/IceCodeNew/telesend/pkg/crypto"
	"github.com/IceCodeNew/telesend/pkg/httpHelper"
	"github.com/samber/lo"
)

type Sender struct {
	// Required; telegram chat id
	Creator int64
	// Required
	DeviceKey string
	// Required; generate by function UniqueID of pkg/uniqueID
	ID string
	// Required
	PreSharedSHA256IV string
	// Required
	PreSharedSHA256Key string
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
	url, err := url.JoinPath(sender.Server, sender.DeviceKey)
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
	defer resp.Body.Close()
	return nil
}

func (sender *Sender) queryFactor(msg *Message) (string, error) {
	plaintext, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}

	iv, key := sender.PreSharedSHA256IV, sender.PreSharedSHA256Key
	ciphertext, err := crypto.EncryptWithAESCBC([]byte(iv), []byte(key), plaintext)
	if err != nil {
		return "", err
	}

	params := url.Values{}
	params.Add("ciphertext", ciphertext)
	params.Add("iv", iv)
	return params.Encode(), nil
}

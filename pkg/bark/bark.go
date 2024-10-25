package bark

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

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

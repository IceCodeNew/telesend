package httpHelper

import (
	"fmt"
	"io"
	"net/http"
)

func HttpReqHelper(method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	// it is not supposed to close the resp.Body under following sequences:
	// 1. resp is nil
	// 2. request completed successfully, and the status code is in the range [200, 300),
	//    in that case, users MUST close the resp.Body outside the function.
	if resp == nil {
		err = fmt.Errorf("response is nil, returned with error:\n%v", err)
		return nil, err
	}
	if err == nil && resp.StatusCode < 300 && resp.StatusCode >= 200 {
		// MAKE SURE the resp.Body is closed outside the function
		return resp, nil
	}
	// The http Client and Transport guarantee that Body is always non-nil, even on responses without a body or responses with a zero-length body.
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		err = fmt.Errorf(`
DEBUG: request returned with HTTP Status "%s",
DEBUG: the request URL is: "%s %s",
check the headers of the response:
%v
the original error is:
%v`,
			resp.Status,
			method, url,
			resp.Header,
			err)
	}
	return nil, err
}

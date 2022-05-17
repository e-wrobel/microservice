package streamer

import (
	"io"
	"io/ioutil"
	"net/http"
)

// sendRequest is functions to send data to Rest API
func sendRequest(addr, method string, body io.Reader) ([]byte, error) {
	c := http.Client{Timeout: timeout}
	req, err := http.NewRequest(method, addr, body)
	if err != nil {
		return nil, ErrNotPrepared
	}
	req.Header.Set("Content-type", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		return nil, ErrNotSent
	}
	defer resp.Body.Close()

	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrBodyCorrupted
	}

	return out, nil
}
